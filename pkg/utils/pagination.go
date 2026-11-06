package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"jianco-api/pkg/response"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

func ParsePagination(c fiber.Ctx) Pagination {
	page := parsePositiveInt(c.Query("page"), DefaultPage)
	perPage := parsePositiveInt(c.Query("per_page"), DefaultPerPage)
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

func (p Pagination) Meta(total int64) response.Pagination {
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(p.PerPage)))
	}

	return response.Pagination{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}
