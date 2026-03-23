package repository

import (
	"warehouse-api/internal/domain/entity"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	FindByCode(code string) (*entity.Product, error)
	FindAll(page, pageSize int) ([]entity.Product, int64, error)
	Update(product *entity.Product) error
	Delete(id uint) error
}

type InventoryRepository interface {
	FindByWarehouse(warehouseID uint) ([]entity.Inventory, error)
	FindByProduct(productID uint) ([]entity.Inventory, error)
	FindByWarehouseAndProduct(warehouseID, productID uint) (*entity.Inventory, error)
}
