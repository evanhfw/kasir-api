package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Constants
const (
	serverPort = ":8080"
	apiPrefix  = "/api"
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

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

var (
	categories = []Category{
		{ID: 1, Name: "Electronics", Description: "Electronic devices and gadgets"},
		{ID: 2, Name: "Furniture", Description: "Home and office furniture"},
	}
	nextCategoryID = 3 // ID counter untuk category baru
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
// getAllCategories godoc
// @Summary      Mendapatkan semua kategori
// @Description  Mengembalikan list array dari semua kategori yang tersedia
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   Category
// @Router       /categories [get]
func getAllCategories(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, categories)
}

// getAllProducts godoc
// @Summary      Mendapatkan semua produk
// @Description  Mengembalikan list array dari semua produk yang tersedia
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   Product
// @Router       /product [get]
func getAllProducts(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, products)
}

// createCategory godoc
// @Summary      Membuat kategori baru
// @Description  Menambahkan data kategori baru ke database
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category body Category true "Data kategori"
// @Success      201  {object}  Category
// @Failure      400  {object}  ErrorResponse
// @Router       /categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newCategory.ID = nextCategoryID
	nextCategoryID++
	categories = append(categories, newCategory)

	respondJSON(w, http.StatusCreated, newCategory)
}

// createProduct godoc
// @Summary      Membuat produk baru
// @Description  Menambahkan data produk baru ke database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product body Product true "Data Produk"
// @Success      201  {object}  Product
// @Failure      400  {object}  ErrorResponse
// @Router       /product [post]
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

// getProductByID godoc
// @Summary      Mendapatkan produk berdasarkan ID
// @Description  Mengambil detail satu produk
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Product
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /product/{id} [get]
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/product/")
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

// updateProduct godoc
// @Summary      Mengupdate produk
// @Description  Mengupdate data produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param        product body Product true "Data Produk"
// @Success      200  {object}  Product
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /product/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/product/")
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

// deleteProduct godoc
// @Summary      Menghapus produk
// @Description  Menghapus produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /product/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/product/")
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

// healthCheck godoc
// @Summary      Health check endpoint
// @Description  Mengecek status API
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

// getCategoryByID godoc
// @Summary      Mendapatkan kategori berdasarkan ID
// @Description  Mengambil detail satu kategori
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  Category
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	for _, category := range categories {
		if category.ID == id {
			respondJSON(w, http.StatusOK, category)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Category not found")
}

// updateCategory godoc
// @Summary      Mengupdate kategori
// @Description  Mengupdate data kategori berdasarkan ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Param        category body Category true "Data kategori"
// @Success      200  {object}  Category
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for i, category := range categories {
		if category.ID == id {
			categories[i].Name = updatedCategory.Name
			categories[i].Description = updatedCategory.Description

			respondJSON(w, http.StatusOK, categories[i])
			return
		}
	}

	respondError(w, http.StatusNotFound, "Category not found")
}

// deleteCategory godoc
// @Summary      Menghapus kategori
// @Description  Menghapus kategori berdasarkan ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, apiPrefix+"/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			respondJSON(w, http.StatusOK, SuccessResponse{
				Message: "Category deleted successfully",
			})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Category not found")
}

// Main function
// @title           Kasir API
// @version         1.0
// @description     API sederhana untuk manajemen kasir (Produk & Kategori).
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      kasir-api.evan-homeserver.my.id
// @BasePath  /api
// @schemes   https http

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// GET /api/product/{id}
	// PUT /api/product/{id}
	// DELETE /api/product/{id}
	http.HandleFunc(apiPrefix+"/product/", func(w http.ResponseWriter, r *http.Request) {
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

	// GET /api/product
	// POST /api/product
	http.HandleFunc(apiPrefix+"/product", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllProducts(w, r)
		case http.MethodPost:
			createProduct(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// GET /api/categories
	// POST /api/categories
	http.HandleFunc(apiPrefix+"/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllCategories(w, r)
		case http.MethodPost:
			createCategory(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// GET /api/categories/{id}
	// PUT /api/categories/{id}
	// DELETE /api/categories/{id}
	http.HandleFunc(apiPrefix+"/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategoryByID(w, r)
		case http.MethodPut:
			updateCategory(w, r)
		case http.MethodDelete:
			deleteCategory(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// GET /health
	http.HandleFunc("/health", healthCheck)

	http.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	fmt.Println("Server running di localhost" + serverPort)
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		fmt.Println("Gagal running server:", err)
	}
}
