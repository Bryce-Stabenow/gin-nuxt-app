package middleware

import (
	"net/http"
	"strings"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/utils"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates JWT tokens and extracts user ID
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// First, try to get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// If not in header, try to get from cookie
		if tokenString == "" {
			cookie, err := r.Cookie("jwt_token")
			if err == nil && cookie != nil {
				tokenString = cookie.Value
			}
		}

		// If still no token, return unauthorized
		if tokenString == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization required. Please sign in.")
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		// Store user ID in context
		r = utils.SetUserID(r, userID)
		next(w, r)
	}
}

// ExtractUserID extracts user ID from JWT token (used for public endpoints that optionally require auth)
func ExtractUserID(r *http.Request) (string, error) {
	var tokenString string

	// First, try to get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}

	// If not in header, try to get from cookie
	if tokenString == "" {
		cookie, err := r.Cookie("jwt_token")
		if err == nil && cookie != nil {
			tokenString = cookie.Value
		}
	}

	// If no token found, return error
	if tokenString == "" {
		return "", jwt.ErrSignatureInvalid
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	// Extract user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	return userID, nil
}
