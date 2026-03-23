package handler

import (
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/usecase"
	"warehouse-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.authUseCase.Login(req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Login successful", result)
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} response.Response{data=dto.UserSummary}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.authUseCase.Register(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "User registered successfully", result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.authUseCase.RefreshToken(req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", result)
}

// Logout godoc
// @Summary User logout
// @Description Logout current user (client should discard tokens)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
// @Security BearerAuth
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// Client should discard the tokens
	// Optionally, implement token blacklist here
	response.Success(c, "Logout successful", gin.H{
		"message": "Please discard your access and refresh tokens",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Response
// @Router /auth/change-password [post]
// @Security BearerAuth
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.authUseCase.ChangePassword(userID.(uint), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "Password changed successfully", nil)
}