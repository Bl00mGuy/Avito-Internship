package router

import (
	"github.com/gorilla/mux"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/delivery/http/handlers"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/delivery/http/middleware"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func SetupRouter(userService service.UserService, jwtSecret string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/auth", handlers.AuthHandler(userService, jwtSecret)).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware(jwtSecret))
	api.HandleFunc("/info", handlers.InfoHandler(userService)).Methods("GET")
	api.HandleFunc("/sendCoin", handlers.SendCoinHandler(userService)).Methods("POST")
	api.HandleFunc("/buy/{item}", handlers.BuyHandler(userService)).Methods("GET")

	return router
}
