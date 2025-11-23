package main

import (
	"log"
	"simple-inventory/internal/infrastructure/config"
	"simple-inventory/internal/infrastructure/database"
	"simple-inventory/internal/infrastructure/persistence"
	"simple-inventory/internal/interfaces/http/handler"
	"simple-inventory/internal/interfaces/http/router"
	"simple-inventory/internal/usecase"
	"time"
)

func main() {
	time.Local = time.UTC
	log.Println("Starting Simple Inventory API...")
	log.Println("Timezone set to UTC")

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := persistence.NewUserRepository(db)
	sessionRepo := persistence.NewSessionRepository(db)
	productRepo := persistence.NewProductRepository(db)
	locationRepo := persistence.NewLocationRepository(db)
	transactionRepo := persistence.NewInventoryTransactionRepository(db)

	authUseCase := usecase.NewAuthUseCase(userRepo, sessionRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	productUseCase := usecase.NewProductUseCase(productRepo)
	locationUseCase := usecase.NewLocationUseCase(locationRepo)
	inventoryUseCase := usecase.NewInventoryUseCase(productRepo, transactionRepo)

	authHandler := handler.NewAuthHandler(authUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	inventoryHandler := handler.NewInventoryHandler(inventoryUseCase)

	r := router.NewRouter(authHandler, productHandler, locationHandler, inventoryHandler, authUseCase)
	engine := r.Setup()

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := engine.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
