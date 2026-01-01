package handlers

import (
	"context"
	"net/http"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/models"
	"bryce-stabenow/grocer-me/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// HandleSignup handles user registration
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		utils.ErrorResponse(w, http.StatusConflict, "Email already exists")
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to check email")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user with profile
	now := time.Now()
	profile := &models.Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if req.AvatarURL != nil && *req.AvatarURL != "" {
		profile.AvatarURL = *req.AvatarURL
	}
	
	user := models.User{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		Username:     req.Email,
		PasswordHash: string(hashedPassword),
		Profile:      profile,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID.Hex())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set JWT as HTTP-only cookie (24 hours expiration to match token)
	utils.SetCookie(w, "jwt_token", token, 3600*24, "/", "", false, true)

	// Return response
	utils.JSONResponse(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User: &models.UserPublic{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			Profile:   user.Profile,
			CreatedAt: user.CreatedAt,
		},
	})
}

// HandleSignin handles user login
func HandleSignin(w http.ResponseWriter, r *http.Request) {
	var req models.SigninRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by email
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to find user")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID.Hex())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set JWT as HTTP-only cookie (24 hours expiration to match token)
	utils.SetCookie(w, "jwt_token", token, 3600*24, "/", "", false, true)

	// Return response
	utils.JSONResponse(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User: &models.UserPublic{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			Profile:   user.Profile,
			CreatedAt: user.CreatedAt,
		},
	})
}

// HandleGetMe returns the current user's information
func HandleGetMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Find user by ID
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(w, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to find user")
		return
	}

	// Return user (password_hash is excluded via json:"-" tag)
	utils.JSONResponse(w, http.StatusOK, user)
}

// HandleLogout handles user logout by clearing the JWT cookie
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT cookie by setting it with an expired expiration time
	utils.SetCookie(w, "jwt_token", "", -1, "/", "", false, true)
	
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// generateToken creates a JWT token for the given user ID
func generateToken(userID string) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}
