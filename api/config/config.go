package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	JWTSecret   string
	MongoClient *mongo.Client
	DB          *mongo.Database
)

func Init() {
	// Load .env file from project root (ignore error if it doesn't exist)
	_ = godotenv.Load("../.env")

	// Load JWT secret
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set. Please create a .env file or set the environment variable.")
	}

	// MongoDB client should be set by main.go after connection
}

func SetMongoClient(client *mongo.Client) {
	MongoClient = client
	DB = client.Database("grocer-me")
}

func GetMongoURI() string {
	// Try to load .env file from project root (ignore error if it doesn't exist)
	_ = godotenv.Load("../.env")

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set. Please create a .env file or set the environment variable.")
	}
	return uri
}

