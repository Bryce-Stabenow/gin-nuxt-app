package utils

import (
	"net/http"
	"strings"
)

// Route represents a single route with its handler
type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// Router handles HTTP routing
type Router struct {
	routes      []Route
	middlewares []func(http.HandlerFunc) http.HandlerFunc
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes:      []Route{},
		middlewares: []func(http.HandlerFunc) http.HandlerFunc{},
	}
}

// Use adds middleware to the router
func (router *Router) Use(middleware func(http.HandlerFunc) http.HandlerFunc) {
	router.middlewares = append(router.middlewares, middleware)
}

// AddRoute adds a route to the router
func (router *Router) AddRoute(method, pattern string, handler http.HandlerFunc) {
	router.routes = append(router.routes, Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
}

// GET adds a GET route
func (router *Router) GET(pattern string, handler http.HandlerFunc) {
	router.AddRoute("GET", pattern, handler)
}

// POST adds a POST route
func (router *Router) POST(pattern string, handler http.HandlerFunc) {
	router.AddRoute("POST", pattern, handler)
}

// PUT adds a PUT route
func (router *Router) PUT(pattern string, handler http.HandlerFunc) {
	router.AddRoute("PUT", pattern, handler)
}

// DELETE adds a DELETE route
func (router *Router) DELETE(pattern string, handler http.HandlerFunc) {
	router.AddRoute("DELETE", pattern, handler)
}

// ServeHTTP implements the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a handler that will process the request
	var finalHandler http.HandlerFunc
	
	// Find matching route
	routeFound := false
	for _, route := range router.routes {
		if route.Method != r.Method {
			continue
		}

		params, matches := matchPattern(route.Pattern, r.URL.Path)
		if matches {
			routeFound = true
			
			// Apply path parameters to request context
			if len(params) > 0 {
				r = SetPathParams(r, params)
			}

			finalHandler = route.Handler
			break
		}
	}

	// If no route found, use 404 handler
	if !routeFound {
		finalHandler = http.NotFound
	}

	// Apply global middlewares to the handler
	for i := len(router.middlewares) - 1; i >= 0; i-- {
		finalHandler = router.middlewares[i](finalHandler)
	}

	// Execute the final handler
	finalHandler(w, r)
}

// matchPattern matches a URL pattern against a path and extracts parameters
// Pattern format: "/lists/:id" matches "/lists/123" with params["id"] = "123"
func matchPattern(pattern, path string) (map[string]string, bool) {
	params := make(map[string]string)
	
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Handle root path
	if pattern == "/" && path == "/" {
		return params, true
	}

	// Different number of parts means no match
	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	for i := 0; i < len(patternParts); i++ {
		patternPart := patternParts[i]
		pathPart := pathParts[i]

		if strings.HasPrefix(patternPart, ":") {
			// This is a parameter
			paramName := strings.TrimPrefix(patternPart, ":")
			params[paramName] = pathPart
		} else if patternPart != pathPart {
			// Not a parameter and doesn't match
			return nil, false
		}
	}

	return params, true
}

