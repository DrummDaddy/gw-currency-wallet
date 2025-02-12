package postgres

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/storages"
)

func (p *PosthgresStorage) CreateWallet(ctx context.Context, userID string) (string, error) {
	var walletID string
	err := p.db.QueryRowContext(ctx,
		`INSERT INTO wallets (user_id)
	VALUES ($1) RETURNING id`, userID).Scan(&walletID)
	if err != nil {
		return " ", err
	}
	return walletID, nil
}

func (p *PosthgresStorage) GetWallet(ctx context.Context, walletID string) (*storages.Wallet, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT id, user_id, FROM wallets WHERE id = $1`, walletID)

	var wallet storages.Wallet
	err := row.Scan(&wallet.ID, &wallet.UserID)
	if err != nil {
		return nil, err
	}

	balanceRows, err := p.db.QueryContext(ctx,
		`SELECT currency, balance FROM wallet_balancec WHERE wallet_id = $1`, walletID)
	if err != nil {
		return nil, err
	}
	defer balanceRows.Close()

	balances := make(map[string]float64)
	for balanceRows.Next() {
		var currency string
		var balance float64
		if err := balanceRows.Scan(&currency, &balance); err != nil {
			return nil, err
		}
		balances[currency] = balance
	}
	wallet.Balances = balances

	return &wallet, nil
}

func (p *PosthgresStorage) AddFunds(ctx context.Context, walletID, currency string, amount float64) error {
	result, err := p.db.ExecContext(ctx,
		`INSERT INTO wallet_balances (wallet_id, currency, balance)
	VALUES ($1, $2, $3)
	ON CONFLICT (wallet_id, currency) DO UPDATE
	SET balance = wallet_balances.balance + $3`, walletID, currency, amount)
	if err != nil {
		return err

	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("wallet not found")
	}

	return nil

}
