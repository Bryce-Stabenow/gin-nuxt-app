package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router)
	assert.Empty(t, router.routes)
	assert.Empty(t, router.middlewares)
}

func TestRouter_AddRoute(t *testing.T) {
	router := NewRouter()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.AddRoute("GET", "/test", handler)

	assert.Len(t, router.routes, 1)
	assert.Equal(t, "GET", router.routes[0].Method)
	assert.Equal(t, "/test", router.routes[0].Pattern)
}

func TestRouter_HTTPMethodHelpers(t *testing.T) {
	router := NewRouter()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	router.GET("/get", handler)
	router.POST("/post", handler)
	router.PUT("/put", handler)
	router.DELETE("/delete", handler)

	assert.Len(t, router.routes, 4)
	assert.Equal(t, "GET", router.routes[0].Method)
	assert.Equal(t, "POST", router.routes[1].Method)
	assert.Equal(t, "PUT", router.routes[2].Method)
	assert.Equal(t, "DELETE", router.routes[3].Method)
}

func TestRouter_Use(t *testing.T) {
	router := NewRouter()
	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "true")
			next(w, r)
		}
	}
	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "true")
			next(w, r)
		}
	}

	router.Use(middleware1)
	router.Use(middleware2)

	assert.Len(t, router.middlewares, 2)
}

func TestRouter_ServeHTTP_SimpleRoute(t *testing.T) {
	router := NewRouter()
	
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_ServeHTTP_RouteNotFound(t *testing.T) {
	router := NewRouter()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/test", handler)

	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouter_ServeHTTP_WithMiddleware(t *testing.T) {
	router := NewRouter()
	
	// Add middleware that sets a header
	middleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "applied")
			next(w, r)
		}
	}
	router.Use(middleware)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "applied", w.Header().Get("X-Middleware"))
}

func TestRouter_ServeHTTP_WithPathParams(t *testing.T) {
	router := NewRouter()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetPathParam(r, "id")
		action := GetPathParam(r, "action")
		
		assert.Equal(t, "123", id)
		assert.Equal(t, "edit", action)
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/items/:id/:action", handler)

	req := httptest.NewRequest("GET", "/items/123/edit", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_ServeHTTP_MethodMismatch(t *testing.T) {
	router := NewRouter()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/test", handler)

	// Try POST on a GET route
	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMatchPattern_ExactMatch(t *testing.T) {
	params, matches := matchPattern("/test", "/test")
	
	assert.True(t, matches)
	assert.Empty(t, params)
}

func TestMatchPattern_WithParameters(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		path        string
		shouldMatch bool
		expected    map[string]string
	}{
		{
			name:        "Single parameter",
			pattern:     "/users/:id",
			path:        "/users/123",
			shouldMatch: true,
			expected:    map[string]string{"id": "123"},
		},
		{
			name:        "Multiple parameters",
			pattern:     "/users/:id/posts/:postId",
			path:        "/users/123/posts/456",
			shouldMatch: true,
			expected:    map[string]string{"id": "123", "postId": "456"},
		},
		{
			name:        "No match - different length",
			pattern:     "/users/:id",
			path:        "/users/123/extra",
			shouldMatch: false,
			expected:    nil,
		},
		{
			name:        "No match - different static parts",
			pattern:     "/users/:id",
			path:        "/posts/123",
			shouldMatch: false,
			expected:    nil,
		},
		{
			name:        "Parameter at end",
			pattern:     "/api/v1/lists/:id",
			path:        "/api/v1/lists/abc123",
			shouldMatch: true,
			expected:    map[string]string{"id": "abc123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, matches := matchPattern(tt.pattern, tt.path)
			
			assert.Equal(t, tt.shouldMatch, matches)
			if tt.shouldMatch {
				assert.Equal(t, tt.expected, params)
			}
		})
	}
}

func TestMatchPattern_RootPath(t *testing.T) {
	params, matches := matchPattern("/", "/")
	
	assert.True(t, matches)
	assert.Empty(t, params)
}

func TestMatchPattern_TrailingSlashes(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		path        string
		shouldMatch bool
	}{
		{
			name:        "Both with trailing slash",
			pattern:     "/test/",
			path:        "/test/",
			shouldMatch: true,
		},
		{
			name:        "Pattern with trailing slash",
			pattern:     "/test/",
			path:        "/test",
			shouldMatch: true,
		},
		{
			name:        "Path with trailing slash",
			pattern:     "/test",
			path:        "/test/",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, matches := matchPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.shouldMatch, matches)
		})
	}
}

func TestRouter_MultipleMiddlewares(t *testing.T) {
	router := NewRouter()
	
	executionOrder := []string{}

	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1")
			next(w, r)
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2")
			next(w, r)
		}
	}

	router.Use(middleware1)
	router.Use(middleware2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})
	router.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Middlewares are applied in reverse order, so they execute in the order they were added
	assert.Equal(t, []string{"middleware1", "middleware2", "handler"}, executionOrder)
}
