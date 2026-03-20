package proxy

import (
	"api-gateway/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/sony/gobreaker"
	"time"
)

var breakers = map[string]*gobreaker.CircuitBreaker{}

func newBreaker(name string, cfg *config.Config) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: cfg.BreakerMaxRequests,
		Interval:    0,
		Timeout:     time.Duration(cfg.BreakerTimeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})
}

func Forward(cfg *config.Config, serviceName, targetURL string) fiber.Handler {
	if _, exists := breakers[serviceName]; !exists {
		breakers[serviceName] = newBreaker(serviceName, cfg)
	}
	return func(c *fiber.Ctx) error {
		breaker := breakers[serviceName]
		_, err := breaker.Execute(func() (interface{}, error) {
			target := fmt.Sprintf("%s%s", targetURL, c.OriginalURL())
			return nil, proxy.Do(c, target)
		})
		if err != nil {
			if err == gobreaker.ErrOpenState {
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"error": fmt.Sprintf("%s servic is unavavilable", serviceName),
				})
			}
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "upstream service error",
			})
		}
		return nil
	}
}
