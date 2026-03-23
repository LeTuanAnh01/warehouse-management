package usecase

import (
	"errors"
	"time"
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"

	"gorm.io/gorm"
)

type ProductUseCase interface {
	GetProducts(page, pageSize int) ([]dto.ProductResponse, int64, error)
	GetProductByID(id uint) (*dto.ProductResponse, error)
	CreateProduct(req dto.CreateProductRequest) (*dto.ProductResponse, error)
	UpdateProduct(id uint, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
	GetInventoryByProduct(productID uint) ([]dto.InventoryStockResponse, error)
}

type productUseCase struct {
	productRepo   repository.ProductRepository
	inventoryRepo repository.InventoryRepository
}

func NewProductUseCase(
	productRepo repository.ProductRepository,
	inventoryRepo repository.InventoryRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *productUseCase) GetProducts(page, pageSize int) ([]dto.ProductResponse, int64, error) {
	products, total, err := u.productRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.ProductResponse
	for _, p := range products {
		response = append(response, dto.ProductResponse{
			ID:          p.ID,
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			Unit:        p.Unit,
			Category:    p.Category,
			MinStock:    p.MinStock,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, total, nil
}

func (u *productUseCase) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Unit:        product.Unit,
		Category:    product.Category,
		MinStock:    product.MinStock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) CreateProduct(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Check if code already exists
	existing, err := u.productRepo.FindByCode(req.Code)
	if err == nil && existing != nil {
		return nil, errors.New("product code already exists")
	}

	product := &entity.Product{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		Category:    req.Category,
		MinStock:    req.MinStock,
		IsActive:    true,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Unit:        product.Unit,
		Category:    product.Category,
		MinStock:    product.MinStock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) UpdateProduct(id uint, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.MinStock > 0 {
		product.MinStock = req.MinStock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Unit:        product.Unit,
		Category:    product.Category,
		MinStock:    product.MinStock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) DeleteProduct(id uint) error {
	_, err := u.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return u.productRepo.Delete(id)
}

func (u *productUseCase) GetInventoryByProduct(productID uint) ([]dto.InventoryStockResponse, error) {
	inventory, err := u.inventoryRepo.FindByProduct(productID)
	if err != nil {
		return nil, err
	}

	var response []dto.InventoryStockResponse
	for _, inv := range inventory {
		status := "ok"
		if inv.Quantity == 0 {
			status = "out"
		} else if inv.Product.MinStock > 0 && inv.Quantity <= inv.Product.MinStock {
			status = "low"
		}

		response = append(response, dto.InventoryStockResponse{
			ProductID:     inv.ProductID,
			ProductCode:   inv.Product.Code,
			ProductName:   inv.Product.Name,
			WarehouseID:   inv.WarehouseID,
			WarehouseName: inv.Warehouse.Name,
			Quantity:      inv.Quantity,
			MinStock:      inv.Product.MinStock,
			Status:        status,
		})
	}

	return response, nil
}
