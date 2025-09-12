package api

import (
	"ecomm/configs"
	"ecomm/internal/api/rest"
	"ecomm/internal/api/rest/handlers"
	"github.com/gofiber/fiber/v3"
)

func StartServer(config configs.AppConfig) {
	//create a new fiber app
	app := fiber.New()

	//intantiate rest handler
	rh := &rest.RestHandler{
		App: app,
	}

	SetupRoutes(rh)
	//start the server
	app.Listen(config.ServerPort)
}

// function to setup all the routes
func SetupRoutes(rh *rest.RestHandler) {
	//setup user routes
	handlers.SetupUserRoutes(rh)
	//transaction routes
	//catalog routes
}
