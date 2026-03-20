package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	JWTSecret       string
	RateLimitMax    int
	RateLimitWindow int // seconds

	BookingURL          string
	CompatibiltyService string
	UserService         string // All the other sevices wil come here first

	BreakerMaxRequests uint32
	BreakerTimeout     uint32
	BreakerFailures    uint32
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Encountered error in loading the Dot env, please check if its there.")
	}
	return &Config{
		Port:                getEnv("PORT", "8080"),
		JWTSecret:           getEnv("JWT_SECERET", ""),
		RateLimitMax:        getEnvInt("RATE_LIMIT_MAX", 100),
		RateLimitWindow:     getEnvInt("RATE_LIMIT_WINDOW", 60),
		BookingURL:          getEnv("BOOKIN_URL", ""),
		CompatibiltyService: getEnv("COMPT_SERV_URL", ""),
		UserService:         getEnv("USER_SERVICE", ""),
		BreakerMaxRequests:  uint32(getEnvInt("BREAKER_MAX_REQUESTS", 3)),
		BreakerTimeout:      uint32(getEnvInt("BREAKER_TIMEOUT", 10)),
		BreakerFailures:     uint32(getEnvInt("BREAKER_FAILURES", 5)),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	var result int
	_, err := fmt.Sscanf(val, "%d", &result)
	if err != nil {
		fmt.Printf("Encountered an error : %v \n", err)
		return fallback
	}
	return result
}
