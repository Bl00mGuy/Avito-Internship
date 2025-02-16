package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/configs"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/delivery/http/router"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/repository"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func main() {
	cfg := configs.Load()
	database, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с БД: %v", err)
		}
	}()

	userRepo := repository.NewUserRepository(database)
	purchaseRepo := repository.NewPurchaseRepository(database)
	transferRepo := repository.NewCoinTransferRepository(database)
	userService := service.NewUserService(userRepo, purchaseRepo, transferRepo, database)

	setupRouter := router.SetupRouter(userService, cfg.JWTSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервис запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, setupRouter))
}
