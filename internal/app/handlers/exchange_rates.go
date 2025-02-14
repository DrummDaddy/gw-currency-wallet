package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	pb "github.com/DrummDaddy/proto-exchange/exchange"
	"google.golang.org/grpc"
)

// ServerHTTP Обработчик для получения курсов валют
// @Summary Получение курсов валют
// @Description Возвращает текущие обменные курсы для заданной валюты
// @Tags rates
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен авторизации"
// @Param base query string true "Базовая валюта"
// @Param symbols query string true "Список целевых валют через запятую (например, USD,EUR)"
// @Success 200 {object} pb.ExchangeRatesResponse
// @Failure 400 {object} struct{error string}
// @Failure 401 {object} struct{error string}
// @Failure 405 {object} struct{error string}
// @Failure 500 {object} struct{error string}
// @Router /rates [get]

// ExchangeHandler структру обмена валют
type ExchangeHandler struct {
	grpcClient pb.ExchangeServiceClient
}

// NewExchangeHandler создает новый экземпляр хендлера
func NewExchangeHandler(grpcCon *grpc.ClientConn) *ExchangeHandler {
	client := pb.NewExchangeServiceClient(grpcCon)
	return &ExchangeHandler{
		grpcClient: client,
	}
}

// ServerHTTP для обработки HTTP - запросов
func (h *ExchangeHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !isValidAuthHeader(authHeader) {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	//Обрабатываем запрос на получение курсов валют
	baseCurrensy := r.URL.Query().Get("base")
	targetCurrencies := r.URL.Query().Get("symbols")

	if baseCurrensy == "" || targetCurrencies == "" {
		http.Error(w, `{"error": "Invalid query parametrs"}`, http.StatusBadRequest)
		return
	}

	//Вызов grpc - клиента
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.GetExchangeRates(ctx, &pb.Empty{})

	if err != nil {
		http.Error(w, `{"error": "Failed to fetch exchange rates"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)

}

func isValidAuthHeader(header string) bool {
	return strings.HasPrefix(header, "Bearer ")
}
