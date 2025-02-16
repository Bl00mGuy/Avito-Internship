package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "shop"),
		JWTSecret:  getEnv("JWT_SECRET", "my_secret_key"),
		Port:       getEnv("PORT", "8080"),
	}
}

func LoadForTests() *Config {
	projectRoot, err := filepath.Abs("../")

	if err != nil {
		log.Fatalf("Ошибка определения корня проекта: %v", err)
	}

	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла по пути %s: %v", envPath, err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "shop"),
		JWTSecret:  getEnv("JWT_SECRET", "my_secret_key"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
