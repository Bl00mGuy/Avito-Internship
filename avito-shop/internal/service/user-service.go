package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/repository"
)

var merchCatalog = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

type UserService interface {
	Auth(username, password string) (*domain.User, error)
	GetUser(userID int64) (*domain.User, error)
	GetInfo(userID int64) (domain.CoinHistory, []domain.Purchase, error)
	TransferCoins(fromUserID int64, toUsername string, amount int) error
	BuyItem(userID int64, item string) error
}

type userService struct {
	userRepo     repository.UserRepository
	purchaseRepo repository.PurchaseRepository
	transferRepo repository.CoinTransferRepository
	db           *sql.DB
}

func NewUserService(userRepo repository.UserRepository, purchaseRepo repository.PurchaseRepository, transferRepo repository.CoinTransferRepository) UserService {
	return &userService{
		userRepo:     userRepo,
		purchaseRepo: purchaseRepo,
		transferRepo: transferRepo,
	}
}

func (s *userService) Auth(username, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		user = &domain.User{Username: username, Password: password, Coins: 1000}
		err = s.userRepo.Create(user)
		if err != nil {
			return nil, err
		}
	} else {
		if user.Password != password {
			return nil, errors.New("неверный пароль")
		}
	}
	return user, nil
}

func (s *userService) GetUser(userID int64) (*domain.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) GetInfo(userID int64) (domain.CoinHistory, []domain.Purchase, error) {
	history := domain.CoinHistory{}
	received, err := s.transferRepo.GetReceived(userID)
	if err != nil {
		return history, nil, err
	}
	sent, err := s.transferRepo.GetSent(userID)
	if err != nil {
		return history, nil, err
	}
	for _, r := range received {
		history.Received = append(history.Received, struct {
			FromUser string
			Amount   int
		}{FromUser: "User#" + itoa(r.FromUserID), Amount: r.Amount})
	}
	for _, t := range sent {
		history.Sent = append(history.Sent, struct {
			ToUser string
			Amount int
		}{ToUser: "User#" + itoa(t.ToUserID), Amount: t.Amount})
	}
	purchases, err := s.purchaseRepo.GetByUserID(userID)
	return history, purchases, err
}

func (s *userService) TransferCoins(fromUserID int64, toUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("количество монет должно быть положительным")
	}
	fromUser, err := s.userRepo.GetByID(fromUserID)
	if err != nil {
		return err
	}
	if fromUser.Coins < amount {
		return errors.New("недостаточно монет")
	}
	toUser, err := s.userRepo.GetByUsername(toUsername)
	if err != nil {
		return errors.New("получатель не найден")
	}
	err = s.userRepo.UpdateCoins(fromUserID, -amount)
	if err != nil {
		return err
	}
	err = s.userRepo.UpdateCoins(toUser.ID, amount)
	if err != nil {
		return err
	}
	transfer := &domain.CoinTransfer{
		FromUserID: fromUserID,
		ToUserID:   toUser.ID,
		Amount:     amount,
	}
	return s.transferRepo.Create(transfer)
}

func (s *userService) BuyItem(userID int64, item string) error {
	price, ok := merchCatalog[item]
	if !ok {
		return errors.New("товар не найден")
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.Coins < price {
		return errors.New("недостаточно монет для покупки")
	}
	err = s.userRepo.UpdateCoins(userID, -price)
	if err != nil {
		return err
	}
	purchase := &domain.Purchase{
		UserID: userID,
		Item:   item,
		Price:  price,
	}
	return s.purchaseRepo.Create(purchase)
}

func itoa(num int64) string {
	return fmt.Sprintf("%d", num)
}
