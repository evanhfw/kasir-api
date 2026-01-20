package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Constants
const (
	serverPort = ":8080"
	apiPrefix  = "/api/product"
)

// Structs
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

// Global variables
var (
	products = []Product{
		{ID: 1, Name: "Laptop", Price: 15000000, Stock: 10},
		{ID: 2, Name: "Smartphone", Price: 5000000, Stock: 25},
		{ID: 3, Name: "Tablet", Price: 3000000, Stock: 15},
	}
	nextID = 4 // ID counter untuk product baru
)

// Helper functions
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

// Handler functions

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)

	respondJSON(w, http.StatusCreated, newProduct)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	for _, product := range products {
		if product.ID == id {
			respondJSON(w, http.StatusOK, product)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Product not found")
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for i, product := range products {
		if product.ID == id {
			products[i].Name = updatedProduct.Name
			products[i].Price = updatedProduct.Price
			products[i].Stock = updatedProduct.Stock

			respondJSON(w, http.StatusOK, products[i])
			return
		}
	}

	respondError(w, http.StatusNotFound, "Product not found")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			respondJSON(w, http.StatusOK, SuccessResponse{
				Message: "Product deleted successfully",
			})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Product not found")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

// Main function
func main() {
	// Route handler untuk /api/product/{id}
	http.HandleFunc(apiPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductByID(w, r)
		case http.MethodPut:
			updateProduct(w, r)
		case http.MethodDelete:
			deleteProduct(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Route handler untuk /api/product
	http.HandleFunc(apiPrefix, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllProducts(w, r)
		case http.MethodPost:
			createProduct(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", healthCheck)

	fmt.Println("Server running di localhost" + serverPort)
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		fmt.Println("Gagal running server:", err)
	}
}
