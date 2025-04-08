package main

import (
	"log"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/framework/driver/db"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/interface_adapter/controller"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/interface_adapter/gateway"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/interface_adapter/routes"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/service"
	"example.com/EVENT-MANAGEMENT-SYSTEM/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file (it will look for .env in the root directory)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}
}

func main() {
	// Load the database configuration from environment variables or .env
	dbConfig := config.LoadDBConfig()

	// Debug: Print the loaded database configuration
	log.Printf("DB Config: Host=%s, Port=%s, User=%s, Password=%s, DBName=%s, SSLMode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	// Connect to the database
	database, err := db.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer database.Close()

	// Create the required tables if they don't exist
	if err := db.CreateTables(database); err != nil {
		log.Fatal("Error creating tables:", err)

	}

	// Initialize the repositories
	userRepository := gateway.NewUserRepository(database)
	tokenRepository := gateway.NewTokenRepository(database)

	// Initialize the services
	userService := service.NewUserService(userRepository, tokenRepository)

	// Initialize the controllers
	userController := controller.NewUserController(userService)

	r := gin.Default()
	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend-domain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterUserRoutes(r, userController, tokenRepository)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
