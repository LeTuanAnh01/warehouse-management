package postgres

import (
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"

	"gorm.io/gorm"
)

type warehouseRepositoryImpl struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) repository.WarehouseRepository {
	return &warehouseRepositoryImpl{db: db}
}

func (r *warehouseRepositoryImpl) Create(warehouse *entity.Warehouse) error {
	return r.db.Create(warehouse).Error
}

func (r *warehouseRepositoryImpl) Update(warehouse *entity.Warehouse) error {
	return r.db.Save(warehouse).Error
}

func (r *warehouseRepositoryImpl) Delete(warehouse *entity.Warehouse) error {
	return r.db.Delete(warehouse).Error
}

func (r *warehouseRepositoryImpl) FindAll(limit, offset int) ([]entity.Warehouse, int64, error) {
	var warehouses []entity.Warehouse
	var total int64

	if err := r.db.Model(&entity.Warehouse{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Manager").Preload("Manager.Role").
		Limit(limit).Offset(offset).
		Find(&warehouses).Error

	return warehouses, total, err
}

func (r *warehouseRepositoryImpl) FindByID(id uint) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.Preload("Manager").Preload("Manager.Role").
		First(&warehouse, id).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepositoryImpl) FindByCode(code string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.Preload("Manager").Preload("Manager.Role").
		Where("code = ?", code).First(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}
