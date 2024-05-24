package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
)

type Config struct {
	DBPort     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	ServerPort string
}

func initConfig() Config {
	return Config{
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBHost:     getEnv("DB_HOST"),
		DBName:     getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
	err := godotenv.Load()
	if os.Getenv("IS_LOCAL") != "true" {
		return os.Getenv(key)
	}
	if err != nil {
		return ""
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func DatabaseUrlBuilder() string {
	dbUser := initConfig().DBUser
	dbPassword := url.QueryEscape(initConfig().DBPassword)
	dbName := initConfig().DBName
	dbPort := initConfig().DBPort
	dbHost := initConfig().DBHost
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
