package main

import (
	"fmt"
	"log"

	"warehouse-api/internal/config"
	httpDelivery "warehouse-api/internal/delivery/http"
	"warehouse-api/internal/delivery/http/handler"
	"warehouse-api/internal/infrastructure/database"
	"warehouse-api/internal/infrastructure/repository/postgres"
	"warehouse-api/internal/usecase"
)

func main() {
	log.Println("🚀 Starting PRODUCTION Server...")
	log.Println("� Authentication ENABLED - All endpoints require login")
	log.Println("")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	userRepo := postgres.NewUserRepository(db)
	warehouseRepo := postgres.NewWarehouseRepository(db)
	permRepo := postgres.NewUserWarehousePermissionRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, inventoryRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg)
	userUseCase := usecase.NewUserUseCase(userRepo, warehouseRepo, permRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase, cfg)
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)

	// Initialize router
	router := httpDelivery.NewRouter(cfg, productHandler, authHandler, userHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🚀 PRODUCTION Server running on http://%s", addr)
	log.Printf("📊 Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("🌍 Environment: %s", cfg.Server.Env)
	log.Println("")
	log.Println("📝 Available Endpoints:")
	log.Println("   Public (No Auth):")
	log.Printf("     - Health Check:  GET  http://%s/health", addr)
	log.Printf("     - Login:         POST http://%s/api/v1/auth/login", addr)
	log.Printf("     - Register:      POST http://%s/api/v1/auth/register", addr)
	log.Printf("     - Refresh Token: POST http://%s/api/v1/auth/refresh", addr)
	log.Println("")
	log.Println("   Protected (Auth Required):")
	log.Printf("     - Products CRUD: /api/v1/products")
	log.Printf("     - User Mgmt:     /api/v1/users (Admin only)")
	log.Printf("     - Logout:        POST http://%s/api/v1/auth/logout", addr)
	log.Printf("     - Change Pass:   POST http://%s/api/v1/auth/change-password", addr)
	log.Println("")
	log.Println("🔒 Protected endpoints require: Authorization: Bearer <token>")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("")

	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}