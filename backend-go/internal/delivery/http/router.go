package http

import (
	"warehouse-api/internal/config"
	"warehouse-api/internal/delivery/http/handler"
	"warehouse-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	productHandler *handler.ProductHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Auth routes
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			// Product routes
			products := protected.Group("/products")
			{
				products.GET("", productHandler.GetProducts)
				products.GET("/:id", productHandler.GetProductByID)
				products.POST("", productHandler.CreateProduct)
				products.PUT("/:id", productHandler.UpdateProduct)
				products.DELETE("/:id", productHandler.DeleteProduct)
				products.GET("/:id/inventory", productHandler.GetInventoryByProduct)
			}

			// User management routes (Admin only)
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetAllUsers)
				users.GET("/:id", userHandler.GetUserByID)
				users.POST("", userHandler.CreateUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.GET("/:id/warehouses", userHandler.GetUserWarehouses)
				users.POST("/assign-warehouse", userHandler.AssignWarehouse)
			}

			// Use separate group to avoid wildcard conflict with /users/:id
			userWarehouses := protected.Group("/user-warehouses")
			{
				userWarehouses.DELETE("/:user_id/warehouses/:warehouse_id", userHandler.RemoveWarehouse)
			}
		}
	}

	return router
}
