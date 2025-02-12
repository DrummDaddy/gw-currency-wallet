package postgres

import (
	"context"
	"gw-currency-wallet/internal/storages"
)

type PostgresAdapter struct {
	Storage *PosthgresStorage
}

func (a *PostgresAdapter) CreateWallet(ctx context.Context, userID string) (string, error) {
	return a.Storage.CreateWallet(ctx, userID)
}

func (a *PostgresAdapter) GetWallet(ctx context.Context, walletID string) (*storages.Wallet, error) {
	return a.Storage.GetWallet(ctx, walletID)
}

// Адаптация сигнатуры AddFunds
func (a *PostgresAdapter) AddFunds(ctx context.Context, currency string, amount float64) error {
	// Используем фиктивный walletID или определяем его логику
	walletID := "default-wallet-id"

	// Вызываем метод AddFunds из PosthgresStorage
	return a.Storage.AddFunds(ctx, walletID, currency, amount)
}
