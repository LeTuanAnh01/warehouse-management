package entity

import "time"

type UserWarehousePermission struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id" gorm:"not null;index"`
    WarehouseID uint      `json:"warehouse_id" gorm:"not null;index"`
    CanView     bool      `json:"can_view" gorm:"default:true"`
    CanEdit     bool      `json:"can_edit" gorm:"default:false"`
    CanDelete   bool      `json:"can_delete" gorm:"default:false"`
    CanManage   bool      `json:"can_manage" gorm:"default:false"`
    CanImport   bool      `json:"can_import" gorm:"default:false"`
    CanExport   bool      `json:"can_export" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // Relations
    User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Warehouse Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

func (UserWarehousePermission) TableName() string {
    return "user_warehouse_permissions"
}