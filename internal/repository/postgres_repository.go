package repository

import (
	"database/sql"
	"fmt"
	"go-product-api/internal/config"
	"log"
)

type postgresRepository struct {
	Products ProductRepository
	Users    UserRepository
	db       *sql.DB
}

func NewPostgresRepository(dsn string) (PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTableProduct(db)
	if err != nil {
		return nil, err
	}

	err = createTableUsers(db)
	if err != nil {
		if err.Error() == "table already exists" {
			fmt.Println("✅ Table users already exists")
		} else {
			return nil, err
		}
	}

	err = createAdminUser(db)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to PostgreSQL successfully")
	return &postgresRepository{db: db}, nil
}

// GetProductRepository возвращает ProductRepository
func (r *postgresRepository) GetProductRepository() ProductRepository {
	return &productRepository{db: r.db}
}

// GetUserRepository возвращает UserRepository  
func (r *postgresRepository) GetUserRepository() UserRepository {
	return &userRepository{db: r.db}
}

// Создаём таблицу products, если она не существует
func createTableProduct(db *sql.DB) error {
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

func createTableUsers(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(email)
	)`

	_, err := db.Exec(query)
	return err
}

func createAdminUser(db *sql.DB) error {
	adminEmail := config.GetsEnv("ADMIN_EMAIL", "admin@example.com")
	adminName := config.GetsEnv("ADMIN_NAME", "Administrator")

	// TODO: ПОМЕНЯТЬ НА СОЛЬ, ПАРОЛИ НЕ ХРАНЯТ В ОТКРЫТОМ ВИДЕ!!!!!!
	adminPassword := config.GetsEnv("ADMIN_PASSWORD", "321admin123")

	query := `
    INSERT INTO users (name, email, password, is_admin) 
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (email) DO UPDATE SET
        name = EXCLUDED.name,
        password = EXCLUDED.password,
        is_admin = EXCLUDED.is_admin,
        updated_at = CURRENT_TIMESTAMP`

	result, err := db.Exec(query, adminName, adminEmail, adminPassword, true)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("✅ Admin user created/updated: %s", adminEmail)
	}

	return nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}
