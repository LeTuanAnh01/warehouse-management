package usecase

import (
	"errors"
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"
	"warehouse-api/pkg/utils"

	"gorm.io/gorm"
)

type UserUseCase interface {
	GetAllUsers(page, limit int) (*dto.PaginatedUsersResponse, error)
	GetUserByID(id uint) (*dto.UserDetailResponse, error)
	CreateUser(req dto.CreateUserRequest) (*dto.UserDetailResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserDetailResponse, error)
	DeleteUser(id uint) error
	AssignWarehouse(req dto.AssignWarehouseRequest) error
	RemoveWarehouse(userID, warehouseID uint) error
	GetUserWarehouses(userID uint) ([]dto.WarehousePermissionResponse, error)
}

type userUseCase struct {
	userRepo      repository.UserRepository
	warehouseRepo repository.WarehouseRepository
	permRepo      repository.UserWarehousePermissionRepository
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	warehouseRepo repository.WarehouseRepository,
	permRepo repository.UserWarehousePermissionRepository,
) UserUseCase {
	return &userUseCase{
		userRepo:      userRepo,
		warehouseRepo: warehouseRepo,
		permRepo:      permRepo,
	}
}

func (u *userUseCase) GetAllUsers(page, limit int) (*dto.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, total, err := u.userRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	userSummaries := make([]dto.UserSummary, len(users))
	for i, user := range users {
		userSummaries[i] = dto.UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			Role:     user.Role.Name,
			IsActive: user.IsActive,
		}
	}

	return &dto.PaginatedUsersResponse{
		Users:       userSummaries,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  (total + int64(limit) - 1) / int64(limit),
	}, nil
}

func (u *userUseCase) GetUserByID(id uint) (*dto.UserDetailResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Get user's warehouse permissions
	permissions, err := u.permRepo.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	warehouses := make([]dto.WarehousePermissionResponse, len(permissions))
	for i, perm := range permissions {
		warehouses[i] = dto.WarehousePermissionResponse{
			WarehouseID:   perm.WarehouseID,
			WarehouseName: perm.Warehouse.Name,
			WarehouseCode: perm.Warehouse.Code,
			CanView:       perm.CanView,
			CanImport:     perm.CanImport,
			CanExport:     perm.CanExport,
			CanManage:     perm.CanManage,
		}
	}

	return &dto.UserDetailResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FullName:   user.FullName,
		Phone:      user.Phone,
		Role:       user.Role.Name,
		RoleID:     user.RoleID,
		IsActive:   user.IsActive,
		Warehouses: warehouses,
	}, nil
}

func (u *userUseCase) CreateUser(req dto.CreateUserRequest) (*dto.UserDetailResponse, error) {
	// Check if username exists
	existing, _ := u.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	existing, _ = u.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign warehouses if provided
	if len(req.WarehouseIDs) > 0 {
		for _, whID := range req.WarehouseIDs {
			perm := &entity.UserWarehousePermission{
				UserID:      user.ID,
				WarehouseID: whID,
				CanView:     true,
				CanImport:   req.RoleID <= 2, // ADMIN or MANAGER
				CanExport:   req.RoleID <= 2,
				CanManage:   req.RoleID <= 2,
			}
			u.permRepo.Create(perm)
		}
	}

	return u.GetUserByID(user.ID)
}

func (u *userUseCase) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserDetailResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if req.Email != "" && req.Email != user.Email {
		existing, _ := u.userRepo.FindByEmail(req.Email)
		if existing != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return u.GetUserByID(user.ID)
}

func (u *userUseCase) DeleteUser(id uint) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Don't allow deleting admin users
	if user.RoleID == 1 {
		return errors.New("cannot delete admin user")
	}

	// Delete warehouse permissions first
	u.permRepo.DeleteByUserID(id)

	return u.userRepo.Delete(user)
}

func (u *userUseCase) AssignWarehouse(req dto.AssignWarehouseRequest) error {
	// Check if user exists
	_, err := u.userRepo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if warehouse exists
	_, err = u.warehouseRepo.FindByID(req.WarehouseID)
	if err != nil {
		return errors.New("warehouse not found")
	}

	// Check if permission already exists
	existing, _ := u.permRepo.FindByUserAndWarehouse(req.UserID, req.WarehouseID)
	if existing != nil {
		// Update existing
		existing.CanView = req.CanView
		existing.CanImport = req.CanImport
		existing.CanExport = req.CanExport
		existing.CanManage = req.CanManage
		return u.permRepo.Update(existing)
	}

	// Create new permission
	perm := &entity.UserWarehousePermission{
		UserID:      req.UserID,
		WarehouseID: req.WarehouseID,
		CanView:     req.CanView,
		CanImport:   req.CanImport,
		CanExport:   req.CanExport,
		CanManage:   req.CanManage,
	}

	return u.permRepo.Create(perm)
}

func (u *userUseCase) RemoveWarehouse(userID, warehouseID uint) error {
	return u.permRepo.Delete(userID, warehouseID)
}

func (u *userUseCase) GetUserWarehouses(userID uint) ([]dto.WarehousePermissionResponse, error) {
	permissions, err := u.permRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.WarehousePermissionResponse, len(permissions))
	for i, perm := range permissions {
		result[i] = dto.WarehousePermissionResponse{
			WarehouseID:   perm.WarehouseID,
			WarehouseName: perm.Warehouse.Name,
			WarehouseCode: perm.Warehouse.Code,
			CanView:       perm.CanView,
			CanImport:     perm.CanImport,
			CanExport:     perm.CanExport,
			CanManage:     perm.CanManage,
		}
	}

	return result, nil
}
