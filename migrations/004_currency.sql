   CREATE TABLE currencies (
       code CHAR(3) PRIMARY KEY, -- Код валюты (например, USD)
       name VARCHAR(50) NOT NULL, -- Название валюты
       exchange_rate NUMERIC(18, 8) DEFAULT 1.0 -- Курс валюты по отношению к базовой валюте
   );