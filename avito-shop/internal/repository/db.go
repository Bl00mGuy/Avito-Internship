package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/configs"
)

func InitDB(cfg *configs.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE,
	password TEXT,
	coins INTEGER
);
CREATE TABLE IF NOT EXISTS purchases (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id),
	item TEXT,
	price INTEGER,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS coin_transfers (
	id SERIAL PRIMARY KEY,
	from_user INTEGER REFERENCES users(id),
	to_user INTEGER REFERENCES users(id),
	amount INTEGER,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err = database.Exec(schema)
	return database, err
}
