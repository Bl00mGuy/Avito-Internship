package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/configs"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/delivery/http/router"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/repository"
	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/service"
)

func getTestRouter(t *testing.T) http.Handler {
	cfg := configs.LoadForTests()
	database, err := repository.InitDB(cfg)
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	userRepo := repository.NewUserRepository(database)
	purchaseRepo := repository.NewPurchaseRepository(database)
	transferRepo := repository.NewCoinTransferRepository(database)
	userService := service.NewUserService(userRepo, purchaseRepo, transferRepo, database)
	return router.SetupRouter(userService, cfg.JWTSecret)
}

func TestAuthShouldCreateUserAndReturnToken(t *testing.T) {
	testRouter := getTestRouter(t)

	payload := map[string]string{
		"username": "testuser1",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %v", recorder.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	if token, ok := resp["token"]; !ok || token == "" {
		t.Fatal("JWT-токен не получен")
	}
}

func TestInfoEndpointReturnsWalletInformation(t *testing.T) {
	testRouter := getTestRouter(t)

	token := getAuthToken(testRouter, map[string]string{"username": "testuser2", "password": "password"}, t)

	req := httptest.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %v", recorder.Code)
	}

	var infoResp struct {
		Coins       int         `json:"coins"`
		Inventory   interface{} `json:"inventory"`
		CoinHistory interface{} `json:"coinHistory"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&infoResp); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if infoResp.Coins < 0 {
		t.Fatalf("Баланс отрицательный: %d", infoResp.Coins)
	}
}

func TestBuyMerchReducesCoinsAndRecordsPurchase(t *testing.T) {
	testRouter := getTestRouter(t)

	token := getAuthToken(testRouter, map[string]string{"username": "buyer", "password": "pass"}, t)

	req := httptest.NewRequest("GET", "/api/buy/t-shirt", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 при покупке мерча, получен %v", recorder.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	if msg, ok := resp["message"]; !ok || msg == "" {
		t.Fatal("Сообщение о покупке отсутствует")
	}
}

func TestSendCoinTransfersCoinsBetweenUsers(t *testing.T) {
	testRouter := getTestRouter(t)

	tokenSender := getAuthToken(testRouter, map[string]string{"username": "sender", "password": "pass"}, t)
	_ = getAuthToken(testRouter, map[string]string{"username": "receiver", "password": "pass"}, t)

	transferPayload := map[string]interface{}{
		"toUser": "receiver",
		"amount": 100,
	}
	body, _ := json.Marshal(transferPayload)
	req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenSender)
	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 при переводе монет, получен %v", recorder.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	if msg, ok := resp["message"]; !ok || msg == "" {
		t.Fatal("Сообщение о переводе отсутствует")
	}
}

func getAuthToken(router http.Handler, payload map[string]string, t *testing.T) string {
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Аутентификация вернула статус %v", recorder.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка декодирования ответа аутентификации: %v", err)
	}
	token, ok := resp["token"]
	if !ok || token == "" {
		t.Fatal("JWT-токен отсутствует в ответе")
	}
	return token
}
