package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"kasir-api/internal/apperrors"
	"kasir-api/internal/domain"
	"kasir-api/internal/service"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories handles GET and POST requests for /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
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
// @Summary      Get all categories
// @Description  Retrieve a list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Category
// @Failure      500  {string}  string  "Failed to fetch categories"
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		log.Println("Error fetching categories:", err)
		WriteError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	WriteJSON(w, http.StatusOK, categories)
}

// Create godoc
// @Summary      Create a new category
// @Description  Create a new category with the provided data
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      domain.CategoryInput  true  "Category data"
// @Success      201       {object}  domain.Category
// @Failure      400       {string}  string  "Invalid request body"
// @Router       /categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(&category); err != nil {
		log.Println("Error creating category:", err)
		WriteError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	WriteJSON(w, http.StatusCreated, category)
}

// HandleCategoryByID handles GET, PUT, DELETE requests for /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
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
// @Summary      Get category by ID
// @Description  Retrieve a single category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  domain.Category
// @Failure      400  {string}  string  "Invalid category ID"
// @Failure      404  {string}  string  "Category not found"
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/categories/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		log.Println("Error fetching category by ID:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch category")
		return
	}

	WriteJSON(w, http.StatusOK, category)
}

// Update godoc
// @Summary      Update a category
// @Description  Update an existing category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path      int                   true  "Category ID"
// @Param        category  body      domain.CategoryInput  true  "Category data"
// @Success      200       {object}  domain.Category
// @Failure      400       {string}  string  "Invalid category ID or request body"
// @Router       /categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/categories/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category.ID = id
	if err := h.service.Update(&category); err != nil {
		log.Println("Error updating category:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	WriteJSON(w, http.StatusOK, category)
}

// Delete godoc
// @Summary      Delete a category
// @Description  Delete a category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  handler.APIResponse  "Category deleted successfully"
// @Failure      400  {string}  string  "Invalid category ID"
// @Failure      404  {string}  string  "Category not found"
// @Failure      500  {string}  string  "Failed to delete category"
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/api/categories/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Println("Error deleting category:", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
