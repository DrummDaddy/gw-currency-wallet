package auth

import (
	"context"
	"errors"
)

// Ключ для хранения данных
type contextKey string

const UserIDKey contextKey = "userID"

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContex(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("context is nil")
	}
	UserID, ok := ctx.Value(UserIDKey).(string)
	if !ok || UserID == "" {
		return "", errors.New("User ID not found in context")
	}
	return UserID, nil
}
