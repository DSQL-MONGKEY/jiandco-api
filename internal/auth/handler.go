package auth

import (
	"github.com/gofiber/fiber/v3"

	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
	"jianco-api/pkg/utils"
)

type Handler struct {
	service *Service
	storage *storage.Client
}

func NewHandler(service *Service, storageClient *storage.Client) *Handler {
	return &Handler{service: service, storage: storageClient}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var request LoginRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	result, err := h.service.Login(c.Context(), LoginInput{
		Username:  request.Username,
		Password:  request.Password,
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
	})
	if err != nil {
		return err
	}

	return response.OK(c, result)
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	var request RefreshRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	result, err := h.service.Refresh(c.Context(), request.RefreshToken)
	if err != nil {
		return err
	}

	return response.OK(c, result)
}

func (h *Handler) Logout(c fiber.Ctx) error {
	session, ok := CurrentSession(c)
	if !ok {
		return apperror.Unauthorized("session not found")
	}

	if err := h.service.Logout(c.Context(), session.ID); err != nil {
		return err
	}

	return response.OK(c, map[string]any{"message": "logged out"})
}

func (h *Handler) Me(c fiber.Ctx) error {
	user, ok := CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("user not found")
	}

	result, err := h.service.Me(c.Context(), user)
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	user, ok := CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("user not found")
	}

	var request UpdateProfileRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	updated, err := h.service.UpdateProfile(c.Context(), user.ID, request.Name, request.Email, request.Phone)
	if err != nil {
		return err
	}

	return response.OK(c, ToUserResponse(updated, h.storage))
}

func (h *Handler) ChangePassword(c fiber.Ctx) error {
	user, ok := CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("user not found")
	}
	session, ok := CurrentSession(c)
	if !ok {
		return apperror.Unauthorized("session not found")
	}

	var request ChangePasswordRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	if err := h.service.ChangePassword(c.Context(), user.ID, session.ID, request.CurrentPassword, request.NewPassword); err != nil {
		return err
	}

	return response.OK(c, map[string]any{"message": "password updated"})
}

func (h *Handler) UploadAvatar(c fiber.Ctx) error {
	user, ok := CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("user not found")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequest("file is required")
	}

	result, err := h.service.UploadAvatar(c.Context(), user.ID, file)
	if err != nil {
		return err
	}

	return response.OK(c, result)
}
