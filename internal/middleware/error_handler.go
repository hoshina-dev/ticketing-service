package middleware

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/hoshina-dev/ticketing-service/internal/apperr"
	"github.com/hoshina-dev/ticketing-service/internal/dto"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		return c.Status(appErr.HTTPStatus).JSON(dto.ErrorResponse{Message: appErr.Message})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(dto.ErrorResponse{Message: fiberErr.Message})
	}

	slog.Error("unhandled error", "error", err.Error(), "path", c.Path())
	return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "internal server error"})
}
