package storages

import (
	"context"
	"errors"
	"strings"
)

// Storage - интерфейс для работы с кошельками
type Storage interface {
	CreateWallet(ctx context.Context, userID string) (string, error)
	GetWallet(ctx context.Context, walletID string) (*Wallet, error)
	AddFunds(ctx context.Context, currency string, amount float64) error
}

var wallets = make(map[string]*Wallet)

// GetUserWallet возвращает кошелек пользователя
func GetUserWallet(userID string) (*Wallet, error) {
	wallet, exists := wallets[userID]
	if !exists {
		return nil, errors.New("wallet not found")
	}
	return wallet, nil
}

// SaveUserWallet сохраняет изменения в кошельке
func SaveUserWallet(wallet *Wallet) error {
	if wallet == nil {
		return errors.New("wallet is nil")
	}
	wallets[wallet.UserID] = wallet
	return nil
}

// CheckBalance проверяет, достаточно ли средсв в кошельке для обмена.
func CheckBalance(wallet *Wallet, currency string, amount float64) bool {
	balance, exists := wallet.Balances[currency]
	if !exists {

		return false

	}
	return balance >= amount
}

// UpadateBalance обновляет баланс кошелька по результатам обмена валют
func UpadateBalance(wallet *Wallet, fromCurrency string, toCurrency string, exchangedAmount, deductedAmount float64) {
	wallet.Balances[fromCurrency] -= deductedAmount
	wallet.Balances[toCurrency] += exchangedAmount
	SaveUserWallet(wallet)
}

// GetUserIDFromToken извлекает userID из токена
func GetUserIDFromToken(token string) string {
	parts := strings.Split(token, "")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return " "
	}
	userID := parts[1]
	if userID == " " {
		return ""
	}

	return userID
}
