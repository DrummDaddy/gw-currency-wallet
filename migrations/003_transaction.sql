   CREATE TABLE transactions (
       id SERIAL PRIMARY KEY,
       wallet_id INT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE, -- Ссылается на кошелёк
       type VARCHAR(50) NOT NULL, -- Тип операции (DEPOSIT, WITHDRAWAL, TRANSFER)
       amount NUMERIC(20, 2) NOT NULL, -- Сумма транзакции
       description TEXT, -- Описание транзакции
       created_at TIMESTAMP DEFAULT NOW() -- Время проведения транзакции
   );