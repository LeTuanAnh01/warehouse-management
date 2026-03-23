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
	log.Println("🧪 Starting TEST Server...")
	log.Println("⚠️  This is for TESTING ONLY - Not for production use")
	log.Println("")

	// Load configuration
	cfg := config.Load()

	// Debug configuration
	log.Printf("📋 Configuration loaded:")
	log.Printf("   Server: %s:%s (env: %s)", cfg.Server.Host, cfg.Server.Port, cfg.Server.Env)
	log.Printf("   Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Println("")

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, inventoryRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase, cfg)

	// Initialize router
	router := httpDelivery.NewRouter(cfg, productHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🧪 TEST Server running on http://%s", addr)
	log.Printf("📊 Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Println("")
	log.Println("📝 Test Endpoints:")
	log.Printf("   - Health Check:    GET  http://%s/health", addr)
	log.Printf("   - Products List:   GET  http://%s/api/v1/products", addr)
	log.Printf("   - Product Detail:  GET  http://%s/api/v1/products/:id", addr)
	log.Printf("   - Create Product:  POST http://%s/api/v1/products", addr)
	log.Println("")
	log.Println("⚠️  No Authentication - All endpoints are open for testing")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("")

	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
