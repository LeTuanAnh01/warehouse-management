package handler

import (
	"net/http"
	"strconv"
	"warehouse-api/internal/config"
	"warehouse-api/internal/domain/dto"
	"warehouse-api/internal/usecase"
	"warehouse-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
	config         *config.Config
}

func NewProductHandler(productUseCase usecase.ProductUseCase, cfg *config.Config) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		config:         cfg,
	}
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.PaginatedResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.config.App.DefaultPageSize)))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > h.config.App.MaxPageSize {
		pageSize = h.config.App.DefaultPageSize
	}

	products, total, err := h.productUseCase.GetProducts(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch products")
		return
	}

	response.Paginated(c, products, page, pageSize, total)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	product, err := h.productUseCase.GetProductByID(uint(id))
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} response.Response{data=dto.ProductResponse}
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	product, err := h.productUseCase.CreateProduct(req)
	if err != nil {
		if err.Error() == "product code already exists" {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to create product")
		return
	}

	response.Created(c, "Product created successfully", product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product data"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	product, err := h.productUseCase.UpdateProduct(uint(id), req)
	if err != nil {
		if err.Error() == "product not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to update product")
		return
	}

	response.Success(c, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	if err := h.productUseCase.DeleteProduct(uint(id)); err != nil {
		if err.Error() == "product not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to delete product")
		return
	}

	response.Success(c, "Product deleted successfully", gin.H{"message": "Product deleted successfully"})
}

// GetInventoryByProduct godoc
// @Summary Get inventory by product
// @Description Get inventory levels for a product across all warehouses
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=[]dto.InventoryStockResponse}
// @Router /products/{id}/inventory [get]
func (h *ProductHandler) GetInventoryByProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	inventory, err := h.productUseCase.GetInventoryByProduct(uint(id))
	if err != nil {
		response.InternalServerError(c, "Failed to fetch inventory")
		return
	}

	response.Success(c, "Inventory retrieved successfully", inventory)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Warehouse API is running",
	})
}
