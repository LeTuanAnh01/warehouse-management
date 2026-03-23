# Warehouse Management API - Backend

Backend API cho hệ thống quản lý kho, được xây dựng với Golang theo Clean Architecture.

## 🏗️ Tech Stack

- **Language**: Go 1.26+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 16
- **ORM**: GORM
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Validation**: go-playground/validator

## 📁 Project Structure

```
backend-go/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/                     # Configuration
│   ├── domain/                     # Domain models & interfaces
│   │   ├── entity/                 # Entities (database models)
│   │   ├── dto/                    # Data Transfer Objects
│   │   └── repository/             # Repository interfaces
│   ├── usecase/                    # Business logic
│   ├── delivery/                   # HTTP handlers
│   │   ├── http/
│   │   │   ├── handler/           # HTTP handlers
│   │   │   ├── middleware/        # Middleware
│   │   │   └── router/            # Routes
│   └── infrastructure/             # External dependencies
│       ├── database/              # Database connection
│       ├── repository/            # Repository implementations
│       └── logger/                # Logging
├── pkg/                            # Public packages
│   ├── utils/                     # Utilities
│   ├── response/                  # Response helpers
│   └── validator/                 # Validation helpers
├── docs/                           # Swagger documentation
├── scripts/                        # Scripts (migration, seed)
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.26 trở lên
- PostgreSQL 16
- Make (optional)

### Installation

1. Clone repository
```bash
cd backend-go
```

2. Install dependencies
```bash
go mod download
```

3. Setup environment
```bash
cp .env.example .env
# Chỉnh sửa .env theo môi trường của bạn
```

4. Run database migrations (database đã setup)
```bash
# Database đã được setup ở folder ../database
```

5. Run the application
```bash
go run cmd/api/main.go
```

Server sẽ chạy tại: `http://localhost:8080`

## 📚 API Documentation

Sau khi chạy server, truy cập Swagger UI tại:
```
http://localhost:8080/swagger/index.html
```

## 🔑 Authentication

API sử dụng JWT cho authentication:

1. **Login**: POST `/api/v1/auth/login`
2. Lấy `access_token` từ response
3. Thêm vào header: `Authorization: Bearer {access_token}`

### Default Test Accounts

| Username | Password | Role |
|----------|----------|------|
| admin | 123456 | ADMIN |
| manager1 | 123456 | WAREHOUSE_MANAGER |
| staff1 | 123456 | WAREHOUSE_STAFF |

## 🛣️ API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/auth/me` - Get current user

### Products
- `GET /api/v1/products` - List products (with pagination)
- `GET /api/v1/products/:id` - Get product detail
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `GET /api/v1/products/search` - Search products

### Warehouses
- `GET /api/v1/warehouses` - List warehouses
- `GET /api/v1/warehouses/:id` - Get warehouse detail
- `POST /api/v1/warehouses` - Create warehouse
- `PUT /api/v1/warehouses/:id` - Update warehouse

### Inventory
- `GET /api/v1/inventory` - List inventory (with pagination)
- `GET /api/v1/inventory/:id` - Get inventory detail
- `GET /api/v1/inventory/warehouse/:warehouse_id` - Get by warehouse
- `GET /api/v1/inventory/product/:product_id` - Get by product
- `GET /api/v1/inventory/low-stock` - Low stock alerts

### Transactions
- `GET /api/v1/transactions` - List transactions (with pagination)
- `GET /api/v1/transactions/:id` - Get transaction detail
- `POST /api/v1/transactions/import` - Create import transaction
- `POST /api/v1/transactions/export` - Create export transaction
- `GET /api/v1/transactions/history/:product_id` - Product history

### Users
- `GET /api/v1/users` - List users (Admin only)
- `GET /api/v1/users/:id` - Get user detail
- `POST /api/v1/users` - Create user (Admin only)
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (Admin only)

### Reports
- `GET /api/v1/reports/inventory-value` - Inventory value report
- `GET /api/v1/reports/transactions-summary` - Transactions summary
- `GET /api/v1/reports/low-stock` - Low stock report

## 🔒 Permissions

Hệ thống có 4 levels quyền:

1. **ADMIN** - Full quyền
2. **WAREHOUSE_MANAGER** - Quản lý kho được phân quyền
3. **WAREHOUSE_STAFF** - Nhập/xuất hàng
4. **VIEWER** - Chỉ xem

Quyền theo kho:
- `can_view` - Xem thông tin
- `can_import` - Nhập hàng
- `can_export` - Xuất hàng
- `can_manage` - Quản lý kho

## 📊 Pagination

Tất cả list APIs đều hỗ trợ pagination:

Query Parameters:
- `page` - Số trang (default: 1)
- `limit` - Số records mỗi trang (default: 20, max: 100)
- `sort` - Sắp xếp (e.g., `created_at desc`)

Response Format:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_records": 150,
    "total_pages": 8,
    "has_next": true,
    "has_previous": false
  }
}
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/usecase/...
```

## 🔨 Development

### Code Structure

- **Domain Layer**: Business entities và interfaces
- **Use Case Layer**: Business logic
- **Delivery Layer**: HTTP handlers, middleware
- **Infrastructure Layer**: Database, external services

### Adding New Feature

1. Define entity trong `internal/domain/entity/`
2. Define DTO trong `internal/domain/dto/`
3. Define repository interface trong `internal/domain/repository/`
4. Implement repository trong `internal/infrastructure/repository/`
5. Implement use case trong `internal/usecase/`
6. Implement handler trong `internal/delivery/http/handler/`
7. Register routes trong `internal/delivery/http/router/`

## 📝 Logging

Application sử dụng structured logging:

```go
log.Info("User logged in", "user_id", userID, "username", username)
log.Error("Database error", "error", err)
```

Log levels: `debug`, `info`, `warn`, `error`, `fatal`

## 🚀 Deployment

### Build

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Run binary
./bin/api
```

### Docker

```bash
# Build image
docker build -t warehouse-api .

# Run container
docker run -p 8080:8080 --env-file .env warehouse-api
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

MIT License

## 👥 Contact

For support or questions, please contact the development team.

---

**Status**: 🚧 Under Development
