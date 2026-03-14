package middleware
import (
	"time"
	"api-gateway/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter(config config.Config) fiber.Handler {
	return limiter.New(
	limiter.Config{
		Max: config.RateLimitMax,
		Expiration: time.Duration(config.RateLimitWindow)*time.Second,
		KeyGenerator: func(context *fiber.Ctx) string{
			return context.IP()
		},
		
	}) 
}
