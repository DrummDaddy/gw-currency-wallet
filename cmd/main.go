package main

import (
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/grpc"
	"gw-currency-wallet/internal/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
