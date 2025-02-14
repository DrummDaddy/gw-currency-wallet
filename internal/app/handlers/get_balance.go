package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GetBalanceHandler Получение баланса кошелька
// @Summary Получить баланс
// @Description Пользователь запрашивает баланс своего кошелька
// @Tags wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен авторизации"
// @Success 200 {object} BalanceResponse "Баланс пользователя"
// @Failure 401 {string} string "Неавторизованный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /wallet/balance [get]

type BalanceResponse struct {
	Balance struct {
		USD float64 `json:"USD"`
		RUB float64 `json:"RUB"`
		EUR float64 `json""EUR"`
	} `json:"balance"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Unathorized", http.StatusUnauthorized)
			return
		}

		//Декодирование и проверка токена
		token := strings.TrimPrefix(authHeader, "Bearer")
		if token == " " {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	balance := BalanceResponse{}
	balance.Balance.USD = 100.50
	balance.Balance.RUB = 5000.00
	balance.Balance.EUR = 85.75

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
