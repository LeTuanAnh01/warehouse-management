package handler

import (
	"strconv"
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/usecase"
	"warehouse-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get paginated list of users (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=dto.PaginatedUsersResponse}
// @Router /users [get]
// @Security BearerAuth
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.userUseCase.GetAllUsers(page, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Users retrieved successfully", result)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get detailed user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=dto.UserDetailResponse}
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.userUseCase.GetUserByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", result)
}

// CreateUser godoc
// @Summary Create new user
// @Description Create a new user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=dto.UserDetailResponse}
// @Router /users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.userUseCase.CreateUser(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "User created successfully", result)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response{data=dto.UserDetailResponse}
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.userUseCase.UpdateUser(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "User updated successfully", result)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Router /users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// AssignWarehouse godoc
// @Summary Assign warehouse to user
// @Description Assign warehouse permissions to user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param assignment body dto.AssignWarehouseRequest true "Assignment data"
// @Success 200 {object} response.Response
// @Router /users/assign-warehouse [post]
// @Security BearerAuth
func (h *UserHandler) AssignWarehouse(c *gin.Context) {
	var req dto.AssignWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.userUseCase.AssignWarehouse(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "Warehouse assigned successfully", nil)
}

// RemoveWarehouse godoc
// @Summary Remove warehouse from user
// @Description Remove warehouse permissions from user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param warehouse_id path int true "Warehouse ID"
// @Success 200 {object} response.Response
// @Router /users/{user_id}/warehouses/{warehouse_id} [delete]
// @Security BearerAuth
func (h *UserHandler) RemoveWarehouse(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	warehouseID, err := strconv.ParseUint(c.Param("warehouse_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid warehouse ID")
		return
	}

	if err := h.userUseCase.RemoveWarehouse(uint(userID), uint(warehouseID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "Warehouse removed successfully", nil)
}

// GetUserWarehouses godoc
// @Summary Get user's warehouses
// @Description Get all warehouses assigned to a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=[]dto.WarehousePermissionResponse}
// @Router /users/{id}/warehouses [get]
// @Security BearerAuth
func (h *UserHandler) GetUserWarehouses(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.userUseCase.GetUserWarehouses(uint(id))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Warehouses retrieved successfully", result)
}
