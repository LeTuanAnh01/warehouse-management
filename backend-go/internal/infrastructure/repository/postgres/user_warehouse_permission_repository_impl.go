package postgres

import (
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"

	"gorm.io/gorm"
)

type userWarehousePermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewUserWarehousePermissionRepository(db *gorm.DB) repository.UserWarehousePermissionRepository {
	return &userWarehousePermissionRepositoryImpl{db: db}
}

func (r *userWarehousePermissionRepositoryImpl) Create(permission *entity.UserWarehousePermission) error {
	return r.db.Create(permission).Error
}

func (r *userWarehousePermissionRepositoryImpl) Update(permission *entity.UserWarehousePermission) error {
	return r.db.Save(permission).Error
}

func (r *userWarehousePermissionRepositoryImpl) Delete(userID, warehouseID uint) error {
	return r.db.Where("user_id = ? AND warehouse_id = ?", userID, warehouseID).
		Delete(&entity.UserWarehousePermission{}).Error
}

func (r *userWarehousePermissionRepositoryImpl) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.UserWarehousePermission{}).Error
}

func (r *userWarehousePermissionRepositoryImpl) FindByUserID(userID uint) ([]entity.UserWarehousePermission, error) {
	var permissions []entity.UserWarehousePermission
	err := r.db.Preload("Warehouse").Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}

func (r *userWarehousePermissionRepositoryImpl) FindByWarehouseID(warehouseID uint) ([]entity.UserWarehousePermission, error) {
	var permissions []entity.UserWarehousePermission
	err := r.db.Preload("User").Preload("User.Role").
		Where("warehouse_id = ?", warehouseID).Find(&permissions).Error
	return permissions, err
}

func (r *userWarehousePermissionRepositoryImpl) FindByUserAndWarehouse(userID, warehouseID uint) (*entity.UserWarehousePermission, error) {
	var permission entity.UserWarehousePermission
	err := r.db.Preload("Warehouse").
		Where("user_id = ? AND warehouse_id = ?", userID, warehouseID).
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *userWarehousePermissionRepositoryImpl) CheckPermission(userID, warehouseID uint, permissionType string) (bool, error) {
	var permission entity.UserWarehousePermission
	err := r.db.Where("user_id = ? AND warehouse_id = ?", userID, warehouseID).
		First(&permission).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	switch permissionType {
	case "view":
		return permission.CanView, nil
	case "import":
		return permission.CanImport, nil
	case "export":
		return permission.CanExport, nil
	case "manage":
		return permission.CanManage, nil
	default:
		return false, nil
	}
}
