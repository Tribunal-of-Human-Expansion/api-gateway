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

	// Ingress + SPA use /api prefix (public booking API — JWT still required on /bookings/* below)
	if cfg.BookingURL != "" {
		app.All("/api/bookings", proxy.Forward(cfg, "booking", cfg.BookingURL))
		app.All("/api/bookings/*", proxy.Forward(cfg, "booking", cfg.BookingURL))
	}

	app.Use(middleware.Auth(cfg))

	if cfg.BookingURL != "" {
		app.All("/bookings/*", proxy.Forward(cfg, "booking", cfg.BookingURL))
	}
	if cfg.CompatibiltyServiceURL != "" {
		app.All("/compatibility/*", proxy.Forward(cfg, "compatibility", cfg.CompatibiltyServiceURL))
	}
	if cfg.UserServiceURL != "" {
		app.All("/users/*", proxy.Forward(cfg, "users", cfg.UserServiceURL))
	}
	if cfg.AuditServiceURL != "" {
		app.All("/audit/*", proxy.Forward(cfg, "audit", cfg.AuditServiceURL))
	}
	if cfg.RouteManagementURL != "" {
		app.All("/routes/*", proxy.Forward(cfg, "routes", cfg.RouteManagementURL))
	}
}
