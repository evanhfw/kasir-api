# AGENTS.md - Kasir API

Guidelines for AI coding agents working in this Go REST API for a Point-of-Sale system.

## Tech Stack

| Aspect | Details |
|--------|---------|
| Language | Go 1.25.6 |
| Framework | Standard library `net/http` |
| Database | PostgreSQL (`lib/pq` driver) |
| Config | Viper (loads `.env` files) |
| Docs | Swagger via `swaggo/swag` |

## Project Structure

```
cmd/api/main.go           # Entry point, dependency wiring
internal/
  config/                 # Configuration loading (Viper)
  database/               # PostgreSQL connection
  domain/                 # Domain models (Product, Category)
  handler/                # HTTP handlers (controllers)
  service/                # Business logic layer
  repository/             # Data access layer + interfaces
  apperrors/              # Custom application errors
  router/                 # HTTP routing setup
docs/                     # Generated Swagger documentation
```

Architecture: `Handler -> Service -> Repository -> Database`

## Build & Run Commands

```bash
# Development
make dev                  # Run locally with hot reload
go run ./cmd/api          # Run application

# Build
go build -o kasir-api ./cmd/api
make build                # Build Docker image

# Dependencies
go mod tidy               # Clean up go.mod/go.sum

# Swagger docs generation
make swagger              # swag init -g cmd/api/main.go -o docs

# Docker
make run                  # Start with docker-compose
make stop                 # Stop application
make logs                 # View container logs
make clean                # Remove containers and images
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/handler
go test ./internal/service

# Run single test by name
go test -run TestGetAllProducts ./internal/handler
go test -run TestCreate ./internal/service

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # View coverage report
```

**Note:** Test files should be created as `*_test.go` alongside source files.

## Code Style Guidelines

### Import Ordering

Group imports in this order, separated by blank lines:
1. Standard library
2. Internal packages (`kasir-api/...`)
3. External dependencies

```go
import (
    "encoding/json"
    "net/http"

    "kasir-api/internal/apperrors"
    "kasir-api/internal/domain"
    "kasir-api/internal/service"

    "github.com/spf13/viper"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Packages | lowercase, singular | `handler`, `service`, `domain` |
| Files | snake_case | `product_handler.go`, `category_service.go` |
| Structs | PascalCase | `ProductHandler`, `CategoryService` |
| Interfaces | PascalCase + descriptive | `ProductRepository`, `CategoryRepository` |
| Public functions | PascalCase | `GetAll`, `Create`, `NewProductHandler` |
| Private functions | camelCase | `parseIDFromPath` |
| Variables | camelCase | `productRepo`, `categoryService` |
| Receivers | Single lowercase letter | `h`, `s`, `r` |
| Constants | PascalCase or ALL_CAPS | `ErrNotFound` |

### Domain Models

Define separate structs for entities and input:

```go
// Entity with ID - used for responses
type Product struct {
    ID         int       `json:"id" example:"1"`
    Name       string    `json:"name" example:"Indomie Goreng"`
    Price      int       `json:"price" example:"3500"`
    Stock      int       `json:"stock" example:"100"`
    CategoryID int       `json:"category_id" example:"1"`
    Category   *Category `json:"category,omitempty"`
}

// Input without ID - used for create/update requests
type ProductInput struct {
    Name       string `json:"name" example:"Indomie Goreng"`
    Price      int    `json:"price" example:"3500"`
    Stock      int    `json:"stock" example:"100"`
    CategoryID int    `json:"category_id" example:"1"`
}
```

### Error Handling

Use sentinel errors from `internal/apperrors`:

```go
// In apperrors package - define custom errors
var (
    ErrNotFound         = errors.New("resource not found")
    ErrCategoryNotFound = errors.New("category not found")
    ErrInvalidInput     = errors.New("invalid input")
)

// In repository - return apperrors for known conditions
if err == sql.ErrNoRows {
    return nil, apperrors.ErrNotFound
}

// In service - check and wrap errors
if errors.Is(err, apperrors.ErrNotFound) {
    return apperrors.ErrCategoryNotFound
}

// In handler - map errors to HTTP status codes
if errors.Is(err, apperrors.ErrNotFound) {
    WriteError(w, http.StatusNotFound, "Product not found")
    return
}
```

### HTTP Handler Pattern

Use method routing within handlers:

```go
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
```

### Response Format

Use `WriteJSON` and `WriteError` from `handler/response.go`:

```go
// Success response - wraps data in APIResponse
WriteJSON(w, http.StatusOK, products)
WriteJSON(w, http.StatusCreated, product)

// Error response
WriteError(w, http.StatusBadRequest, "Invalid request body")
WriteError(w, http.StatusNotFound, "Product not found")
```

### Swagger Annotations

Add godoc comments before handler methods:

```go
// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      500  {string}  string  "Failed to fetch products"
// @Router       /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
```

### Dependency Injection

Use constructor functions with explicit dependencies:

```go
// Repository - takes *sql.DB, returns interface
func NewProductRepository(db *sql.DB) ProductRepository {
    return &productRepository{db: db}
}

// Service - takes repository interfaces
func NewProductService(productRepo ProductRepository, categoryRepo CategoryRepository) *ProductService {
    return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

// Handler - takes service pointer
func NewProductHandler(service *service.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}
```

### Database Queries

- Use PostgreSQL placeholders: `$1`, `$2`, etc.
- Always `defer rows.Close()` after Query
- Check `rows.Err()` after iteration loop
- Check `RowsAffected()` for UPDATE/DELETE operations
- Use `RETURNING id` for INSERT to get generated ID

```go
func (r *productRepository) Create(product *domain.Product) error {
    query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
    return r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_CONN` | PostgreSQL connection string | `postgres://user:pass@host:5432/db?sslmode=disable` |

Copy `.env.example` to `.env` for local development.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | List all products |
| POST | `/api/products` | Create a product |
| GET | `/api/products/{id}` | Get product by ID |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |
| GET | `/api/categories` | List all categories |
| POST | `/api/categories` | Create a category |
| GET | `/api/categories/{id}` | Get category by ID |
| PUT | `/api/categories/{id}` | Update category |
| DELETE | `/api/categories/{id}` | Delete category |
| GET | `/swagger/` | Swagger UI documentation |
