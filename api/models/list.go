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
	Checked  bool               `json:"checked" bson:"checked"`
	Details  string             `json:"details,omitempty" bson:"details,omitempty"`
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

// AddListItemRequest represents the request body for adding an item to a list
type AddListItemRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity"`
	Details  string `json:"details,omitempty" binding:"max=512"`
}

// UpdateListItemCheckedRequest represents the request body for updating an item's checked state
type UpdateListItemCheckedRequest struct {
	Index   *int `json:"index" binding:"required"`
	Checked bool `json:"checked"`
}

// UpdateListItemRequest represents the request body for updating an item's name, details, and quantity
type UpdateListItemRequest struct {
	Index    *int    `json:"index" binding:"required"`
	Name     string  `json:"name,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
	Details  *string `json:"details,omitempty"`
}

// DeleteListItemRequest represents the request body for deleting an item from a list
type DeleteListItemRequest struct {
	Index *int `json:"index" binding:"required"`
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

