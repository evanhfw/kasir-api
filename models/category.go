package models

type Category struct {
	ID				int    `json:"id" example:"1"`
	Name 			string `json:"name" example:"Makanan Ringan"`
	Description 	string `json:"description,omitempty" example:"Kategori untuk makanan ringan seperti keripik, biskuit, dll."`
}

type CategoryInput struct {
	Name 			string `json:"name" example:"Makanan Ringan"`
	Description 	string `json:"description,omitempty" example:"Kategori untuk makanan ringan seperti keripik, biskuit, dll."`
}