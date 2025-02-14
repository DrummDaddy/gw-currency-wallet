package postgres

import (
	"context"
	"database/sql"
)

type Balance struct {
	USD float64 `json:"USD"`
	RUB float64 `json:"RUB"`
	EUR float64 `json""EUR"`
}

func (p *PostgresAdapter) GetUserBalance(ctx context.Context, userID int64) (*Balance, error) {
	query := `
		SELECT usd_balance, rub_balance, eur_balance
		FROM wallet_balances
		WHERE user_id = $1;
	`

	row := p.Storage.db.QueryRowContext(ctx, query, userID)

	var balance Balance
	if err := row.Scan(&balance.USD, &balance.RUB, &balance.EUR); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}

	return &balance, nil
}
