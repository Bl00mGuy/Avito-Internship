package repository

import (
	"database/sql"
	"time"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

type PurchaseRepository interface {
	Create(purchase *domain.Purchase) error
	GetByUserID(userID int64) ([]domain.Purchase, error)
}

type purchaseRepo struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) PurchaseRepository {
	return &purchaseRepo{db: db}
}

func (r *purchaseRepo) Create(purchase *domain.Purchase) error {
	return r.db.QueryRow("INSERT INTO purchases(user_id, item, price, created_at) VALUES($1,$2,$3,$4) RETURNING id",
		purchase.UserID, purchase.Item, purchase.Price, time.Now()).Scan(&purchase.ID)
}

func (r *purchaseRepo) GetByUserID(userID int64) ([]domain.Purchase, error) {
	rows, err := r.db.Query("SELECT id, user_id, item, price, created_at FROM purchases WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var purchases []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.Item, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}
