package domain

import "time"

type User struct {
	ID       int64
	Username string
	Password string
	Coins    int
}

type Purchase struct {
	ID        int64
	UserID    int64
	Item      string
	Price     int
	CreatedAt time.Time
}

type CoinTransfer struct {
	ID         int64
	FromUserID int64
	ToUserID   int64
	Amount     int
	CreatedAt  time.Time
}

type CoinHistory struct {
	Received []struct {
		FromUser string
		Amount   int
	}
	Sent []struct {
		ToUser string
		Amount int
	}
}
