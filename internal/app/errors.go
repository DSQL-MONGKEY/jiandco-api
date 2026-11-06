package app

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
)

func errorHandler(c fiber.Ctx, err error) error {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		if appErr.Status >= http.StatusInternalServerError {
			slog.Error("request failed",
				"request_id", c.GetRespHeader(fiber.HeaderXRequestID),
				"method", c.Method(),
				"path", c.Path(),
				"error", appErr.Error(),
			)
		}
		return response.Fail(c, appErr.Status, appErr.Message)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		message := fiberErr.Message
		if message == "" {
			message = http.StatusText(fiberErr.Code)
		}
		return response.Fail(c, fiberErr.Code, message)
	}

	slog.Error("unhandled request error",
		"request_id", c.GetRespHeader(fiber.HeaderXRequestID),
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)
	return response.Fail(c, http.StatusInternalServerError, "internal server error")
}
