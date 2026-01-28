package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port	string `mapstructure:"PORT"`
	DBConn  string `mapstructure:"DB_CONN"`
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

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Error reading config file:", err)
		}
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		log.Fatal("PORT belum diset")
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	addr := "0.0.0.0:" + config.Port
	log.Println("Server running di", addr)

	if err := http.ListenAndServe(addr, nil);  err != nil {
		log.Fatal("Gagal running server:", err)
	}
}