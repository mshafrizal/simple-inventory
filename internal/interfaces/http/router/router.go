package router

import (
	"simple-inventory/internal/interfaces/http/handler"
	"simple-inventory/internal/interfaces/http/middleware"
	"simple-inventory/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	authHandler      *handler.AuthHandler
	productHandler   *handler.ProductHandler
	locationHandler  *handler.LocationHandler
	inventoryHandler *handler.InventoryHandler
	authUseCase      *usecase.AuthUseCase
}

func NewRouter(
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	locationHandler *handler.LocationHandler,
	inventoryHandler *handler.InventoryHandler,
	authUseCase *usecase.AuthUseCase,
) *Router {
	return &Router{
		engine:           gin.Default(),
		authHandler:      authHandler,
		productHandler:   productHandler,
		locationHandler:  locationHandler,
		inventoryHandler: inventoryHandler,
		authUseCase:      authUseCase,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.RecoveryMiddleware())
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggingMiddleware())
	r.engine.Use(middleware.ErrorMiddleware())

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.engine.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/logout", r.authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(r.authUseCase), r.authHandler.GetCurrentUser)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.authUseCase))
		{
			products := protected.Group("/products")
			{
				products.POST("", r.productHandler.CreateProduct)
				products.GET("", r.productHandler.ListProducts)
				products.GET("/search", r.productHandler.SearchProducts)
				products.GET("/low-stock", r.productHandler.GetLowStockProducts)
				products.GET("/scan", r.productHandler.GetProductByBarcode)
				products.GET("/:id", r.productHandler.GetProduct)
				products.PUT("/:id", r.productHandler.UpdateProduct)
				products.DELETE("/:id", r.productHandler.DeleteProduct)
			}

			locations := protected.Group("/locations")
			{
				locations.POST("", r.locationHandler.CreateLocation)
				locations.GET("", r.locationHandler.ListLocations)
				locations.GET("/search", r.locationHandler.SearchLocations)
				locations.GET("/:id", r.locationHandler.GetLocation)
				locations.PUT("/:id", r.locationHandler.UpdateLocation)
				locations.DELETE("/:id", r.locationHandler.DeleteLocation)
			}

			inventory := protected.Group("/inventory")
			{
				inventory.POST("/receive", r.inventoryHandler.ReceiveInventory)
				inventory.POST("/issue", r.inventoryHandler.IssueInventory)
				inventory.POST("/adjust", r.inventoryHandler.AdjustInventory)
				inventory.POST("/transfer", r.inventoryHandler.TransferInventory)
				inventory.GET("/transactions/product/:product_id", r.inventoryHandler.GetProductTransactions)
			}
		}
	}

	return r.engine
}
