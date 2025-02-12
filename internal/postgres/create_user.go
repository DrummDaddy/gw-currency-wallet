package postgres

import (
	"database/sql"
)

type PostgresUser struct {
	Storage *sql.DB
}

// Проверка существования пользователя по имени или email
func (p *PostgresUser) UserExists(username, email string) (bool, error) {
	var exists bool
	query :=
		`SELECT EXISTS(
	SELECT 1
	FROM users
	WHERE username = $1 OR email = $2 
	)`

	err := p.Storage.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err

	}
	return exists, nil
}

// Создание нового пользователя
func (p *PostgresUser) CreateUser(username, password, email string) error {
	query :=
		`INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)`

	_, err := p.Storage.Exec(query, username, password, email)
	return err
}
