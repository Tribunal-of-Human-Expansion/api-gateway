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
	app.Use(middleware.Redirector());

	if cfg.BookingURL != "" {
		registerProxy(app, "/bookings", proxy.ForwardPrefix(cfg, "booking", cfg.BookingURL, "/bookings", "/api/bookings"))
	}
	if cfg.CompatibiltyServiceURL != "" {
		registerProxy(app, "/compatibility", proxy.ForwardPrefix(cfg, "compatibility", cfg.CompatibiltyServiceURL, "/compatibility", "/api/v1"))
	}
	if cfg.RouteServiceURL != "" {
		registerProxy(app, "/routes", proxy.ForwardPrefix(cfg, "routes", cfg.RouteServiceURL, "/routes", "/api/routes"))
	}
	if cfg.UserServiceURL != "" {
		registerProxy(app, "/users", proxy.Forward(cfg, "users", cfg.UserServiceURL))
		registerProxy(app, "/notifications", proxy.Forward(cfg, "notifications", cfg.UserServiceURL))
	}
	if cfg.AuditServiceURL != "" {
		registerProxy(app, "/audit", proxy.ForwardPrefix(cfg, "audit", cfg.AuditServiceURL, "/audit", "/api/v1/audit"))
	}
	if cfg.AuthorityServiceURL != "" {
		registerProxy(app, "/authority", proxy.Forward(cfg, "authority", cfg.AuthorityServiceURL))
	}
	if cfg.RouteManagementURL != "" {
		app.All("/routes/*", proxy.Forward(cfg, "routes", cfg.RouteManagementURL))
	}
}

func registerProxy(app *fiber.App, prefix string, handler fiber.Handler) {
	app.All(prefix, handler)
	app.All(prefix+"/*", handler)
}
