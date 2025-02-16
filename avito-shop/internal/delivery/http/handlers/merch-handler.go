package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func BuyHandler(svc service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID := request.Context().Value("userID").(int64)
		vars := mux.Vars(request)
		item := vars["item"]
		if item == "" {
			http.Error(writer, "Не указан товар", http.StatusBadRequest)
			return
		}
		err := svc.BuyItem(userID, item)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(writer).Encode(map[string]string{"message": "Покупка прошла успешно"})
	}
}
