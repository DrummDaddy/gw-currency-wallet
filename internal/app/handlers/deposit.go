package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gw-currency-wallet/internal/auth"
)

type DepositRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type DepositResponse struct {
	Message    string             `json:"message"`
	NewBalance map[string]float64 `json:"new_balnce"`
}

// Доступная валюта
var allowedCurencies = map[string]bool{
	"USD": true,
	"RUB": true,
	"EUR": true,
}

// DepositHandler Обработка депозита
// @Summary Добавить депозит
// @Description Пользователь добавляет депозит в свой кошелек
// @Tags wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен авторизации"
// @Param request body DepositRequest true "Данные депозита"
// @Success 200 {object} DepositResponse "Депозит успешно выполнен"
// @Failure 400 {string} string "Неверные данные запроса"
// @Failure 401 {string} string "Неавторизованный запрос"
// @Failure 404 {string} string "Кошелек не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /wallet/deposit [post]

func DepositHandler(db *sql.DB, jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, авторизован ли пользователь
		userID, err := auth.GetUserIDFromContex(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Разбираем запрос на депонирование
		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем, задана ли валюта среди разрешенных
		req.Currency = strings.ToUpper(req.Currency)
		if !allowedCurencies[req.Currency] {
			http.Error(w, "Unsupported currency", http.StatusBadRequest)
			return
		}

		// Убедимся, что сумма депозита положительная
		if req.Amount <= 0 {
			http.Error(w, "Deposit amount must be greater than zero", http.StatusBadRequest)
			return
		}

		// Получаем текущее состояние балансов пользователя
		var usdBalance, rubBalance, eurBalance float64
		query := `SELECT usd_balance, rub_balance, eur_balance FROM wallets WHERE user_id = $1`
		err = db.QueryRow(query, userID).Scan(&usdBalance, &rubBalance, &eurBalance)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Wallet not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Обновляем баланс в зависимости от валюты
		switch req.Currency {
		case "USD":
			usdBalance += req.Amount
		case "RUB":
			rubBalance += req.Amount
		case "EUR":
			eurBalance += req.Amount
		}

		// Сохраняем обновленный баланс в базу данных
		updateQuery := `UPDATE wallets SET usd_balance = $1, rub_balance = $2, eur_balance = $3 WHERE user_id = $4`
		_, err = db.Exec(updateQuery, usdBalance, rubBalance, eurBalance, userID)
		if err != nil {
			http.Error(w, "Failed to update wallet balance", http.StatusInternalServerError)
			return
		}

		// Подготавливаем ответный JSON
		newBalance := map[string]float64{
			"USD": usdBalance,
			"RUB": rubBalance,
			"EUR": eurBalance,
		}
		resp := DepositResponse{
			Message:    "Deposit successful",
			NewBalance: newBalance,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
