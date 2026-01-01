package handlers

import (
	"context"
	"net/http"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/middleware"
	"bryce-stabenow/grocer-me/models"
	"bryce-stabenow/grocer-me/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// HandleCreateList handles creating a new list
func HandleCreateList(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.CreateListRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
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
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create list")
		return
	}

	// Fetch the created list to return
	var createdList models.List
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve created list")
		return
	}

	// Convert to response format
	response := listToResponse(&createdList)
	utils.JSONResponse(w, http.StatusCreated, response)
}

// HandleGetLists handles getting all lists for the authenticated user
func HandleGetLists(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
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
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch lists")
		return
	}
	defer cursor.Close(ctx)

	var lists []models.List
	if err = cursor.All(ctx, &lists); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to decode lists")
		return
	}

	// Convert to response format
	responses := make([]models.ListResponse, len(lists))
	for i, list := range lists {
		responses[i] = listToResponse(&list)
	}

	utils.JSONResponse(w, http.StatusOK, responses)
}

// HandleGetList handles getting a single list by ID
func HandleGetList(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Fetch list
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Convert to response format
	response := listToResponse(list)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleUpdateList handles updating a list
func HandleUpdateList(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.UpdateListRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch list and verify access
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Build update document
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{"$set": update},
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update list")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleAddListItem handles adding an item to a list
func HandleAddListItem(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.AddListItemRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set default quantity to 1 if not provided or 0
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	// Fetch list and verify access
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Create new item
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	newItem := models.ListItem{
		Name:     req.Name,
		Quantity: quantity,
		Checked:  false,
		Details:  req.Details,
		AddedBy:  userID,
		AddedAt:  now,
	}

	// Add item to list and update updated_at
	update := bson.M{
		"$push": bson.M{"items": newItem},
		"$set":  bson.M{"updated_at": now},
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		update,
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to add item to list")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleUpdateListItemChecked handles updating an item's checked state
func HandleUpdateListItemChecked(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.UpdateListItemCheckedRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch list and verify access
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Validate index
	if req.Index == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Index is required")
		return
	}

	index := *req.Index
	if index < 0 || index >= len(list.Items) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item index")
		return
	}

	// Update the item's checked state
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	
	// Update the item in the slice
	list.Items[index].Checked = req.Checked

	// Update the entire items array and updated_at in the database
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{
			"$set": bson.M{
				"items":      list.Items,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleUpdateListItem handles updating an item's name, details, and quantity
func HandleUpdateListItem(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.UpdateListItemRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch list and verify access
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Validate index
	if req.Index == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Index is required")
		return
	}

	index := *req.Index
	if index < 0 || index >= len(list.Items) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item index")
		return
	}

	// Validate details length if provided
	if req.Details != nil && len(*req.Details) > 512 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Details must be 512 characters or less")
		return
	}

	// Update the item's fields
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	
	// Update fields if provided
	if req.Name != "" {
		list.Items[index].Name = req.Name
	}
	if req.Quantity != nil && *req.Quantity > 0 {
		list.Items[index].Quantity = *req.Quantity
	}
	if req.Details != nil {
		// Allow empty string to clear the details field
		list.Items[index].Details = *req.Details
	}

	// Update the entire items array and updated_at in the database
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{
			"$set": bson.M{
				"items":      list.Items,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleDeleteListItem handles deleting an item from a list
func HandleDeleteListItem(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Parse request body
	var req models.DeleteListItemRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch list and verify access
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user has access
	if !utils.CheckListAccess(w, list, userID) {
		return // Error response already sent
	}

	// Validate index
	if req.Index == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Index is required")
		return
	}

	index := *req.Index
	if index < 0 || index >= len(list.Items) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item index")
		return
	}

	// Remove the item from the slice
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updatedItems := make([]models.ListItem, 0, len(list.Items)-1)
	updatedItems = append(updatedItems, list.Items[:index]...)
	updatedItems = append(updatedItems, list.Items[index+1:]...)

	// Update the items array and updated_at in the database
	now := time.Now()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{
			"$set": bson.M{
				"items":      updatedItems,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleDeleteList handles deleting a list
func HandleDeleteList(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := utils.GetAuthenticatedUser(w, r)
	if !ok {
		return // Error response already sent
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Fetch list and verify ownership
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Only the owner can delete the list
	if !utils.CheckListOwnership(w, list, userID) {
		return // Error response already sent
	}

	// Delete the list
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": listID})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete list")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "List deleted successfully"})
}

// HandleShareList handles adding the current user to a list's shared_with array
// This endpoint is public but requires authentication (checked internally)
func HandleShareList(w http.ResponseWriter, r *http.Request) {
	// Try to extract user ID from JWT (manual check for this public endpoint)
	userIDStr, err := middleware.ExtractUserID(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Authentication required. Please sign in to join this list.")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Get and validate list ID
	listID, ok := utils.GetAndValidateListID(w, r)
	if !ok {
		return // Error response already sent
	}

	// Fetch list
	list, ok := utils.FetchList(w, listID)
	if !ok {
		return // Error response already sent
	}

	// Check if user is already the owner
	if list.UserID == userID {
		utils.ErrorResponse(w, http.StatusBadRequest, "You are already the owner of this list")
		return
	}

	// Check if user is already in shared_with array
	alreadyShared := false
	for _, sharedUserID := range list.SharedWith {
		if sharedUserID == userID {
			alreadyShared = true
			break
		}
	}

	if alreadyShared {
		// User is already shared, return the list anyway (idempotent)
		response := listToResponse(list)
		utils.JSONResponse(w, http.StatusOK, response)
		return
	}

	// Add user to shared_with array
	collection := config.DB.Collection("lists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": listID},
		bson.M{
			"$addToSet": bson.M{"shared_with": userID},
			"$set":      bson.M{"updated_at": now},
		},
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to add user to shared list")
		return
	}

	// Fetch the updated list to return
	var updatedList models.List
	err = collection.FindOne(ctx, bson.M{"_id": listID}).Decode(&updatedList)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated list")
		return
	}

	// Convert to response format
	response := listToResponse(&updatedList)
	utils.JSONResponse(w, http.StatusOK, response)
}

// listToResponse converts a List model to ListResponse
func listToResponse(list *models.List) models.ListResponse {
	// Fetch user emails for shared_with users
	sharedWith := make([]models.SharedUser, 0, len(list.SharedWith))
	if len(list.SharedWith) > 0 {
		userCollection := config.DB.Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Fetch all users in a single query
		cursor, err := userCollection.Find(ctx, bson.M{"_id": bson.M{"$in": list.SharedWith}})
		if err == nil {
			defer cursor.Close(ctx)
			
			// Create a map of user ID to email for quick lookup
			userMap := make(map[primitive.ObjectID]string)
			var user models.User
			for cursor.Next(ctx) {
				if err := cursor.Decode(&user); err == nil {
					userMap[user.ID] = user.Email
				}
			}

			// Build sharedWith array maintaining the original order
			for _, userID := range list.SharedWith {
				if email, exists := userMap[userID]; exists {
					sharedWith = append(sharedWith, models.SharedUser{
						ID:    userID.Hex(),
						Email: email,
					})
				} else {
					// If user not found, still include the ID but with empty email
					sharedWith = append(sharedWith, models.SharedUser{
						ID:    userID.Hex(),
						Email: "",
					})
				}
			}
		} else {
			// If query fails, fall back to just IDs
			for _, userID := range list.SharedWith {
				sharedWith = append(sharedWith, models.SharedUser{
					ID:    userID.Hex(),
					Email: "",
				})
			}
		}
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
