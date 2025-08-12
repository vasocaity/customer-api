package main

import (
	"customer-api/pkg/db"
	"customer-api/pkg/handler"
	"customer-api/pkg/model"
	"customer-api/pkg/repository"
	"customer-api/pkg/service"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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
	if err := database.AutoMigrate(&model.Customer{},
		&model.Product{},
		&model.Feedback{},
		&model.Interaction{},
	); err != nil {
		log.Fatalf("Migrate failed: %v", err)
	}

	// inject dependencies
	cusRepo := repository.NewRepository(database)
	cusService := service.NewService(cusRepo)
	productRepository := repository.NewProductRepository(database)
	cusHandler := handler.NewCustomerHandler(cusService, productRepository)
	feedbackHandler := handler.NewFeedbackHandler(database)

	// Middleware
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	customer := r.Group("customers")
	customer.GET("", cusHandler.Get)
	customer.POST("", cusHandler.CreateCustomer)
	customer.DELETE("/:id", cusHandler.DeleteByID)
	customer.PUT("/:id", cusHandler.UpdateByID)
	customer.GET("/:id", cusHandler.GetByID)

	feedbackGroup := r.Group("/feedbacks")
	{
		feedbackGroup.POST("", feedbackHandler.CreateFeedback)
		feedbackGroup.GET("", feedbackHandler.ListFeedbacks)
		feedbackGroup.GET("/:id", feedbackHandler.GetFeedback)
		feedbackGroup.PUT("/:id", feedbackHandler.UpdateFeedback)
		feedbackGroup.DELETE("/:id", feedbackHandler.DeleteFeedback)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
