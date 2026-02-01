package main

import (
	"log"
	"net/http"

	"kasir-api/internal/config"
	"kasir-api/internal/database"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/router"
	"kasir-api/internal/service"

	_ "kasir-api/docs"
)

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
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if cfg.Port == "" {
		log.Fatal("PORT belum diset")
	}

	// Initialize database
	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	healthHandler := handler.NewHealthHandler(db)

	// Setup router
	r := router.New(productHandler, categoryHandler, healthHandler)

	// Start server
	addr := "0.0.0.0:" + cfg.Port
	log.Println("Server running di", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Gagal running server:", err)
	}
}
