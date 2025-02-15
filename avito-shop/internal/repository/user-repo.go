package repository

import (
	"database/sql"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
	UpdateCoins(id int64, delta int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE username=$1", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Coins)
	return u, err
}

func (r *userRepo) GetByID(id int64) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Username, &u.Password, &u.Coins)
	return u, err
}

func (r *userRepo) Create(user *domain.User) error {
	return r.db.QueryRow("INSERT INTO users(username, password, coins) VALUES($1,$2,1000) RETURNING id",
		user.Username, user.Password).Scan(&user.ID)
}

func (r *userRepo) UpdateCoins(id int64, delta int) error {
	_, err := r.db.Exec("UPDATE users SET coins = coins + $1 WHERE id=$2", delta, id)
	return err
}
