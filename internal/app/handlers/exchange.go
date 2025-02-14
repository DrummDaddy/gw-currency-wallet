package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gw-currency-wallet/internal/storages"

	exchange_grpc "github.com/DrummDaddy/proto-exchange/exchange"
	"google.golang.org/grpc"
)

// HandleExchange Обработчик API обмена валют
// @Summary Обмен валют
// @Description Производит обмен одной валюты на другую
// @Tags exchange
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен авторизации"
// @Param request body struct{from_currency string; to_currency string; amount float64} true "Данные для обмена"
// @Success 200 {object} struct{message string; exchanged_amount float64; new_balance map[string]float64}
// @Failure 400 {object} struct{error string}
// @Failure 401 {object} struct{error string}
// @Failure 405 {object} struct{error string}
// @Router /exchange [post]

// gRPCClient содержит подключение к ExchangeService
type gRPCClient struct {
	client exchange_grpc.ExchangeServiceClient
	conn   *grpc.ClientConn
}

// Cache структура для хранения курсов валют
var exchangeCache = make(map[string]float64)
var cacheExpiration = make(map[string]time.Time)

// getCacheRate возвращает курс из кэша или ошибку
func getCacheRate(fromCurrency, toCurrency string) (float64, bool) {
	key := fromCurrency + "_" + toCurrency
	rate, exists := exchangeCache[key]
	if !exists || time.Now().After(cacheExpiration[key]) {
		return 0, false
	}
	return rate, true
}

// setCachedRate сохраняет курс валют в кэшэ
func setCachedRate(fromCurrency, toCurrency string, rate float64) {
	key := fromCurrency + "_" + toCurrency
	exchangeCache[key] = rate
	cacheExpiration[key] = time.Now().Add(5 * time.Minute)
}

// NewExchangeClient создает новое подключение к gRPC серверу
func NewExchangeClient(serverAddres string) (*gRPCClient, error) {
	conn, err := grpc.Dial(serverAddres, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect gRPC server: %w, err")
	}

	client := exchange_grpc.NewExchangeServiceClient(conn)
	return &gRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

var grpcClient *gRPCClient

// initGrpcClient инициализирует gRPC клиент (один раз)
func initGrpcClient() error {
	var err error
	serverAddres := "localhost:50051"
	grpcClient, err = NewExchangeClient(serverAddres)
	if err != nil {
		return fmt.Errorf("failed to initialize gRPCclient: %w", err)
	}
	return nil
}

// fetchExchangeRate получает курс валют через gRPC
func fetchExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	if grpcClient == nil {
		if err := initGrpcClient(); err != nil {
			return 0, err
		}
	}

	return 0, nil
}

// handleExchange метод для обработки обмена валют
func HandleExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		FromCurrency string `json:"from_currency"`
		ToCurrency   string `json:"to_currency"`
		Amount       float64
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Получение кошелька пользователя
	userID := storages.GetUserIDFromToken(token)
	wallet, err := storages.GetUserWallet(userID)
	if err != nil {
		http.Error(w, "Wallet not found", http.StatusBadRequest)
		return
	}

	if !storages.CheckBalance(wallet, request.FromCurrency, request.Amount) {
		http.Error(w, `{"error": Insufficient funds or invalid currencies}`, http.StatusBadRequest)
		return
	}

	//Получение курса валют
	var exchangeRate float64

	exchangeRate, cached := getCacheRate(request.FromCurrency, request.ToCurrency)
	if !cached {
		exchangeRate, err = fetchExchangeRate(request.FromCurrency, request.ToCurrency)
		if err != nil {
			http.Error(w, `{"error":"Failed to retrieve exchange rate"}`, http.StatusBadRequest)
			return
		}
		setCachedRate(request.FromCurrency, request.ToCurrency, exchangeRate)

		//Расчет сумму обмена
		exhchangedAmount := request.Amount * exchangeRate

		//Обновление кошелька
		storages.UpadateBalance(wallet, request.FromCurrency, request.ToCurrency, exhchangedAmount, request.Amount)

		response := struct {
			Message         string             `json:"message"`
			ExchangedAmount float64            `json:"exchanged_amount"`
			NewBalance      map[string]float64 `json:"new_balance"`
		}{
			Message:         "Exchange succeful",
			ExchangedAmount: exhchangedAmount,
			NewBalance:      wallet.Balances,
		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)

	}

}
