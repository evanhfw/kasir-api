package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 15000000, Stock: 10},
	{ID: 2, Name: "Smartphone", Price: 5000000, Stock: 25},
	{ID: 3, Name: "Tablet", Price: 3000000, Stock: 15},
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var updateProduct Product
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			products[i].Name = updateProduct.Name
			products[i].Price = updateProduct.Price
			products[i].Stock = updateProduct.Stock

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products[i])
			return
		}
	}

	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent) //204
			return
		}
	}

	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func main() {
	// GET /api/products/{id}
	// PUT /api/products/{id}
	// DELETE /api/products/{id}
	http.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			getProductByID(w, r)
		case "PUT":
			updateProduct(w, r)
		case "DELETE":
			deleteProduct(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/products
	// Post /api/products
	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			json.NewEncoder(w).Encode(products)
		case "POST":
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			newProduct.ID = len(products) + 1
			products = append(products, newProduct)

			w.WriteHeader(http.StatusCreated) //201
			json.NewEncoder(w).Encode(newProduct)
		}

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal Running Server")
	}
}
