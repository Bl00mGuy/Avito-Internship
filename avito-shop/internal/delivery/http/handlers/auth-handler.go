package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/auth"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(svc service.UserService, jwtSecret string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			http.Error(writer, "Неверный запрос", http.StatusBadRequest)
			return
		}
		user, err := svc.Auth(req.Username, req.Password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := auth.GenerateJWT(user, jwtSecret)
		if err != nil {
			http.Error(writer, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(writer).Encode(AuthResponse{Token: token})
	}
}
