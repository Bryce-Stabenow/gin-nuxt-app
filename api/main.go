package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/handlers"
	"bryce-stabenow/grocer-me/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	// Initialize config (loads JWT_SECRET)
	config.Init()

	// Get MongoDB URI from environment variable
	mongoURI := config.GetMongoURI()

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	// Set MongoDB client in config
	config.SetMongoClient(client)

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes - HTML forms
	router.GET("/signup", handlers.HandleSignupForm)
	router.GET("/signin", handlers.HandleSigninForm)
	
	// Public routes - API endpoints
	router.POST("/signup", handlers.HandleSignup)
	router.POST("/signin", handlers.HandleSignin)

	// Protected routes (require JWT)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/me", handlers.HandleGetMe)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
