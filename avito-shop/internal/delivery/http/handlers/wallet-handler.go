package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   interface{} `json:"inventory"`
	CoinHistory interface{} `json:"coinHistory"`
}

func InfoHandler(svc service.UserService) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)
		user, err := svc.GetUser(userID)
		if err != nil {
			http.Error(responseWriter, "Пользователь не найден", http.StatusInternalServerError)
			return
		}
		history, purchases, err := svc.GetInfo(userID)
		if err != nil {
			http.Error(responseWriter, "Ошибка получения информации", http.StatusInternalServerError)
			return
		}
		resp := InfoResponse{
			Coins:       user.Coins,
			Inventory:   purchases,
			CoinHistory: history,
		}
		if err := json.NewEncoder(responseWriter).Encode(resp); err != nil {
			http.Error(responseWriter, "Ошибка при отправке ответа", http.StatusInternalServerError)
		}
	}
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func SendCoinHandler(svc service.UserService) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)
		var req SendCoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(responseWriter, "Неверный запрос", http.StatusBadRequest)
			return
		}
		if req.Amount <= 0 {
			http.Error(responseWriter, "Количество должно быть положительным", http.StatusBadRequest)
			return
		}
		err := svc.TransferCoins(userID, req.ToUser, req.Amount)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(responseWriter).Encode(map[string]string{"message": "Монеты переведены"}); err != nil {
			http.Error(responseWriter, "Ошибка при отправке ответа", http.StatusInternalServerError)
		}
	}
}
