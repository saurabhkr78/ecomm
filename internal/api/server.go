package api

import (
	"github.com/gofiber/fiber/v3"
)

type AppConfig struct {
	ServerPort string
}

func startServer(config AppConfig) {
	//create a new fiber app
	app := fiber.New()
	//app is listening on port 9000
	app.Listen(config.ServerPort)

}
