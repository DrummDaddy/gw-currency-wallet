package storages

type Wallet struct {
	ID       string
	UserID   string
	Balances map[string]float64
}
