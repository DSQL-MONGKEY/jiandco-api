package rbac

import (
	"github.com/gofiber/fiber/v3"

	"jianco-api/internal/auth"
	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
)

func RequirePermission(repo *Repository, page, action string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := auth.CurrentUser(c)
		if !ok {
			return response.Fail(c, fiber.StatusUnauthorized, "authentication is required")
		}

		if user.Role.Level <= models.RoleLevelSuperAdmin {
			return c.Next()
		}

		allowed, err := repo.HasPermission(c.Context(), user.RoleID, page, action)
		if err != nil {
			return apperror.Internal(err)
		}
		if !allowed {
			return response.Fail(c, fiber.StatusForbidden, "you do not have permission to access this resource")
		}

		return c.Next()
	}
}
