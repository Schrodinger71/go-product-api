package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"go-product-api/internal/models"
	"go-product-api/internal/repository"
	"sync"
	"time"
)

type Session struct {
	UserID   uint
	UserName string
	IsAdmin  bool
	Expires  time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	userRepo repository.UserRepository
}

func NewSessionManager(userRepo repository.UserRepository) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		userRepo: userRepo,
	}
}

// CreateSession создает новую сессию
func (sm *SessionManager) CreateSession(user *models.User) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Генерируем случайный токен сессии
	token := generateSessionToken()

	sm.sessions[token] = &Session{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin,
		Expires:  time.Now().Add(24 * time.Hour), // Сессия на 24 часа
	}

	return token
}

// GetSession возвращает сессию по токену
func (sm *SessionManager) GetSession(token string) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[token]
	if !exists || time.Now().After(session.Expires) {
		return nil
	}

	return session
}

// DeleteSession удаляет сессию
func (sm *SessionManager) DeleteSession(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, token)
}

// CleanupExpiredSessions очищает просроченные сессии
func (sm *SessionManager) CleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.After(session.Expires) {
			delete(sm.sessions, token)
		}
	}
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
