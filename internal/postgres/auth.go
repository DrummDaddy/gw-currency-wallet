package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func initDatabase() {
	var err error
	connStr := "user=postgres password=159357 dbname=currency_wallet sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func authenticateUserDB(username, password string) bool {
	var passwordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username=$1", username).Scan(&passwordHash)
	if err != nil {
		return false //Пользователь не найден

	}

	//Сравнение пароля
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil

}
