package main

import (
    "log"

    "api-gateway/config"
    "api-gateway/router"

    "github.com/gofiber/fiber/v2"
)

func main() {
    cfg := config.Load()
    app := fiber.New(fiber.Config{
        AppName: "GTBS API Gateway",
    })
    router.Setup(app, cfg)
    log.Fatal(app.Listen(":" + cfg.Port))
}
