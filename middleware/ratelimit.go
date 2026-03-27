package middleware

import (
	"api-gateway/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func RateLimiter(cfg *config.Config) fiber.Handler {
	return limiter.New(
		limiter.Config{
			Max:        cfg.RateLimitMax,
			Expiration: time.Duration(cfg.RateLimitWindow) * time.Second,
			KeyGenerator: func(context *fiber.Ctx) string {
				return context.IP()
			},
			LimitReached: func(context *fiber.Ctx) error {
				return context.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests"})
			},
		})
}
