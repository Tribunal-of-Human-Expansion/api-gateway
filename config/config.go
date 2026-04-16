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

	BookingURL             string
	CompatibiltyServiceURL string
	RouteServiceURL        string
	UserServiceURL         string
	AuditServiceURL        string
	RouteManagementURL     string
	AuthorityServiceURL    string

	BreakerMaxRequests uint32
	BreakerTimeout     uint32
	BreakerFailures    uint32
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Encountered error in loading the Dot env, please check if its there.")
	}
	jwt := getEnv("JWT_SECRET", "")
	if jwt == "" {
		jwt = getEnv("JWT_SECERET", "secret")
	}
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		JWTSecret:              jwt,
		RateLimitMax:           getEnvInt("RATE_LIMIT_MAX", 100),
		RateLimitWindow:        getEnvInt("RATE_LIMIT_WINDOW", 60),
		BookingURL:             getEnvAny([]string{"BOOKING_URL", "BOOKING_SERVICE_URL"}, ""),
		CompatibiltyServiceURL: getEnvAny([]string{"COMPT_SERV_URL", "COMPATIBILITY_SERVICE_URL"}, ""),
		RouteServiceURL:        getEnvAny([]string{"ROUTE_SERVICE_URL"}, ""),
		UserServiceURL:         getEnvAny([]string{"USER_SERVICE", "USER_SERVICE_URL"}, ""),
		AuditServiceURL:        getEnvAny([]string{"AUDIT_SERVICE_URL"}, ""),
    RouteManagementURL:     getEnv("ROUTE_MGMT_URL", ""),
		AuthorityServiceURL:    getEnvAny([]string{"AUTHORITY_SERVICE_URL", "ADMIN_AUTHORITY_SERVICE_URL"}, ""),
		BreakerMaxRequests:     uint32(getEnvInt("BREAKER_MAX_REQUESTS", 3)),
		BreakerTimeout:         uint32(getEnvInt("BREAKER_TIMEOUT", 10)),
		BreakerFailures:        uint32(getEnvInt("BREAKER_FAILURES", 5)),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
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
