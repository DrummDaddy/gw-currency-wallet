# gw-currency-wallet

GW Currency Wallet - это приложение, предоставляющее функционал для управления электронным кошельком с поддержкой операций через gRPC и REST API. Кошелек позволяет пользователям хранить баланс в нескольких валютах и выполнять операции пополнения счета.

---

## Функциональные возможности:
1. **GRPC-Сервер**:
   - Реализован сервер, который обрабатывает запросы через gRPC.
   - Управление балансами пользователей.

2. **Операции с аккаунтом**:
   - Пополнение кошелька пользователя с валидацией валюты.
   - Хранение баланса в нескольких валютах: `USD`, `RUB`, `EUR`.
   - Обмен валют

3. **Интеграция c PostgreSQL**:
   - Данные пользователей и кошельков управляются через базу данных.

4. **Логирование**:
   - Поддержка настраиваемого уровня логирования.

---

## Основные зависимости:

- **Go Modules**: Все модули управления зависимостями извлечены с использованием `go mod`.
- **PostgreSQL**: Хранение данных пользователей и их кошельков.
- **gRPC**: Поддержка gRPC для сетевого взаимодействия.
- **Viper**: Чтение конфигов из переменных окружения или конфигурационных файлов.

---

## Установка и запуск

### 1. Клонировать репозиторий:
```bash
git clone &lt;url-репозитория&gt;
cd &lt;название проекта&gt;  

### 2. Настройка конфигурации: 
Создайте файл config.env в корневой директории с содержимым: 
GRPC_PORT=50051                # Порт для gRPC сервера
POSTGRES_DSN=postgres://user:password@localhost:5432/walletdb # Строка подключения к PostgreSQL
LOG_LEVEL=info                 # Уровень логирования: debug, info, warn, error

### 3. Сборка и запуск проекта локально: 
Сборка: 
go build -o gw-currency-wallet main.go 

### 4. Запуск с помощью Docker: 
Проект поддерживает Docker для развертывания сервиса в контейнере 
Сборка Docker-образа: 
docker build -t gw-currency-wallet .

Запуск контейнера: 
docker run -d --env-file config.env -p 50051:50051 qw-currency-wallet 

### Структура проекта: 
.
├── config/                 # Настройки и конфигурация
│   └── config.go           # Загрузка конфигурации из конфиг-файла/переменных окружения
├── internal/
│   ├── grpc/               # Реализация GRPC функциональности
│   │   └── server.go       # GRPC-сервер с регистрацией сервисов
│   ├── handlers/           # HTTP Handlers (REST API)
│   │   └── deposit.go      # Обработка операций пополнения
│   ├── storages/           # Модели и хранилища данных
│       └── model.go        # Модели данных (User, Wallet)
└── Dockerfile              # Файл для создания Docker образа 

### Запуск тестов: 
go test ./... -v