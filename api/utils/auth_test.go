package utils

import (
	"net/http/httptest"
	"testing"

	"bryce-stabenow/grocer-me/models"
	"bryce-stabenow/grocer-me/testutil"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAuthenticatedUser(t *testing.T) {
	t.Run("Valid user ID in context", func(t *testing.T) {
		userIDStr := "507f1f77bcf86cd799439011"
		
		req := testutil.CreateAuthenticatedRequest(t, "GET", "/test", nil)
		req = SetUserID(req, userIDStr)
		w := httptest.NewRecorder()

		userID, ok := GetAuthenticatedUser(w, req)

		assert.True(t, ok)
		assert.Equal(t, userIDStr, userID.Hex())
	})

	t.Run("No user ID in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		userID, ok := GetAuthenticatedUser(w, req)

		assert.False(t, ok)
		assert.Equal(t, primitive.ObjectID{}, userID)
		testutil.AssertErrorResponse(t, w, 401)
	})

	t.Run("Invalid user ID format", func(t *testing.T) {
		req := testutil.CreateAuthenticatedRequest(t, "GET", "/test", nil)
		req = SetUserID(req, "invalid-id")
		w := httptest.NewRecorder()

		userID, ok := GetAuthenticatedUser(w, req)

		assert.False(t, ok)
		assert.Equal(t, primitive.ObjectID{}, userID)
		testutil.AssertErrorResponse(t, w, 400)
	})
}

func TestGetAndValidateListID(t *testing.T) {
	t.Run("Valid list ID", func(t *testing.T) {
		listIDStr := "507f1f77bcf86cd799439011"
		
		req := testutil.CreateRequestWithPathParams(t, "GET", "/lists/"+listIDStr, nil)
		req = SetPathParams(req, map[string]string{"id": listIDStr})
		w := httptest.NewRecorder()

		listID, ok := GetAndValidateListID(w, req)

		assert.True(t, ok)
		assert.Equal(t, listIDStr, listID.Hex())
	})

	t.Run("Missing list ID", func(t *testing.T) {
		req := testutil.CreateRequestWithPathParams(t, "GET", "/lists/", nil)
		req = SetPathParams(req, map[string]string{})
		w := httptest.NewRecorder()

		listID, ok := GetAndValidateListID(w, req)

		assert.False(t, ok)
		assert.Equal(t, primitive.ObjectID{}, listID)
		testutil.AssertErrorResponse(t, w, 400)
	})

	t.Run("Invalid list ID format", func(t *testing.T) {
		req := testutil.CreateRequestWithPathParams(t, "GET", "/lists/invalid", nil)
		req = SetPathParams(req, map[string]string{"id": "invalid-id"})
		w := httptest.NewRecorder()

		listID, ok := GetAndValidateListID(w, req)

		assert.False(t, ok)
		assert.Equal(t, primitive.ObjectID{}, listID)
		testutil.AssertErrorResponse(t, w, 400)
	})
}

func TestCheckListAccess(t *testing.T) {
	userID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	sharedUserID := primitive.NewObjectID()

	t.Run("Owner has access", func(t *testing.T) {
		list := &models.List{
			ID:         primitive.NewObjectID(),
			UserID:     userID,
			SharedWith: []primitive.ObjectID{},
		}

		w := httptest.NewRecorder()
		hasAccess := CheckListAccess(w, list, userID)

		assert.True(t, hasAccess)
	})

	t.Run("Shared user has access", func(t *testing.T) {
		list := &models.List{
			ID:         primitive.NewObjectID(),
			UserID:     userID,
			SharedWith: []primitive.ObjectID{sharedUserID},
		}

		w := httptest.NewRecorder()
		hasAccess := CheckListAccess(w, list, sharedUserID)

		assert.True(t, hasAccess)
	})

	t.Run("Non-owner and non-shared user does not have access", func(t *testing.T) {
		list := &models.List{
			ID:         primitive.NewObjectID(),
			UserID:     userID,
			SharedWith: []primitive.ObjectID{},
		}

		w := httptest.NewRecorder()
		hasAccess := CheckListAccess(w, list, otherUserID)

		assert.False(t, hasAccess)
		testutil.AssertErrorResponse(t, w, 403)
	})

	t.Run("Multiple shared users", func(t *testing.T) {
		user1 := primitive.NewObjectID()
		user2 := primitive.NewObjectID()
		user3 := primitive.NewObjectID()

		list := &models.List{
			ID:         primitive.NewObjectID(),
			UserID:     userID,
			SharedWith: []primitive.ObjectID{user1, user2, user3},
		}

		// Each shared user should have access
		for _, uid := range []primitive.ObjectID{user1, user2, user3} {
			w := httptest.NewRecorder()
			hasAccess := CheckListAccess(w, list, uid)
			assert.True(t, hasAccess)
		}
	})
}

func TestCheckListOwnership(t *testing.T) {
	ownerID := primitive.NewObjectID()
	nonOwnerID := primitive.NewObjectID()

	t.Run("Owner check succeeds", func(t *testing.T) {
		list := &models.List{
			ID:     primitive.NewObjectID(),
			UserID: ownerID,
		}

		w := httptest.NewRecorder()
		isOwner := CheckListOwnership(w, list, ownerID)

		assert.True(t, isOwner)
	})

	t.Run("Non-owner check fails", func(t *testing.T) {
		list := &models.List{
			ID:     primitive.NewObjectID(),
			UserID: ownerID,
		}

		w := httptest.NewRecorder()
		isOwner := CheckListOwnership(w, list, nonOwnerID)

		assert.False(t, isOwner)
		testutil.AssertErrorResponse(t, w, 403)
	})
}
