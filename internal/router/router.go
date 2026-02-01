package router

import (
	"net/http"

	"kasir-api/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

// New creates and configures the HTTP router with all routes
func New(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler.Health)

	// Category routes
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Product routes
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}
