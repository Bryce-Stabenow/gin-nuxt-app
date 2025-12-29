package handlers

import (
	"context"
	"net/http"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/middleware"
	"bryce-stabenow/grocer-me/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// HandleCreateList handles creating a new list
func HandleCreateList(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse request body
	var req models.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create list
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	list := models.List{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Items:       []models.ListItem{},
		SharedWith:  []primitive.ObjectID{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := collection.InsertOne(ctx, list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create list"})
		return
	}

	// Fetch the created list to return
	var createdList models.List
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created list"})
		return
	}

	// Convert to response format
	response := listToResponse(&createdList)
	c.JSON(http.StatusCreated, response)
}

// HandleGetLists handles getting all lists for the authenticated user
func HandleGetLists(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Find lists where user is owner or in shared_with array
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userID},
			{"shared_with": userID},
		},
	}

	// Sort by created_at descending
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lists"})
		return
	}
	defer cursor.Close(ctx)

	var lists []models.List
	if err = cursor.All(ctx, &lists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode lists"})
		return
	}

	// Convert to response format
	responses := make([]models.ListResponse, len(lists))
	for i, list := range lists {
		responses[i] = listToResponse(&list)
	}

	c.JSON(http.StatusOK, responses)
}

// HandleGetList handles getting a single list by ID
func HandleGetList(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get list ID from URL parameter
	listIDStr := c.Param("id")
	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID format"})
		return
	}

	// Find list
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var list models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find list"})
		return
	}

	// Check if user has access (owner or in shared_with)
	hasAccess := false
	if list.UserID == userID {
		hasAccess = true
	} else {
		for _, sharedUserID := range list.SharedWith {
			if sharedUserID == userID {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this list"})
		return
	}

	// Convert to response format
	response := listToResponse(&list)
	c.JSON(http.StatusOK, response)
}

// HandleUpdateList handles updating a list
func HandleUpdateList(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get list ID from URL parameter
	listIDStr := c.Param("id")
	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID format"})
		return
	}

	// Parse request body
	var req models.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find list and verify access
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var list models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find list"})
		return
	}

	// Check if user has access (only owner can update)
	if list.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this list"})
		return
	}

	// Build update document
	update := bson.M{
		"updated_at": time.Now(),
	}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}

	// Update the list
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{"$set": update},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update list"})
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated list"})
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	c.JSON(http.StatusOK, response)
}

// HandleDeleteList handles deleting a list (stub)
func HandleDeleteList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Delete list not yet implemented"})
}

// listToResponse converts a List model to ListResponse
func listToResponse(list *models.List) models.ListResponse {
	sharedWith := make([]string, len(list.SharedWith))
	for i, id := range list.SharedWith {
		sharedWith[i] = id.Hex()
	}

	return models.ListResponse{
		ID:          list.ID.Hex(),
		UserID:      list.UserID.Hex(),
		Name:        list.Name,
		Description: list.Description,
		Items:       list.Items,
		SharedWith:  sharedWith,
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
	}
}

