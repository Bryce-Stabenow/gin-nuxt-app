package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/handlers"
	"bryce-stabenow/grocer-me/middleware"
	"bryce-stabenow/grocer-me/utils"

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

	// Initialize router
	router := utils.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(middleware.CORS)

	// Health check endpoint
	router.GET("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Public routes - API endpoints
	router.POST("/signup", handlers.HandleSignup)
	router.POST("/signin", handlers.HandleSignin)
	router.POST("/lists/share/:id", handlers.HandleShareList)

	// Protected routes (require JWT)
	router.GET("/me", withAuth(handlers.HandleGetMe))
	router.POST("/logout", withAuth(handlers.HandleLogout))

	// List routes
	router.POST("/lists", withAuth(handlers.HandleCreateList))
	router.GET("/lists", withAuth(handlers.HandleGetLists))
	router.GET("/lists/:id", withAuth(handlers.HandleGetList))
	router.PUT("/lists/:id", withAuth(handlers.HandleUpdateList))
	router.DELETE("/lists/:id", withAuth(handlers.HandleDeleteList))
	router.POST("/lists/:id/items", withAuth(handlers.HandleAddListItem))
	router.PUT("/lists/:id/items", withAuth(handlers.HandleUpdateListItem))
	router.DELETE("/lists/:id/items", withAuth(handlers.HandleDeleteListItem))
	router.PUT("/lists/:id/items/checked", withAuth(handlers.HandleUpdateListItemChecked))

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// withAuth wraps a handler with JWT authentication middleware
func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.JWTAuth(handler)
}
