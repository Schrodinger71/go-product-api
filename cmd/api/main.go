package main

import (
	"fmt"
	"go-product-api/internal/config"
	"go-product-api/internal/handlers"
	"go-product-api/internal/repository"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("Configuration Loaded!")
	fmt.Println("Database URL:", cfg.DatabaseUrl)
	fmt.Println("Server Port:", cfg.ServerPort)

	repo, err := repository.NewPostgresRepository(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer repo.Close()

	fmt.Println("Database connection test passed!")

	// Инициализация обработчиков
	productHandler := handlers.NewProductHandler(repo)

	// Настройка маршрутов
	http.HandleFunc("/", productHandler.IndexPage)
	http.HandleFunc("/api/products", productHandler.GetProductsAPI)
	http.HandleFunc("/api/products/create", productHandler.CreateProductAPI)
	http.HandleFunc("/api/products/delete/", productHandler.DeleteProductAPI)

	// Поддержка статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Server starting on port %s...\n", cfg.ServerPort)
	fmt.Printf("Open http://localhost:%s in your browser\n", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
