package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserIDKey is the context key for storing user ID
	UserIDKey ContextKey = "user_id"
	// PathParamsKey is the context key for storing path parameters
	PathParamsKey ContextKey = "path_params"
)

// JSONResponse sends a JSON response with the given status code
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// ErrorResponse sends a JSON error response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"error": message})
}

// DecodeJSON decodes JSON from request body
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// GetUserID retrieves user ID from request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}

// SetUserID sets user ID in context
func SetUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}

// GetPathParam retrieves a path parameter from context
func GetPathParam(r *http.Request, key string) string {
	params, ok := r.Context().Value(PathParamsKey).(map[string]string)
	if !ok {
		return ""
	}
	return params[key]
}

// SetPathParams sets path parameters in context
func SetPathParams(r *http.Request, params map[string]string) *http.Request {
	ctx := context.WithValue(r.Context(), PathParamsKey, params)
	return r.WithContext(ctx)
}

// SetCookie sets an HTTP cookie
func SetCookie(w http.ResponseWriter, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(w, cookie)
}

