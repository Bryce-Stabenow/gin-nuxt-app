package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bryce-stabenow/grocer-me/models"
	"bryce-stabenow/grocer-me/testutil"
	"bryce-stabenow/grocer-me/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note: These tests demonstrate the expected behavior of list handlers
// For full integration tests, you would need a test MongoDB instance

func TestHandleCreateList_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	req := models.CreateListRequest{
		Name:        "Grocery List",
		Description: "Weekly groceries",
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "POST", "/lists", req)
	httpReq = utils.SetUserID(httpReq, userID)
	w := httptest.NewRecorder()

	HandleCreateList(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Description, response.Description)
	assert.Empty(t, response.Items)
}

func TestHandleCreateList_Unauthenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	req := models.CreateListRequest{
		Name: "Test List",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/lists", req, "")
	w := httptest.NewRecorder()

	HandleCreateList(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestHandleCreateList_InvalidJSON(t *testing.T) {
	testutil.SetupTestConfig(t)

	userID := "507f1f77bcf86cd799439011"
	httpReq := testutil.CreateAuthenticatedRequest(t, "POST", "/lists", nil)
	httpReq = utils.SetUserID(httpReq, userID)
	w := httptest.NewRecorder()

	HandleCreateList(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleGetLists_Authenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/lists", nil)
	httpReq = utils.SetUserID(httpReq, userID)
	w := httptest.NewRecorder()

	HandleGetLists(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.NotNil(t, response)
}

func TestHandleGetLists_Unauthenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	httpReq := httptest.NewRequest("GET", "/lists", nil)
	w := httptest.NewRecorder()

	HandleGetLists(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestHandleGetList_ValidID(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/lists/"+listID, nil)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleGetList(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.Equal(t, listID, response.ID)
}

func TestHandleGetList_InvalidID(t *testing.T) {
	testutil.SetupTestConfig(t)

	userID := "507f1f77bcf86cd799439011"

	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/lists/invalid", nil)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": "invalid"})
	w := httptest.NewRecorder()

	HandleGetList(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleGetList_NotFound(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	nonExistentListID := "507f1f77bcf86cd799439999"

	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/lists/"+nonExistentListID, nil)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": nonExistentListID})
	w := httptest.NewRecorder()

	HandleGetList(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusNotFound)
}

func TestHandleUpdateList_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	req := models.UpdateListRequest{
		Name:        "Updated List Name",
		Description: "Updated description",
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID, req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateList(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Description, response.Description)
}

func TestHandleUpdateList_Forbidden(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// User trying to update another user's list
	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	req := models.UpdateListRequest{
		Name: "Hacked List",
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID, req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateList(w, httpReq)

	assert.Equal(t, http.StatusForbidden, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusForbidden)
}

func TestHandleDeleteList_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	httpReq := testutil.CreateAuthenticatedRequest(t, "DELETE", "/lists/"+listID, nil)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleDeleteList(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	testutil.ParseJSONResponse(t, w, &response)

	assert.Equal(t, "List deleted successfully", response["message"])
}

func TestHandleDeleteList_NotOwner(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// User trying to delete another user's list
	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	httpReq := testutil.CreateAuthenticatedRequest(t, "DELETE", "/lists/"+listID, nil)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleDeleteList(w, httpReq)

	assert.Equal(t, http.StatusForbidden, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusForbidden)
}

func TestHandleAddListItem_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	req := models.AddListItemRequest{
		Name:     "Milk",
		Quantity: 2,
		Details:  "2% organic",
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "POST", "/lists/"+listID+"/items", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleAddListItem(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.NotEmpty(t, response.Items)
	// The new item should be in the list
	found := false
	for _, item := range response.Items {
		if item.Name == req.Name {
			found = true
			assert.Equal(t, req.Quantity, item.Quantity)
			assert.Equal(t, req.Details, item.Details)
			break
		}
	}
	assert.True(t, found, "New item should be in the list")
}

func TestHandleAddListItem_DefaultQuantity(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	req := models.AddListItemRequest{
		Name: "Bread",
		// Quantity not specified or 0
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "POST", "/lists/"+listID+"/items", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleAddListItem(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	// The new item should have quantity 1 by default
	found := false
	for _, item := range response.Items {
		if item.Name == req.Name {
			found = true
			assert.Equal(t, 1, item.Quantity)
			break
		}
	}
	assert.True(t, found)
}

func TestHandleUpdateListItemChecked_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	index := 0
	req := models.UpdateListItemCheckedRequest{
		Index:   &index,
		Checked: true,
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID+"/items/checked", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateListItemChecked(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.True(t, response.Items[0].Checked)
}

func TestHandleUpdateListItemChecked_InvalidIndex(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	index := 999 // Invalid index
	req := models.UpdateListItemCheckedRequest{
		Index:   &index,
		Checked: true,
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID+"/items/checked", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateListItemChecked(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleUpdateListItem_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	index := 0
	quantity := 5
	details := "Updated details"
	req := models.UpdateListItemRequest{
		Index:    &index,
		Name:     "Updated Item Name",
		Quantity: &quantity,
		Details:  &details,
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID+"/items", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateListItem(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.Equal(t, req.Name, response.Items[0].Name)
	assert.Equal(t, *req.Quantity, response.Items[0].Quantity)
	assert.Equal(t, *req.Details, response.Items[0].Details)
}

func TestHandleUpdateListItem_DetailsMaxLength(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	index := 0
	longDetails := make([]byte, 513) // More than 512 characters
	for i := range longDetails {
		longDetails[i] = 'a'
	}
	detailsStr := string(longDetails)

	req := models.UpdateListItemRequest{
		Index:   &index,
		Details: &detailsStr,
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "PUT", "/lists/"+listID+"/items", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleUpdateListItem(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleDeleteListItem_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"

	index := 0
	req := models.DeleteListItemRequest{
		Index: &index,
	}

	httpReq := testutil.CreateAuthenticatedRequest(t, "DELETE", "/lists/"+listID+"/items", req)
	httpReq = utils.SetUserID(httpReq, userID)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleDeleteListItem(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	// Item should be removed
	assert.NotNil(t, response.Items)
}

func TestHandleShareList_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := primitive.NewObjectID()
	listID := "507f1f77bcf86cd799439012"
	
	token, err := testutil.GenerateTestToken(userID.Hex(), "test-secret-key-for-testing-purposes-only")
	assert.NoError(t, err)

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/lists/share/"+listID, nil, token)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleShareList(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ListResponse
	testutil.ParseJSONResponse(t, w, &response)

	// User should be added to shared_with
	found := false
	for _, sharedUser := range response.SharedWith {
		if sharedUser.ID == userID.Hex() {
			found = true
			break
		}
	}
	assert.True(t, found, "User should be in shared_with array")
}

func TestHandleShareList_AlreadyOwner(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// User trying to share a list they own
	userID := "507f1f77bcf86cd799439011"
	listID := "507f1f77bcf86cd799439012"
	
	token, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	assert.NoError(t, err)

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/lists/share/"+listID, nil, token)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleShareList(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleShareList_Unauthenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	listID := "507f1f77bcf86cd799439012"

	httpReq := httptest.NewRequest("POST", "/lists/share/"+listID, nil)
	httpReq = utils.SetPathParams(httpReq, map[string]string{"id": listID})
	w := httptest.NewRecorder()

	HandleShareList(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestListToResponse(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// Test the listToResponse helper function
	userID := primitive.NewObjectID()
	sharedUserID := primitive.NewObjectID()

	list := &models.List{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        "Test List",
		Description: "Test Description",
		Items: []models.ListItem{
			{
				Name:     "Item 1",
				Quantity: 1,
				Checked:  false,
			},
		},
		SharedWith: []primitive.ObjectID{sharedUserID},
	}

	response := listToResponse(list)

	assert.Equal(t, list.ID.Hex(), response.ID)
	assert.Equal(t, list.UserID.Hex(), response.UserID)
	assert.Equal(t, list.Name, response.Name)
	assert.Equal(t, list.Description, response.Description)
	assert.Len(t, response.Items, 1)
	assert.Len(t, response.SharedWith, 1)
}
