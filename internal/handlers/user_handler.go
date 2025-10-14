package handlers

import (
	"encoding/json"
	"fmt"
	"go-product-api/internal/middleware"
	"go-product-api/internal/models"
	"go-product-api/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// CreateUserAPI - только для администраторов
func (h *UserHandler) CreateUserAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	isAdminStr := r.FormValue("is_admin")

	isAdmin := false
	if isAdminStr == "true" || isAdminStr == "1" {
		isAdmin = true
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password, // В реальности нужно хешировать!
		IsAdmin:  isAdmin,
	}

	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ вместо редиректа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUsersAPI - только для администраторов
func (h *UserHandler) GetUsersAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	users, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUserAPI - только для администраторов
func (h *UserHandler) UpdateUserAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/users/update/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	isAdminStr := r.FormValue("is_admin")

	// Получаем существующего пользователя
	existingUser, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Обновляем только переданные поля
	if name != "" {
		existingUser.Name = name
	}
	if email != "" {
		existingUser.Email = email
	}
	if password != "" {
		existingUser.Password = password // В реальности хешировать!
	}
	if isAdminStr != "" {
		existingUser.IsAdmin = (isAdminStr == "true" || isAdminStr == "1")
	}

	if err := h.userRepo.Update(existingUser); err != nil {
		http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
		"user":    existingUser,
	})
}

// DeleteUserAPI - только для администраторов
func (h *UserHandler) DeleteUserAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/users/delete/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Не позволяем удалить самого себя
	currentUserID := middleware.GetUserID(r.Context())
	if uint(id) == currentUserID {
		http.Error(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("User with ID %d deleted successfully", id),
	})
}

// GetUserByIDAPI - только для администраторов
func (h *UserHandler) GetUserByIDAPI(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
