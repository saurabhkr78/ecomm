package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func ErrorMessage(ctx fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func InternalError(ctx fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
}
func BadRequestError(ctx fiber.Ctx, message string) error {
	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
		"message": message,
	})
}
func SuccessResponse(ctx fiber.Ctx, msg string, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}
