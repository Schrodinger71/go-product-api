package middleware

import (
	"context"
	"fmt"
	"go-product-api/internal/repository"
	"net/http"
	"strings"
	"time"
)

// Константы для контекстных ключей
type contextKey string

const (
	UserIDKey   contextKey = "userID"
	IsAdminKey  contextKey = "isAdmin"
	UserNameKey contextKey = "userName"
)

var (
	sessionManager *SessionManager
)

// InitializeAuth инициализирует менеджер сессий
func InitializeAuth(userRepo repository.UserRepository) {
	sessionManager = NewSessionManager(userRepo)

	// Запускаем очистку просроченных сессий каждые 30 минут
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			sessionManager.CleanupExpiredSessions()
		}
	}()
}

// AuthMiddleware проверяет аутентификацию
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем статические файлы и страницу логина
		if strings.HasPrefix(r.URL.Path, "/static/") ||
			r.URL.Path == "/login" ||
			r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем сессию
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			// Нет сессии - редирект на логин
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session := sessionManager.GetSession(sessionToken.Value)
		if session == nil {
			// Невалидная сессия - редирект на логин
			http.SetCookie(w, &http.Cookie{
				Name:   "session_token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Обновляем контекст с данными пользователя
		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, IsAdminKey, session.IsAdmin)
		ctx = context.WithValue(ctx, UserNameKey, session.UserName)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AdminOnlyMiddleware проверяет права администратора
func AdminOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
		if !ok || !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// LoginHandler обрабатывает логин
func LoginHandler(userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Ищем пользователя по email
		user, err := userRepo.GetByEmail(email)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// TODO: ДОБАВИТЬ ХЭШИРОВАНИЕ ПАРОЛЯ
		// НЕ ДЛЯ ПРОДА
		if user.Password != password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Создаем сессию
		token := sessionManager.CreateSession(user)

		// Устанавливаем cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		// Возвращаем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": true, "message": "Login successful"}`)
	}
}

// LogoutHandler обрабатывает выход
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err == nil {
			sessionManager.DeleteSession(sessionToken.Value)
		}

		// Удаляем cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Вспомогательные функции для handlers
func GetUserID(ctx context.Context) uint {
	if userID, ok := ctx.Value(UserIDKey).(uint); ok {
		return userID
	}
	return 0
}

func IsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(IsAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}

func GetUserName(ctx context.Context) string {
	if userName, ok := ctx.Value(UserNameKey).(string); ok {
		return userName
	}
	return ""
}
