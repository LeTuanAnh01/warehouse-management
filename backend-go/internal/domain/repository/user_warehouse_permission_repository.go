package repository

import "warehouse-api/internal/domain/entity"

type UserWarehousePermissionRepository interface {
	Create(permission *entity.UserWarehousePermission) error
	Update(permission *entity.UserWarehousePermission) error
	Delete(userID, warehouseID uint) error
	DeleteByUserID(userID uint) error
	FindByUserID(userID uint) ([]entity.UserWarehousePermission, error)
	FindByWarehouseID(warehouseID uint) ([]entity.UserWarehousePermission, error)
	FindByUserAndWarehouse(userID, warehouseID uint) (*entity.UserWarehousePermission, error)
	CheckPermission(userID, warehouseID uint, permissionType string) (bool, error)
}
