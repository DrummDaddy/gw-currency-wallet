package main

import (
	"database/sql"
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/app/handlers"
	"gw-currency-wallet/internal/auth"
	"gw-currency-wallet/internal/grpc"
	"gw-currency-wallet/internal/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/authenticate", handlers.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/balance", handlers.GetBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/deposit", handlers.DepositHandler(&sql.DB{}, &auth.JWTManager{})).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/wallet/withdraw", handlers.Withdraw).Methods(http.MethodPost)
	r.HandleFunc("swagger/", httpSwagger.WrapHandler)
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Проверка конфигурации
	if cfg.PostgresDSN == "" {
		log.Fatalf("Не указан DSN для подключения к PostgreSQL")
	}
	if cfg.GRPCPort == "" {
		log.Fatalf("Не указан порт для запуска gRPC сервера")
	}

	// Подключение к базе данных
	db, err := postgres.NewPostgresStorage(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer func() {
		log.Println("Закрываю соединение с базой данных...")
		db.Close()
	}()

	// Создание gRPC-сервера
	adapter := &postgres.PostgresAdapter{Storage: db}
	grpcServer, err := grpc.NewServer(cfg, adapter)

	if err != nil {
		log.Fatalf("Ошибка создания gRPC-сервера: %v", err)
	}

	// Канал для управления завершением работы сервиса
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск gRPC-сервера в отдельной горутине
	go func() {
		log.Printf("Запуск gRPC-сервера на порту %s...", cfg.GRPCPort)
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
		}
	}()

	log.Println("Сервис gw-currency-wallet успешно запущен.")

	// Ожидание сигнала завершения
	<-quit

	log.Println("Получен сигнал завершения, завершаю работу...")

}
