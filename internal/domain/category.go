package domain

// Category represents a product category
// @Description Category information
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Makanan Ringan"`
	Description string `json:"description,omitempty" example:"Kategori untuk makanan ringan seperti keripik, biskuit, dll."`
}

// CategoryInput is used for create/update requests
// @Description Category input for create/update
type CategoryInput struct {
	Name        string `json:"name" example:"Makanan Ringan"`
	Description string `json:"description,omitempty" example:"Kategori untuk makanan ringan seperti keripik, biskuit, dll."`
}
