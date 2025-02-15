package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/configs"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/handlers"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/middleware"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/repository"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func main() {
	cfg := configs.Load()
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	transferRepo := repository.NewCoinTransferRepository(db)
	userService := service.NewUserService(userRepo, purchaseRepo, transferRepo)

	r := mux.NewRouter()
	r.HandleFunc("/api/auth", handlers.AuthHandler(userService, cfg.JWTSecret)).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	api.HandleFunc("/info", handlers.InfoHandler(userService)).Methods("GET")
	api.HandleFunc("/sendCoin", handlers.SendCoinHandler(userService)).Methods("POST")
	api.HandleFunc("/buy/{item}", handlers.BuyHandler(userService)).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервис запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
