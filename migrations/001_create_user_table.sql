CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL, -- Уникальное имя пользователя
    email VARCHAR(100) UNIQUE NOT NULL,  -- Уникальный email
    password_hash TEXT NOT NULL,         -- Хэш пароля
    created_at TIMESTAMP DEFAULT NOW()   -- Время создания
);