package main

import (
    "log"

    "api-gateway/config"
    "api-gateway/router"

    "github.com/gofiber/fiber/v2"
)

func main() {
    // Step 1: Load all environment config
    cfg := config.Load()

    // Step 2: Create the Fiber app
    app := fiber.New(fiber.Config{
        AppName: "GTBS API Gateway",
    })

    // Step 3: Register all routes and middleware
    router.Setup(app, cfg)

    // Step 4: Start listening
    log.Fatal(app.Listen(":" + cfg.Port))
}
