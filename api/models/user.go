package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document in MongoDB
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	Profile      *Profile           `json:"profile,omitempty" bson:"profile,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// Profile represents user profile information
type Profile struct {
	FirstName string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
}

// SignupRequest represents the request body for signup
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SigninRequest represents the request body for signin
type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response for signup/signin
type AuthResponse struct {
	Token string      `json:"token"`
	User  *UserPublic `json:"user"`
}

// UserPublic represents public user information (without password)
type UserPublic struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

