package repository

import (
	"database/sql"
	"time"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

type CoinTransferRepository interface {
	Create(transfer *domain.CoinTransfer) error
	GetReceived(userID int64) ([]domain.CoinTransfer, error)
	GetSent(userID int64) ([]domain.CoinTransfer, error)
}

type coinTransferRepo struct {
	db *sql.DB
}

func NewCoinTransferRepository(db *sql.DB) CoinTransferRepository {
	return &coinTransferRepo{db: db}
}

func (r *coinTransferRepo) Create(transfer *domain.CoinTransfer) error {
	return r.db.QueryRow("INSERT INTO coin_transfers(from_user, to_user, amount, created_at) VALUES($1,$2,$3,$4) RETURNING id",
		transfer.FromUserID, transfer.ToUserID, transfer.Amount, time.Now()).Scan(&transfer.ID)
}

func (r *coinTransferRepo) GetReceived(userID int64) ([]domain.CoinTransfer, error) {
	rows, err := r.db.Query("SELECT id, from_user, to_user, amount, created_at FROM coin_transfers WHERE to_user=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transfers []domain.CoinTransfer
	for rows.Next() {
		var t domain.CoinTransfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)
	}
	return transfers, nil
}

func (r *coinTransferRepo) GetSent(userID int64) ([]domain.CoinTransfer, error) {
	rows, err := r.db.Query("SELECT id, from_user, to_user, amount, created_at FROM coin_transfers WHERE from_user=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transfers []domain.CoinTransfer
	for rows.Next() {
		var t domain.CoinTransfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)
	}
	return transfers, nil
}
