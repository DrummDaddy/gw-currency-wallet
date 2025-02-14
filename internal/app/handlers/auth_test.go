package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler_Success(t *testing.T) {
	// Создаем правильный запрос
	requestBody := LoginRequest{
		Username: "admin",
		Password: "password123",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	// Выполняем тестовый запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус: %v, получен: %v", http.StatusOK, rr.Code)
	}

	// Проверяем тело ответа
	expectedToken := "your_generated_token" // Подставьте правильное значение
	var responseData LoginResponse
	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Не удалось декодировать ответ: %v", err)
	}
	if responseData.Token == "" {
		t.Errorf("Ответ не содержит токена", expectedToken)
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	requestBody := LoginRequest{
		Username: "wrong_user",
		Password: "wrong_password",
	}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус: %v, получен: %v", http.StatusUnauthorized, rr.Code)
	}
}
