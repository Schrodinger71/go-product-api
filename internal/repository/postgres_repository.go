package repository

import (
	"database/sql"
	"fmt"
	"go-product-api/internal/models"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (ProductRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to PostgreSQL successfully")
	return &postgresRepository{db: db}, nil
}

// Создаём таблицу products, если она не существует
func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		stock INTEGER NOT NULL DEFAULT 0,
		price DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}

func (r *postgresRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, name, description, stock, price, created_at, updated_at FROM products ORDER BY id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Stock,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *postgresRepository) GetByID(id uint) (*models.Product, error) {
	query := `SELECT id, name, description, stock, price, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var product models.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Stock,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *postgresRepository) Create(product *models.Product) error {
	query := `
        INSERT INTO products (name, description, stock, price) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Stock,
		product.Price,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func (r *postgresRepository) Update(product *models.Product) error {
	query := `
        UPDATE products 
        SET name = $1, description = $2, stock = $3, price = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
        RETURNING updated_at
    `

	i := 0

	for i < 10000000 {
		i++
	}

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Stock,
		product.Price,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("product not found")
	}
	return err
}

func (r *postgresRepository) Delete(id uint) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}
