package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// List represents a list document in MongoDB
type List struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	Items       []ListItem           `json:"items" bson:"items"`
	SharedWith  []primitive.ObjectID `json:"shared_with" bson:"shared_with"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

// ListItem represents an item in a list
type ListItem struct {
	Name     string             `json:"name" bson:"name"`
	Quantity int                `json:"quantity" bson:"quantity"`
	Unit     string             `json:"unit,omitempty" bson:"unit,omitempty"`
	Checked  bool               `json:"checked" bson:"checked"`
	AddedBy  primitive.ObjectID `json:"added_by" bson:"added_by"`
	AddedAt  time.Time          `json:"added_at" bson:"added_at"`
}

// CreateListRequest represents the request body for creating a list
type CreateListRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// UpdateListRequest represents the request body for updating a list
type UpdateListRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListResponse represents the response for list operations
type ListResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Items       []ListItem `json:"items"`
	SharedWith  []string   `json:"shared_with"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

