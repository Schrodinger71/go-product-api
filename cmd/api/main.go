package main

import (
	"fmt"
	"go-product-api/internal/config"
	"go-product-api/internal/models"
	"go-product-api/internal/repository"
)

func main() {
	var cfg = config.LoadConfig()

	fmt.Println("Configuration Loaded!")
	fmt.Println("Database URL:", cfg.DatabaseUrl)
	fmt.Println("Server Port:", cfg.ServerPort)
	fmt.Println("Redis URL:", cfg.RedisUrl)

	repo, err := repository.NewPostgresRepository(cfg.DatabaseUrl)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	fmt.Println("Database connection test passed!")

	testProduct := &models.Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Stock:       100,
		Price:       29.99,
	}

	err = repo.Create(testProduct)
	if err != nil {
		fmt.Println("Error creating product:", err)
		return
	}

	products, err := repo.GetAll()
	if err != nil {
		fmt.Println("Error retrieving products:", err)
		return
	}
	fmt.Println("Products in database:")
	for _, p := range products {
		fmt.Printf("- ID: %d, Name: %s, Price: %f\n", p.ID, p.Name, p.Price)
	}

}
