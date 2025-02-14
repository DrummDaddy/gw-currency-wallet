   CREATE TABLE wallets (
       id SERIAL PRIMARY KEY,
       user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Ссылается на пользователя
       currency_code CHAR(3) NOT NULL, -- Код валюты (например, USD, EUR)
       balance NUMERIC(20, 2) DEFAULT 0.00, -- Баланс кошелька
       created_at TIMESTAMP DEFAULT NOW() -- Время создания
   );