package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gw-currency-wallet/internal/storages"
)

// Withdraw Обработчик API вывода средств
// @Summary Вывод средств
// @Description Уменьшает баланс пользователя на указанную сумму
// @Tags withdraw
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен авторизации"
// @Param request body WithdrawRequest true "Данные о сумме вывода"
// @Success 200 {object} WithdrawResponse
// @Failure 400 {object} struct{error string}
// @Failure 401 {object} struct{error string}
// @Router /withdraw [post]

// WithdrawRequest описывает структуру тела запроса и снятия средств
type WithdrawRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// WithdrawResponce описывает структуру ответа
type WithdrawResponce struct {
	Message    string             `json:"message"`
	NewBalance map[string]float64 `json:"new_balance"`
}

// withdrawFunds - функция обработки вывода средств
func withdrawFunds(wallet *storages.Wallet, amount float64, currency string) error {
	balance, exists := wallet.Balances[currency]
	if !exists {
		return errors.New("currency not supported")
	}

	if balance < amount || amount < 0 {
		return errors.New("insufficient funds or invalid amount")

	}

	wallet.Balances[currency] -= amount
	return nil
}

// Withdraw - обработчик для вывода средств
func Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	//Проверка поддерживаемой валюты
	if req.Currency != "USD" && req.Currency != "RUB" && req.Currency != "EUR" {
		http.Error(w, `{"error": "unsupported currency"}`, http.StatusBadRequest)
		return
	}

	//Загрузка кошелька пользователя из хранилища
	wallet, err := storages.GetUserWallet(userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	res := WithdrawResponce{
		Message:    "Withdraw succeful",
		NewBalance: wallet.Balances,
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	//Сохраняем новый балас
	err = storages.SaveUserWallet(wallet)
	if err != nil {
		http.Error(w, "error for save ", http.StatusBadRequest)
	}

}
