package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	apperrors "kasir-api/errors"
	"kasir-api/models"
	"kasir-api/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts
// GET /api/products
// POST /api/products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Product
// @Failure      500  {string}  string  "Failed to fetch products"
// @Router       /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		log.Println("Error fetching products:", err)
		WriteError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	WriteJSON(w, http.StatusOK, products)
}

// Create godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided data
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.ProductInput  true  "Product data"
// @Success      201      {object}  models.Product
// @Failure      400      {string}  string  "Invalid request body"
// @Failure      400      {string}  string  "Category not found"
// @Router       /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(&product); err != nil {
		log.Println("Error creating product:", err)
		if errors.Is(err, apperrors.ErrCategoryNotFound) {
			WriteError(w, http.StatusBadRequest, "Category not found")
			return
		}
		WriteError(w, http.StatusBadRequest, "Failed to create product")
		return
	}

	WriteJSON(w, http.StatusCreated, product)
}

// HandleProductByID
// GET /api/products/{id}
// PUT /api/products/{id}
// DELETE /api/products/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetByID godoc
// @Summary      Get product by ID
// @Description  Retrieve a single product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {string}  string  "Invalid product ID"
// @Failure      404  {string}  string  "Product not found"
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/products/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		log.Println("Error fetching product by ID:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	WriteJSON(w, http.StatusOK, product)
}

// Update godoc
// @Summary      Update a product
// @Description  Update an existing product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Product ID"
// @Param        product  body      models.ProductInput  true  "Product data"
// @Success      200      {object}  models.Product
// @Failure      400      {string}  string  "Invalid product ID or request body"
// @Failure      400      {string}  string  "Category not found"
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/products/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product.ID = id
	if err := h.service.Update(&product); err != nil {
		log.Println("Error updating product:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		if errors.Is(err, apperrors.ErrCategoryNotFound) {
			WriteError(w, http.StatusBadRequest, "Category not found")
			return
		}
		WriteError(w, http.StatusBadRequest, "Failed to update product")
		return
	}

	WriteJSON(w, http.StatusOK, product)
}

// Delete godoc
// @Summary      Delete a product
// @Description  Delete a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  handlers.APIResponse  "Product deleted successfully"
// @Failure      400  {string}  string  "Invalid product ID"
// @Failure      404  {string}  string  "Product not found"
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/products/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Println("Error deleting product:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		WriteError(w, http.StatusBadRequest, "Failed to delete product")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
