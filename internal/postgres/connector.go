package postgres

import (
	"database/sql"
	"errors"

	"gw-currency-wallet/internal/storages"

	_ "github.com/lib/pq"
)

type PosthgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage инициализирует подключение к БД
func NewPostgresStorage(dsn string) (*PosthgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PosthgresStorage{
		db: db,
	}, nil
}

func (p *PosthgresStorage) Close() {
	p.db.Close()
}

func CreateWallet(userID int, currencyCode string, balance float64) error {
	_, err := db.Exec(
		"INSERT INTO wallets (user_id, currency_code, balance) VALUES ($1, $2, $3)",
		userID, currencyCode, balance,
	)
	return err
}

func GetWalletByID(id int) (storages.Wallet, error) {
	var wallet storages.Wallet
	err := db.QueryRow(
		"SELECT id, user_id, currency_code, balance FROM wallets WHERE id=$1", id,
	).Scan(&wallet.ID, &wallet.UserID, &wallet.Balances)
	if err != nil {
		return storages.Wallet{}, errors.New("wallet not found")
	}
	return wallet, nil
}
