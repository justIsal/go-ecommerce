package main

import (
	"go-ecommerce/config"
	"go-ecommerce/internal/domain"
	"go-ecommerce/internal/handler"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/repository"
	"go-ecommerce/internal/service"
	"log"
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	err := db.AutoMigrate(
		&domain.User{}, 
		&domain.RefreshToken{}, 
		&domain.Product{},
		&domain.Order{},     
		&domain.OrderItem{}, 
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db) 

	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo) 

	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService) 
	uploadHandler := handler.NewUploadHandler()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, 
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")
	api := r.Group("/api/v1")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.RefreshToken)
	api.GET("/products", productHandler.GetAll)
	api.GET("/products/:id", productHandler.GetByID)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/auth/logout", authHandler.Logout)
		
		protected.POST("/checkout", orderHandler.Checkout)

		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware(domain.RoleAdmin))
		{
			admin.POST("/upload", uploadHandler.UploadImage)
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)
		}
	}

	log.Println("Server running on port 8080")
	r.Run(":8080")
}