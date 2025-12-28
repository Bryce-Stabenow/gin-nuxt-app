package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	// Try to load .env file (ignore error if it doesn't exist)
	_ = godotenv.Load("../../../.env");

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set. Please create a .env file or set the environment variable.")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Get the database (you can change "grocer-me" to your preferred database name)
	db := client.Database("grocer-me")

	// Create User collection with indexes
	if err := createUserCollection(db); err != nil {
		log.Fatal("Error creating User collection:", err)
	}

	// Create List collection with indexes
	if err := createListCollection(db); err != nil {
		log.Fatal("Error creating List collection:", err)
	}

	fmt.Println("Successfully created User and List collections with indexes!")
}

func createUserCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("users")

	// Create indexes for User collection
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"email", 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true).SetName("username_unique"),
		},
		{
			Keys:    bson.D{{"created_at", 1}},
			Options: options.Index().SetName("created_at_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	fmt.Println("✓ User collection created with indexes (email, username, created_at)")

	// Create a sample document structure comment (optional - for documentation)
	// User document structure:
	// {
	//   "_id": ObjectId,
	//   "email": "user@example.com",
	//   "username": "username",
	//   "password_hash": "hashed_password",
	//   "profile": {
	//     "first_name": "John",
	//     "last_name": "Doe",
	//     "avatar_url": "https://..."
	//   },
	//   "created_at": ISODate,
	//   "updated_at": ISODate
	// }

	return nil
}

func createListCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("lists")

	// Create indexes for List collection
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"user_id", 1}},
			Options: options.Index().SetName("user_id_idx"),
		},
		{
			Keys:    bson.D{{"created_at", 1}},
			Options: options.Index().SetName("created_at_idx"),
		},
		{
			Keys:    bson.D{{"shared_with", 1}},
			Options: options.Index().SetName("shared_with_idx"),
		},
		{
			Keys:    bson.D{{"user_id", 1}, {"created_at", -1}},
			Options: options.Index().SetName("user_id_created_at_idx"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	fmt.Println("✓ List collection created with indexes (user_id, created_at, shared_with, user_id+created_at)")

	// Create a sample document structure comment (optional - for documentation)
	// List document structure:
	// {
	//   "_id": ObjectId,
	//   "user_id": ObjectId, // Reference to users collection
	//   "name": "Grocery List",
	//   "description": "Weekly shopping list",
	//   "items": [
	//     {
	//       "name": "Milk",
	//       "quantity": 1,
	//       "unit": "gallon",
	//       "checked": false,
	//       "added_by": ObjectId,
	//       "added_at": ISODate
	//     }
	//   ],
	//   "shared_with": [ObjectId], // Array of user IDs who have access
	//   "created_at": ISODate,
	//   "updated_at": ISODate
	// }

	return nil
}

