# API Testing Guide

Server đã chạy thành công! 🎉

## Server Information
- **Base URL**: `http://localhost:8080`
- **Environment**: Development
- **Database**: PostgreSQL (warehouse_db)

## Available Endpoints

### 1. Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Warehouse API is running"
}
```

---

### 2. Get All Products (với phân trang)
```bash
# Trang 1, 20 items
curl http://localhost:8080/api/v1/products

# Trang 2, 10 items mỗi trang
curl http://localhost:8080/api/v1/products?page=2&page_size=10
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "SP001",
      "name": "Laptop Dell XPS 15",
      "description": "Laptop cao cấp Dell XPS 15 inch",
      "unit": "Cái",
      "category": "Electronics",
      "min_stock": 10,
      "is_active": true,
      "created_at": "2026-03-09T11:00:28Z",
      "updated_at": "2026-03-09T11:00:28Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_records": 10,
    "total_pages": 1,
    "has_next": false,
    "has_previous": false
  }
}
```

---

### 3. Get Product By ID
```bash
curl http://localhost:8080/api/v1/products/1
```

**Response:**
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "code": "SP001",
    "name": "Laptop Dell XPS 15",
    "description": "Laptop cao cấp Dell XPS 15 inch",
    "unit": "Cái",
    "category": "Electronics",
    "min_stock": 10,
    "is_active": true,
    "created_at": "2026-03-09T11:00:28Z",
    "updated_at": "2026-03-09T11:00:28Z"
  }
}
```

---

### 4. Create New Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SP011",
    "name": "MacBook Pro M3",
    "description": "MacBook Pro 16 inch M3 Max",
    "unit": "Cái",
    "category": "Electronics",
    "min_stock": 5
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 11,
    "code": "SP011",
    "name": "MacBook Pro M3",
    "description": "MacBook Pro 16 inch M3 Max",
    "unit": "Cái",
    "category": "Electronics",
    "min_stock": 5,
    "is_active": true,
    "created_at": "2026-03-09T13:32:39+07:00",
    "updated_at": "2026-03-09T13:32:39+07:00"
  }
}
```

---

### 5. Update Product
```bash
curl -X PUT http://localhost:8080/api/v1/products/11 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro M3 Updated",
    "description": "MacBook Pro 16 inch M3 Max - Updated",
    "min_stock": 10
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": 11,
    "code": "SP011",
    "name": "MacBook Pro M3 Updated",
    "description": "MacBook Pro 16 inch M3 Max - Updated",
    "unit": "Cái",
    "category": "Electronics",
    "min_stock": 10,
    "is_active": true,
    "created_at": "2026-03-09T13:32:39+07:00",
    "updated_at": "2026-03-09T13:35:21+07:00"
  }
}
```

---

### 6. Delete Product
```bash
curl -X DELETE http://localhost:8080/api/v1/products/11
```

**Response:**
```json
{
  "success": true,
  "message": "Product deleted successfully",
  "data": {
    "message": "Product deleted successfully"
  }
}
```

---

### 7. Get Product Inventory (xem tồn kho theo sản phẩm)
```bash
curl http://localhost:8080/api/v1/products/1/inventory
```

**Response:**
```json
{
  "success": true,
  "message": "Inventory retrieved successfully",
  "data": [
    {
      "product_id": 1,
      "product_code": "SP001",
      "product_name": "Laptop Dell XPS 15",
      "warehouse_id": 1,
      "warehouse_name": "Kho Hà Nội",
      "quantity": 30,
      "min_stock": 10,
      "status": "ok"
    },
    {
      "product_id": 1,
      "product_code": "SP001",
      "product_name": "Laptop Dell XPS 15",
      "warehouse_id": 2,
      "warehouse_name": "Kho Hồ Chí Minh",
      "quantity": 20,
      "min_stock": 10,
      "status": "ok"
    }
  ]
}
```

**Inventory Status:**
- `ok`: Số lượng đủ (> min_stock)
- `low`: Số lượng thấp (<= min_stock)
- `out`: Hết hàng (quantity = 0)

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": "Product not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## Testing with Postman

1. Import các request vào Postman
2. Set base URL: `http://localhost:8080`
3. Đối với POST/PUT requests, set header: `Content-Type: application/json`

---

## Next Steps

✅ **Đã hoàn thành:**
- Health check endpoint
- Product CRUD endpoints
- Product inventory endpoint
- Pagination support
- Error handling
- Database connection

⏳ **Cần làm tiếp:**
- Authentication (Login/Logout)
- Warehouse endpoints
- Transaction endpoints (Nhập/Xuất hàng)
- User management & permissions
- Reports
- Swagger documentation
- Unit tests
