package rest

import (
	"github.com/gofiber/fiber/v3"
)

// her we need to define pointer of fiber application
// so that we can register our routes with the fiber app instance
// and then we can call the setup user routes function using that struct instance
// so we created httpHandler.go
type RestHandler struct {
	App *fiber.App
}
