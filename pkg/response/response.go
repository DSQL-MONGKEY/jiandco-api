package response

import "github.com/gofiber/fiber/v3"

type Success struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type Error struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type List struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func OK(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Success{Success: true, Data: data})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Success{Success: true, Data: data})
}

func Fail(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Error{Success: false, Code: status, Message: message})
}
