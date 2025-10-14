package handlers

import (
	"encoding/json"
	"go-product-api/internal/middleware"
	"go-product-api/internal/models"
	"go-product-api/internal/repository"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	repo repository.ProductRepository
	tmpl *template.Template
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	return &ProductHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// HTML страница
func (h *ProductHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	if !middleware.IsAdmin(r.Context()) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	products, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}

	h.tmpl.Execute(w, products)
}

// API: Получить все товары JSON
func (h *ProductHandler) GetProductsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// API: Создать товар только для админа
func (h *ProductHandler) CreateProductAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid price"}`, http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid stock"}`, http.StatusBadRequest)
		return
	}

	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	err = h.repo.Create(product)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Редирект обратно на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// API: Удалить товар по ID
func (h *ProductHandler) DeleteProductAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/products/delete/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, `{"error": "Invalid product ID"}`, http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(uint(id))
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Редирект обратно на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
