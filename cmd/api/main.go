package main

import (
	"fmt"
	"go-product-api/internal/config"
	"go-product-api/internal/handlers"
	"go-product-api/internal/middleware"
	"go-product-api/internal/repository"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	// Создаем главный репозиторий
	repo, err := repository.NewPostgresRepository(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer repo.Close()

	// Получаем конкретные репозитории
	productRepo := repo.GetProductRepository()
	userRepo := repo.GetUserRepository()

	// Инициализируем аутентификацию
	middleware.InitializeAuth(userRepo)

	// Создаем handlers
	productHandler := handlers.NewProductHandler(productRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Настройка маршрутов
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})
	http.HandleFunc("/api/login", middleware.LoginHandler(userRepo))
	http.HandleFunc("/logout", middleware.LogoutHandler())

	// Защищенные роуты
	http.HandleFunc("/", middleware.AuthMiddleware(productHandler.IndexPage))
	http.HandleFunc("/api/products", middleware.AuthMiddleware(productHandler.GetProductsAPI))
	http.HandleFunc("/api/products/create", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(productHandler.CreateProductAPI)))
	http.HandleFunc("/api/products/delete/", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(productHandler.DeleteProductAPI)))

	http.HandleFunc("/api/users", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(userHandler.GetUsersAPI)))
	http.HandleFunc("/api/users/create", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(userHandler.CreateUserAPI)))
	http.HandleFunc("/api/users/delete/", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(userHandler.DeleteUserAPI)))

	http.HandleFunc("/admin/users", middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(userHandler.UsersPage)))

	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Server starting on port %s...\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
