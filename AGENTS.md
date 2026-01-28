# AGENTS.md - Kasir API

Guidelines for AI coding agents working in this Go REST API for a Point-of-Sale system.

## Tech Stack

| Aspect | Details |
|--------|---------|
| Language | Go 1.25.6 |
| Framework | Standard library `net/http` |
| Database | PostgreSQL (`lib/pq` driver) |
| Config | Viper |
| Docs | Swagger via `swaggo/swag` |

## Architecture

Clean/Layered Architecture: `Handler -> Service -> Repository -> Database`

```
main.go              # Entry point, wiring, routes
├── handlers/        # HTTP layer (controllers)
├── services/        # Business logic layer
├── repositories/    # Data access layer
├── models/          # Data structures
├── database/        # DB connection utilities
└── docs/            # Generated Swagger docs
```

## Build & Run Commands

```bash
# Local development
go build -o kasir-api .
go run main.go
go mod tidy

# Generate Swagger docs
swag init

# Docker (via Makefile)
make build     # Build image
make run       # Start with docker-compose
make stop      # Stop application
make logs      # View logs
```

## Testing

```bash
go test ./...                              # Run all tests
go test ./handlers                         # Test specific package
go test -run TestFunctionName ./package    # Run single test
go test -v ./...                           # Verbose output
go test -cover ./...                       # With coverage
```

**Note:** No test files exist yet. Create `*_test.go` files alongside source files.

## Code Style Guidelines

### Import Ordering

Group imports (separated by blank lines): 1) Standard lib, 2) Internal (`kasir-api/...`), 3) External

```go
import (
    "encoding/json"
    "net/http"

    "kasir-api/models"
    "kasir-api/services"

    "github.com/spf13/viper"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Packages | lowercase | `handlers`, `services` |
| Files | snake_case | `product_handler.go` |
| Structs/Functions | PascalCase | `ProductHandler`, `GetAll` |
| Variables | camelCase | `productRepository` |
| Receivers | Single letter | `h`, `s`, `r` |

### Struct Definitions

```go
type Product struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"Indomie Goreng"`
    Price int    `json:"price" example:"3500"`
}

// Separate Input types for requests (no ID)
type ProductInput struct {
    Name  string `json:"name" example:"Indomie Goreng"`
    Price int    `json:"price" example:"3500"`
}
```

### Error Handling

```go
// Repository: return errors directly
if err == sql.ErrNoRows {
    return nil, errors.New("product not found")
}

// Handler: log and return HTTP error
if err != nil {
    log.Println("Error:", err)
    http.Error(w, "Failed", http.StatusInternalServerError)
    return
}

// Check rows affected for UPDATE/DELETE
if rowsAffected == 0 {
    return errors.New("not found")
}
```

### HTTP Handler Pattern

```go
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.GetAll(w, r)
    case http.MethodPost:
        h.Create(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

### Swagger Annotations

```go
// GetAll godoc
// @Summary      Get all products
// @Tags         products
// @Success      200  {array}  models.Product
// @Router       /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
```

### Dependency Injection

```go
func NewProductHandler(service *services.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}
```

### Database Queries

- Use `$1, $2, ...` placeholders (PostgreSQL)
- Always `defer rows.Close()`
- Check `rows.Err()` after iteration

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | Server port (e.g., `8080`) |
| `DB_CONN` | PostgreSQL connection string |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET/POST | `/api/products` | List/Create products |
| GET/PUT/DELETE | `/api/products/{id}` | Get/Update/Delete product |
| GET/POST | `/api/categories` | List/Create categories |
| GET/PUT/DELETE | `/api/categories/{id}` | Get/Update/Delete category |
| GET | `/swagger/` | Swagger UI |
