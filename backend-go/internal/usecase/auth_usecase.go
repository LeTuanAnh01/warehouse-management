package usecase

import (
	"errors"
	"warehouse-api/internal/config"
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"
	"warehouse-api/pkg/utils"

	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) (*dto.UserSummary, error)
	RefreshToken(req dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	ChangePassword(userID uint, req dto.ChangePasswordRequest) error
}

type authUseCase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, cfg *config.Config) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *authUseCase) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by username
	user, err := u.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Username, user.RoleID, u.cfg)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, user.RoleID, u.cfg)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Update last login
	_ = u.userRepo.UpdateLastLogin(user.ID)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.cfg.JWT.Expiration.Seconds()),
		User: dto.UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role.Name,
		},
	}, nil
}

func (u *authUseCase) Register(req dto.RegisterRequest) (*dto.UserSummary, error) {
	// Check if username already exists
	if _, err := u.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := u.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user with default role (Viewer = 4)
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		RoleID:       4, // Default: Viewer
		IsActive:     true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Reload user with role
	user, err = u.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserSummary{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role.Name,
	}, nil
}

func (u *authUseCase) RefreshToken(req dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, u.cfg)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate new tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Username, user.RoleID, u.cfg)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, user.RoleID, u.cfg)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.cfg.JWT.Expiration.Seconds()),
		User: dto.UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role.Name,
		},
	}, nil
}

func (u *authUseCase) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = hashedPassword
	return u.userRepo.Update(user)
}
