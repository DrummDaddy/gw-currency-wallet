package postgres

import (
	"database/sql"

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
