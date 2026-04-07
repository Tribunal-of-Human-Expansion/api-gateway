package proxy

import (
	"api-gateway/config"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/sony/gobreaker"
)

var breakers = map[string]*gobreaker.CircuitBreaker{}

func newBreaker(name string, cfg *config.Config) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: cfg.BreakerMaxRequests,
		Interval:    0,
		Timeout:     time.Duration(cfg.BreakerTimeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= uint32(cfg.BreakerFailures)
		},
	})
}

func Forward(cfg *config.Config, serviceName, targetURL string) fiber.Handler {
	return ForwardPrefix(cfg, serviceName, targetURL, "", "")
}

func ForwardPrefix(cfg *config.Config, serviceName, targetURL, sourcePrefix, targetPrefix string) fiber.Handler {
	if _, exists := breakers[serviceName]; !exists {
		breakers[serviceName] = newBreaker(serviceName, cfg)
	}
	return func(c *fiber.Ctx) error {
		breaker := breakers[serviceName]
		_, err := breaker.Execute(func() (interface{}, error) {
			target := fmt.Sprintf("%s%s", strings.TrimRight(targetURL, "/"), rewritePath(c.OriginalURL(), sourcePrefix, targetPrefix))
			return nil, proxy.Do(c, target)
		})
		if err != nil {
			if err == gobreaker.ErrOpenState {
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"error": fmt.Sprintf("%s service is unavailable", serviceName),
				})
			}
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "upstream service error",
			})
		}
		return nil
	}
}

func rewritePath(originalURL, sourcePrefix, targetPrefix string) string {
	if sourcePrefix == "" || targetPrefix == "" {
		return originalURL
	}

	path := originalURL
	query := ""
	if queryStart := strings.Index(path, "?"); queryStart >= 0 {
		query = path[queryStart:]
		path = path[:queryStart]
	}

	if path == sourcePrefix {
		return targetPrefix + query
	}
	if strings.HasPrefix(path, sourcePrefix+"/") {
		return targetPrefix + strings.TrimPrefix(path, sourcePrefix) + query
	}
	return originalURL
}
