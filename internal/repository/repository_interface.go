package repository

import "go-product-api/internal/models"

type PostgresRepository interface {
	GetProductRepository() ProductRepository
	GetUserRepository() UserRepository
	// Закрыть соединение с базой данных
	Close() error
}

type ProductRepository interface {
	// Получить все продукты
	GetAll() ([]models.Product, error)

	// Получить продукт по ID
	GetByID(id uint) (*models.Product, error)

	// Добавить новый продукт
	Create(product *models.Product) error

	// Обновить существующий продукт
	Update(product *models.Product) error

	// Удалить продукт по ID
	Delete(id uint) error
}

type UserRepository interface {
	// Получить все пользователи
	GetAll() ([]models.User, error)

	// Получить пользователя по ID
	GetByID(id uint) (*models.User, error)

	// Получить пользователя по email
	GetByEmail(email string) (*models.User, error)

	// Создать нового пользователя
	Create(user *models.User) error

	// Обновить существующего пользователя
	Update(user *models.User) error

	// Удалить пользователя по ID
	Delete(id uint) error
}
