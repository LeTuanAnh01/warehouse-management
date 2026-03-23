package postgres

import (
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByCode(code string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Where("code = ?", code).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(page, pageSize int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	// Count total
	if err := r.db.Model(&entity.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) repository.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) FindByWarehouse(warehouseID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	err := r.db.Where("warehouse_id = ?", warehouseID).
		Preload("Product").
		Preload("Warehouse").
		Find(&inventory).Error
	return inventory, err
}

func (r *inventoryRepository) FindByProduct(productID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	err := r.db.Where("product_id = ?", productID).
		Preload("Product").
		Preload("Warehouse").
		Find(&inventory).Error
	return inventory, err
}

func (r *inventoryRepository) FindByWarehouseAndProduct(warehouseID, productID uint) (*entity.Inventory, error) {
	var inventory entity.Inventory
	err := r.db.Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
		Preload("Product").
		Preload("Warehouse").
		First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}
