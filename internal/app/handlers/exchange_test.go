package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleExchange_Success(t *testing.T) {
	// Подготовка данных запроса
	requestBody := struct {
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		Amount       float64 `json:"amount"`
	}{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       100,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/exchange", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Authorization", "Bearer valid_token")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleExchange)

	// Выполняем тестовый запрос
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус: %v, получен: %v", http.StatusOK, rr.Code)
	}

	// Проверяем тело ответа
	expectedResponse := `{"message":"Exchange successful"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Ожидался ответ: %v, получен: %v", expectedResponse, rr.Body.String())
	}
}
