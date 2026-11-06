package auth

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo       *Repository
	jwt        *JWTManager
	storage    *storage.Client
	refreshTTL time.Duration
}

func NewService(repo *Repository, jwtManager *JWTManager, storageClient *storage.Client, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		jwt:        jwtManager,
		storage:    storageClient,
		refreshTTL: refreshTTL,
	}
}

type LoginInput struct {
	Username  string
	Password  string
	UserAgent string
	IPAddress string
}

type AvatarResult struct {
	Key string `json:"avatar_key"`
	URL string `json:"avatar_url"`
}

func (s *Service) Me(ctx context.Context, user *models.User) (*MeResponse, error) {
	if user == nil {
		return nil, apperror.Unauthorized("user not found")
	}

	var (
		permissions []string
		err         error
	)
	if user.Role.Level <= models.RoleLevelSuperAdmin {
		permissions, err = s.repo.AllPermissionKeys(ctx)
	} else {
		permissions, err = s.repo.PermissionKeys(ctx, user.RoleID)
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &MeResponse{
		User:        ToUserResponse(user, s.storage),
		Permissions: permissions,
	}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := s.repo.FindUserByUsername(ctx, strings.ToLower(strings.TrimSpace(input.Username)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid username or password")
		}
		return nil, apperror.Internal(err)
	}

	if !user.IsActive {
		return nil, apperror.Forbidden("user is inactive")
	}

	if !CheckPassword(user.PasswordHash, input.Password) {
		return nil, apperror.Unauthorized("invalid username or password")
	}

	refreshToken, refreshHash, err := newRefreshToken()
	if err != nil {
		return nil, apperror.Internal(err)
	}

	now := time.Now()
	session := &models.AuthSession{
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		UserAgent:        truncate(input.UserAgent, 255),
		IPAddress:        truncate(input.IPAddress, 45),
		ExpiresAt:        now.Add(s.refreshTTL),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, apperror.Internal(err)
	}

	accessToken, _, err := s.jwt.Generate(user.ID, session.ID, user.Role.Level, user.Username)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	_ = s.repo.UpdateLastLogin(ctx, user.ID, now)
	user.LastLoginAt = &now

	return &AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
		User:         ToUserResponse(user, s.storage),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, apperror.Unauthorized("refresh token is required")
	}

	session, err := s.repo.FindSessionByHash(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid refresh token")
		}
		return nil, apperror.Internal(err)
	}

	now := time.Now()
	if session.RevokedAt != nil {
		return nil, apperror.Unauthorized("session has been revoked")
	}
	if !now.Before(session.ExpiresAt) {
		return nil, apperror.Unauthorized("session has expired")
	}
	if !session.User.IsActive {
		return nil, apperror.Forbidden("user is inactive")
	}

	newToken, newHash, err := newRefreshToken()
	if err != nil {
		return nil, apperror.Internal(err)
	}

	rotated, err := s.repo.RotateSession(ctx, session.ID, newHash, now)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if !rotated {
		return nil, apperror.Unauthorized("session has been revoked")
	}

	accessToken, _, err := s.jwt.Generate(session.UserID, session.ID, session.User.Role.Level, session.User.Username)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &AuthResponse{
		Token:        accessToken,
		RefreshToken: newToken,
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
		User:         ToUserResponse(&session.User, s.storage),
	}, nil
}

func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.repo.RevokeSession(ctx, sessionID, time.Now()); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email, phone string) (*models.User, error) {
	err := s.repo.UpdateProfile(ctx, userID, strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email)), strings.TrimSpace(phone))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("email is already used")
		}
		return nil, apperror.Internal(err)
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return user, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, currentSessionID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return apperror.Internal(err)
	}

	if !CheckPassword(user.PasswordHash, currentPassword) {
		return apperror.BadRequest("current password is incorrect")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return apperror.Internal(err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, hash); err != nil {
		return apperror.Internal(err)
	}

	if err := s.repo.RevokeAllSessionsExcept(ctx, userID, currentSessionID, time.Now()); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) (*AvatarResult, error) {
	if !s.storage.Enabled() {
		return nil, apperror.New(http.StatusServiceUnavailable, "file storage is not configured")
	}

	upload, err := storage.OpenValidated(fileHeader, storage.ImageTypes, storage.MaxAvatarSize)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyFile):
			return nil, apperror.BadRequest("avatar file is required")
		case errors.Is(err, storage.ErrFileTooLarge):
			return nil, apperror.BadRequest("avatar must be at most 2 MB")
		case errors.Is(err, storage.ErrUnsupportedType):
			return nil, apperror.BadRequest("avatar must be a JPEG, PNG, or WebP image")
		default:
			return nil, apperror.BadRequest("failed to read uploaded file")
		}
	}
	defer upload.File.Close()

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	key := storage.BuildKey("user/avatars", userID.String(), "avatar"+upload.Extension)
	url, err := s.storage.Upload(ctx, key, upload.File, upload.ContentType, upload.Size)
	if err != nil {
		return nil, apperror.Wrap(http.StatusBadGateway, "failed to upload avatar", err)
	}

	if err := s.repo.UpdateAvatar(ctx, userID, key); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, apperror.Internal(err)
	}

	if user.AvatarKey != "" {
		_ = s.storage.Delete(ctx, user.AvatarKey)
	}

	return &AvatarResult{Key: key, URL: url}, nil
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}
