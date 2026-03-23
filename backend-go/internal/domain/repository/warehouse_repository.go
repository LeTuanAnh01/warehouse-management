package repository

import "warehouse-api/internal/domain/entity"

type WarehouseRepository interface {
	Create(warehouse *entity.Warehouse) error
	Update(warehouse *entity.Warehouse) error
	Delete(warehouse *entity.Warehouse) error
	FindAll(limit, offset int) ([]entity.Warehouse, int64, error)
	FindByID(id uint) (*entity.Warehouse, error)
	FindByCode(code string) (*entity.Warehouse, error)
}
