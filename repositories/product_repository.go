package repositories

import (
	"database/sql"

	apperrors "kasir-api/errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
		       c.id, c.name, c.description
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var c models.Category
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
			&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		p.Category = &c
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
		       c.id, c.name, c.description
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	row := r.db.QueryRow(query, id)

	var p models.Product
	var c models.Category
	if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
		&c.ID, &c.Name, &c.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	p.Category = &c

	return &p, nil
}

func (r *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
