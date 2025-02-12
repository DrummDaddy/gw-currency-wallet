package storages

import "context"

//Storage - интерфейс для работы с кошельками
type Storage interface {
	CreateWallet(ctx context.Context, userID string) (string, error)
	GetWallet(ctx context.Context, walletID string) (*Wallet, error)
	AddFunds(ctx context.Context, currency string, amount float64) error
}
