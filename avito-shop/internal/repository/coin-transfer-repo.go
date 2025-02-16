package repository

import (
	"database/sql"
	"fmt"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

type CoinTransferRepository interface {
	CreateTx(tx *sql.Tx, transfer *domain.CoinTransfer) error
	GetReceived(userID int64) ([]domain.CoinTransfer, error)
	GetSent(userID int64) ([]domain.CoinTransfer, error)
}

type coinTransferRepo struct {
	db *sql.DB
}

func NewCoinTransferRepository(db *sql.DB) CoinTransferRepository {
	return &coinTransferRepo{db: db}
}

func (r *coinTransferRepo) CreateTx(tx *sql.Tx, transfer *domain.CoinTransfer) error {
	err := tx.QueryRow(
		"INSERT INTO coin_transfers(from_user, to_user, amount) VALUES($1, $2, $3) RETURNING id",
		transfer.FromUserID, transfer.ToUserID, transfer.Amount,
	).Scan(&transfer.ID)
	if err != nil {
		return fmt.Errorf("не удалось создать перевод в транзакции: %w", err)
	}
	return nil
}

func (r *coinTransferRepo) GetReceived(userID int64) ([]domain.CoinTransfer, error) {
	rows, err := r.db.Query("SELECT id, from_user, to_user, amount FROM coin_transfers WHERE to_user=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []domain.CoinTransfer
	for rows.Next() {
		var transfer domain.CoinTransfer
		if err := rows.Scan(&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func (r *coinTransferRepo) GetSent(userID int64) ([]domain.CoinTransfer, error) {
	rows, err := r.db.Query("SELECT id, from_user, to_user, amount FROM coin_transfers WHERE from_user=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []domain.CoinTransfer
	for rows.Next() {
		var transfer domain.CoinTransfer
		if err := rows.Scan(&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}
