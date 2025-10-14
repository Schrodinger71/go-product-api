package repository

import (
	"database/sql"
	"fmt"
	"go-product-api/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at FROM users ORDER BY id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, is_admin) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.IsAdmin,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password = $3, is_admin = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at`

	result, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
