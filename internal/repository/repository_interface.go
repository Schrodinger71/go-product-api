package repository

import "go-product-api/internal/models"

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
