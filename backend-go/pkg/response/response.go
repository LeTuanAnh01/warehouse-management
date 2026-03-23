package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// Success responses
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Paginated response
func Paginated(c *gin.Context, data interface{}, page, limit int, totalRecords int64) {
	totalPages := int(totalRecords) / limit
	if int(totalRecords)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:         page,
			Limit:        limit,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
			HasNext:      page < totalPages,
			HasPrevious:  page > 1,
		},
	})
}

// Error responses
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message)
}
