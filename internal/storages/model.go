package storages

import "time"

type Wallet struct {
	ID       string
	UserID   string
	Balances map[string]float64
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
