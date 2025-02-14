package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Секретный ключ для JWT
var jwtSecret = []byte("your_secret_key")

// Mock - функция для проверки пользователя.
func mockAuthenticateUser(username, password string) bool {
	mockUsername := "admin"
	mockHashPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	//Проверка имени пользователя и пароля
	if username == mockUsername && bcrypt.CompareHashAndPassword(mockHashPassword, []byte(password)) == nil {
		return true
	}
	return false
}

// Генерация JWT токена
func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SigningString()
}

// LoginHandler Авторизация пользователя
// @Summary Авторизация
// @Description Авторизация пользователя с использованием логина и пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse "Успешная авторизация"
// @Failure 401 {object} ErrorResponse "Неверные данные для входа"
// @Failure 405 {string} string "Метод не дозволен"
// @Router /auth/login [post]

// Обработчик для авторизации
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	// Аутентификация пользователя

	if !mockAuthenticateUser(req.Username, req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid username or password"})
		return

	}

	//Генерация JWT токена
	token, err := generateJWT(req.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
