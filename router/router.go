package router

import (
	"api-gateway/config"
	"api-gateway/middleware"
	"api-gateway/proxy"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, cfg *config.Config) {

	// Health check — no auth, registered first
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Use(middleware.RateLimiter(cfg))
	app.Use(middleware.Auth(cfg))
	app.All("/bookings/*", proxy.Forward(cfg, "booking", cfg.BookingURL))
	app.All("/compatibility/*", proxy.Forward(cfg, "compatibility", cfg.CompatibiltyServiceURL))
	app.All("/users/*", proxy.Forward(cfg, "users", cfg.UserServiceURL))
	app.All("/audit/*", proxy.Forward(cfg, "audit", cfg.AuditServiceURL))

}
