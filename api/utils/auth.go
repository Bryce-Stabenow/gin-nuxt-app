package utils

import (
	"context"
	"net/http"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetAuthenticatedUser retrieves the authenticated user ID from context and validates it
func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, bool) {
	userIDStr, ok := GetUserID(r)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User ID not found in context")
		return primitive.ObjectID{}, false
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return primitive.ObjectID{}, false
	}

	return userID, true
}

// GetAndValidateListID extracts and validates the list ID from path parameters
func GetAndValidateListID(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, bool) {
	listIDStr := GetPathParam(r, "id")
	if listIDStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "List ID is required")
		return primitive.ObjectID{}, false
	}

	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid list ID format")
		return primitive.ObjectID{}, false
	}

	return listID, true
}

// FetchList retrieves a list by ID from the database
func FetchList(w http.ResponseWriter, listID primitive.ObjectID) (*models.List, bool) {
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var list models.List
	err := collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ErrorResponse(w, http.StatusNotFound, "List not found")
			return nil, false
		}
		ErrorResponse(w, http.StatusInternalServerError, "Failed to find list")
		return nil, false
	}

	return &list, true
}

// CheckListAccess verifies if a user has access to a list (owner or shared with)
func CheckListAccess(w http.ResponseWriter, list *models.List, userID primitive.ObjectID) bool {
	if list.UserID == userID {
		return true
	}

	for _, sharedUserID := range list.SharedWith {
		if sharedUserID == userID {
			return true
		}
	}

	ErrorResponse(w, http.StatusForbidden, "You do not have access to this list")
	return false
}

// CheckListOwnership verifies if a user is the owner of a list
func CheckListOwnership(w http.ResponseWriter, list *models.List, userID primitive.ObjectID) bool {
	if list.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to perform this action")
		return false
	}
	return true
}

