package auth

import (
	"github.com/gofiber/fiber/v3"

	"jianco-api/internal/models"
)

const (
	userKey    = "auth.user"
	sessionKey = "auth.session"
	claimsKey  = "auth.claims"
)

func SetUser(c fiber.Ctx, user *models.User) {
	c.SetLocal(userKey, user)
}

func CurrentUser(c fiber.Ctx) (*models.User, bool) {
	user, ok := c.Locals(userKey).(*models.User)
	return user, ok && user != nil
}

func SetSession(c fiber.Ctx, session *models.AuthSession) {
	c.SetLocal(sessionKey, session)
}

func CurrentSession(c fiber.Ctx) (*models.AuthSession, bool) {
	session, ok := c.Locals(sessionKey).(*models.AuthSession)
	return session, ok && session != nil
}

func SetClaims(c fiber.Ctx, claims *Claims) {
	c.SetLocal(claimsKey, claims)
}

func CurrentClaims(c fiber.Ctx) (*Claims, bool) {
	claims, ok := c.Locals(claimsKey).(*Claims)
	return claims, ok && claims != nil
}
