package middleware

import (
    "fmt"
    "strings"

    "github.com/gofiber/fiber/v2"
)

func Redirector() fiber.Handler {
    return func(c *fiber.Ctx) error {

        // Preserve body for downstream proxy — must do this before BodyParser
        rawBody := c.Body()

        routeID := c.Query("routeId")

        if routeID == "" {
            var body map[string]interface{}
            if err := c.BodyParser(&body); err == nil {
                if val, ok := body["routeId"].(string); ok {
                    routeID = val
                }
            }
            // Restore body so proxy.Do can still forward it
            c.Request().SetBody(rawBody)
        }

        if routeID != "" {
            parts := strings.Split(routeID, "-")
            // Expected format: ROUTE-SOURCE-DESTINATION
            if len(parts) >= 3 {
                source := parts[1]
                fmt.Printf("[GTBS Gateway] routeId=%s | source=%s | path=%s\n",
                    routeID, source, c.Path())
            } else {
                fmt.Printf("[GTBS Gateway] malformed routeId=%s | path=%s\n",
                    routeID, c.Path())
            }
        }

        return c.Next()
    }
}
