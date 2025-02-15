package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func BuyHandler(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)
		vars := mux.Vars(r)
		item := vars["item"]
		if item == "" {
			http.Error(w, "Не указан товар", http.StatusBadRequest)
			return
		}
		err := svc.BuyItem(userID, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Покупка прошла успешно"})
	}
}
