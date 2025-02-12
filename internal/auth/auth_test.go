package auth_test

import (
	"gw-currency-wallet/internal/auth"
	"testing"
	"time"
)

func TestJWTManager(t *testing.T) {
	jwtManager := auth.NewJWTManager("secret-key", time.Minute)

	userID := "123"
	token, err := jwtManager.Generate(userID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	verifiedUserID, err := jwtManager.Verify(token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
	if verifiedUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, verifiedUserID)
	}
}
