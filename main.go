package main

import (
	"customer-api/pkg/db"
	"customer-api/pkg/handler"
	"customer-api/pkg/model"
	"customer-api/pkg/repository"
	"customer-api/pkg/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	database := db.NewPostgresDB()

	// Auto migrate
	if err := database.AutoMigrate(&model.Customer{}); err != nil {
		log.Fatalf("Migrate failed: %v", err)
	}

	// inject dependencies
	cusRepo := repository.NewRepository(database)
	cusService := service.NewService(cusRepo)
	cusHandler := handler.NewCustomerHandler(cusService)

	// Middleware
	r.Use(gin.Logger(), gin.Recovery())

	customer := r.Group("customers")
	customer.POST("", cusHandler.CreateCustomer)
	customer.GET("/:id", cusHandler.Get)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
