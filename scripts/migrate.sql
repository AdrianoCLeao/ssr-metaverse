DROP USER IF EXISTS platable_user;
CREATE USER platable_user WITH PASSWORD 'senha123';

DROP DATABASE IF EXISTS platable_db;
CREATE DATABASE platable_db OWNER platable_user;

GRANT ALL PRIVILEGES ON DATABASE platable_db TO platable_user;


CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ingredients TEXT NOT NULL,
    instructions TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);