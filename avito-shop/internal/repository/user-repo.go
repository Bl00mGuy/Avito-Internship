package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
	UpdateCoins(id int64, delta int) error
	UpdateCoinsTx(tx *sql.Tx, id int64, delta int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("пользователь с именем %s не найден", username)
	}
	return user, err
}

func (r *userRepo) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("пользователь с ID %d не найден", id)
	}
	return user, err
}

func (r *userRepo) Create(user *domain.User) error {
	err := r.db.QueryRow("INSERT INTO users(username, password, coins) VALUES($1,$2,1000) RETURNING id",
		user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return nil
}

func (r *userRepo) UpdateCoins(id int64, delta int) error {
	_, err := r.db.Exec("UPDATE users SET coins = coins + $1 WHERE id=$2", delta, id)
	return err
}

func (r *userRepo) UpdateCoinsTx(tx *sql.Tx, id int64, delta int) error {
	_, err := tx.Exec("UPDATE users SET coins = coins + $1 WHERE id=$2", delta, id)
	return err
}
