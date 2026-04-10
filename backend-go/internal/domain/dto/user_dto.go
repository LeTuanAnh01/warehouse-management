package dto

import "time"

type UserSummary struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FullName     string `json:"full_name" binding:"required"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	RoleID       uint   `json:"role_id" binding:"required"`
	WarehouseIDs []uint `json:"warehouse_ids"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
	IsActive *bool  `json:"is_active"`
}

type UserDetailResponse struct {
	ID         uint                          `json:"id"`
	Username   string                        `json:"username"`
	Email      string                        `json:"email"`
	FullName   string                        `json:"full_name"`
	Phone      string                        `json:"phone"`
	Role       string                        `json:"role"`
	RoleID     uint                          `json:"role_id"`
	IsActive   bool                          `json:"is_active"`
	Warehouses []WarehousePermissionResponse `json:"warehouses"`
}

type PaginatedUsersResponse struct {
	Users      []UserSummary `json:"users"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int64         `json:"total_pages"`
}

type AssignWarehouseRequest struct {
	UserID      uint `json:"user_id" binding:"required"`
	WarehouseID uint `json:"warehouse_id" binding:"required"`
	CanView     bool `json:"can_view"`
	CanImport   bool `json:"can_import"`
	CanExport   bool `json:"can_export"`
	CanManage   bool `json:"can_manage"`
}

type WarehousePermissionResponse struct {
	WarehouseID   uint   `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	WarehouseCode string `json:"warehouse_code"`
	CanView       bool   `json:"can_view"`
	CanImport     bool   `json:"can_import"`
	CanExport     bool   `json:"can_export"`
	CanManage     bool   `json:"can_manage"`
}

type UserWarehousePermissionRequest struct {
	WarehouseID uint `json:"warehouse_id" binding:"required"`
	CanView     bool `json:"can_view"`
	CanEdit     bool `json:"can_edit"`
	CanDelete   bool `json:"can_delete"`
	CanManage   bool `json:"can_manage"`
	CanImport   bool `json:"can_import"`
	CanExport   bool `json:"can_export"`
}
