package models

// Product represents a product in the store
// @Description Product information
type Product struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"Indomie Goreng"`
	Price int    `json:"price" example:"3500"`
	Stock int    `json:"stock" example:"100"`
}

// ProductInput is used for create/update requests
// @Description Product input for create/update
type ProductInput struct {
	Name  string `json:"name" example:"Indomie Goreng"`
	Price int    `json:"price" example:"3500"`
	Stock int    `json:"stock" example:"100"`
}