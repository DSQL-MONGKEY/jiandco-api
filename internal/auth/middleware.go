package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
)

func RequireAuth(repo *Repository, jwtManager *JWTManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return response.Fail(c, fiber.StatusUnauthorized, "authorization header is required")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return response.Fail(c, fiber.StatusUnauthorized, "authorization header must use Bearer token")
		}

		claims, err := jwtManager.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		session, err := repo.FindSessionByID(c.Context(), claims.SessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.Fail(c, fiber.StatusUnauthorized, "session not found")
			}
			return apperror.Internal(err)
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil || session.UserID != userID {
			return response.Fail(c, fiber.StatusUnauthorized, "invalid token subject")
		}
		if session.RevokedAt != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "session has been revoked")
		}
		if !time.Now().Before(session.ExpiresAt) {
			return response.Fail(c, fiber.StatusUnauthorized, "session has expired")
		}
		if !session.User.IsActive {
			return response.Fail(c, fiber.StatusForbidden, "user is inactive")
		}

		SetUser(c, &session.User)
		SetSession(c, session)
		SetClaims(c, claims)

		return c.Next()
	}
}
