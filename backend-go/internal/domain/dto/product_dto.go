package dto

type CreateProductRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Category    string `json:"category"`
	MinStock    int    `json:"min_stock"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Category    string `json:"category"`
	MinStock    int    `json:"min_stock"`
	IsActive    *bool  `json:"is_active"`
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Category    string `json:"category"`
	MinStock    int    `json:"min_stock"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type InventoryStockResponse struct {
	ProductID     uint   `json:"product_id"`
	ProductCode   string `json:"product_code"`
	ProductName   string `json:"product_name"`
	WarehouseID   uint   `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	Quantity      int    `json:"quantity"`
	MinStock      int    `json:"min_stock"`
	Status        string `json:"status"` // "ok", "low", "out"
}
