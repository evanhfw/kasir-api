# Kasir API

REST API for a Point-of-Sale (POS) system built with Go. Manages products and categories with PostgreSQL database.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25.6 |
| Framework | Standard library `net/http` |
| Database | PostgreSQL |
| Config | Viper |
| Documentation | Swagger (swaggo/swag) |

## Features

- CRUD operations for Products and Categories
- Product-Category relationship
- Health check endpoint with database connectivity check
- Swagger UI documentation
- Docker support with multi-stage build

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd kasir-api
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Setup database**
   ```bash
   # Connect to PostgreSQL and run the schema (see Database Schema section)
   psql -U postgres -d kasir -f schema.sql
   ```

4. **Run the application**
   ```bash
   make dev
   # or
   go run ./cmd/api
   ```

5. **Access the API**
   - API: http://localhost:8080/api/products
   - Health: http://localhost:8080/health
   - Swagger UI: http://localhost:8080/swagger/

### Using Docker

```bash
# Build the image
make build

# Start with docker-compose
make run

# View logs
make logs

# Stop
make stop
```

## API Endpoints

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check (server + database) |

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | List all products |
| POST | `/api/products` | Create a new product |
| GET | `/api/products/{id}` | Get product by ID |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/categories` | List all categories |
| POST | `/api/categories` | Create a new category |
| GET | `/api/categories/{id}` | Get category by ID |
| PUT | `/api/categories/{id}` | Update category |
| DELETE | `/api/categories/{id}` | Delete category |

### Documentation

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/swagger/` | Swagger UI |

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_CONN` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/kasir?sslmode=disable` |

## Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run application locally |
| `make swagger` | Generate Swagger documentation |
| `make build` | Build Docker image |
| `make run` | Start with docker-compose |
| `make stop` | Stop containers |
| `make logs` | View container logs |
| `make clean` | Remove containers and images |
| `go test ./...` | Run all tests |

## Project Structure

```
kasir-api/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── apperrors/            # Custom error definitions
│   ├── config/               # Configuration loading
│   ├── database/             # Database connection
│   ├── domain/               # Domain models (Product, Category)
│   ├── handler/              # HTTP handlers
│   ├── repository/           # Data access layer
│   ├── router/               # HTTP routing
│   └── service/              # Business logic layer
├── docs/                     # Generated Swagger documentation
├── .env.example              # Environment variables template
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # Docker Compose configuration
├── Makefile                  # Development commands
└── go.mod                    # Go module definition
```

## Architecture

```
HTTP Request
     │
     ▼
┌─────────┐    ┌─────────┐    ┌────────────┐    ┌──────────┐
│ Handler │───▶│ Service │───▶│ Repository │───▶│ Database │
└─────────┘    └─────────┘    └────────────┘    └──────────┘
     │
     ▼
HTTP Response
```

- **Handler**: HTTP layer, request/response handling, validation
- **Service**: Business logic, data orchestration
- **Repository**: Data access, SQL queries
- **Database**: PostgreSQL storage

## Database Schema

```sql
-- Create categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Create products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER NOT NULL REFERENCES categories(id)
);

-- Create index for faster product lookups by category
CREATE INDEX idx_products_category_id ON products(category_id);
```

## API Response Format

All API responses follow a consistent format:

### Success Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Indomie Goreng",
    "price": 3500,
    "stock": 100,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Makanan Ringan",
      "description": "Kategori untuk makanan ringan"
    }
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Product not found"
}
```

### Health Check Response

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "database": "connected"
  }
}
```

## Example Requests

### Create a Category

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Makanan Ringan", "description": "Snacks and light food"}'
```

### Create a Product

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Indomie Goreng", "price": 3500, "stock": 100, "category_id": 1}'
```

### Get All Products

```bash
curl http://localhost:8080/api/products
```

### Health Check

```bash
curl http://localhost:8080/health
```
