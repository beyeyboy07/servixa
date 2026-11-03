package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	JWTSecret               string
	JWTAccessTokenDuration  time.Duration
	JWTRefreshTokenDuration time.Duration
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return Config{
		AppPort:    getEnv("APP_PORT", "8089"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "servixa"),
		JWTSecret: getEnv(
			"JWT_SECRET",
			"change-this-secret",
		),
		JWTAccessTokenDuration: getDurationEnv(
			"JWT_ACCESS_TOKEN_DURATION",
			15*time.Minute,
		),
		JWTRefreshTokenDuration: getDurationEnv(
			"JWT_REFRESH_TOKEN_DURATION",
			7*24*time.Hour,
		),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf(
			"invalid duration for %s: %v, using fallback %s",
			key,
			err,
			fallback,
		)

		return fallback
	}

	return duration
}
