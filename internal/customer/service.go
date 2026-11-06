package customer

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Customer, int64, error) {
	customers, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return customers, total, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("customer not found")
		}
		return nil, apperror.Internal(err)
	}
	return customer, nil
}

func (s *Service) Create(ctx context.Context, request Request) (*models.Customer, error) {
	customer := &models.Customer{
		Name:       strings.TrimSpace(request.Name),
		Phone:      strings.TrimSpace(request.Phone),
		Email:      strings.ToLower(strings.TrimSpace(request.Email)),
		Address:    strings.TrimSpace(request.Address),
		City:       strings.TrimSpace(request.City),
		Province:   strings.TrimSpace(request.Province),
		PostalCode: strings.TrimSpace(request.PostalCode),
		Notes:      strings.TrimSpace(request.Notes),
	}
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, apperror.Internal(err)
	}
	return customer, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, request Request) (*models.Customer, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}

	err := s.repo.Update(ctx, id, map[string]any{
		"name":        strings.TrimSpace(request.Name),
		"phone":       strings.TrimSpace(request.Phone),
		"email":       strings.ToLower(strings.TrimSpace(request.Email)),
		"address":     strings.TrimSpace(request.Address),
		"city":        strings.TrimSpace(request.City),
		"province":    strings.TrimSpace(request.Province),
		"postal_code": strings.TrimSpace(request.PostalCode),
		"notes":       strings.TrimSpace(request.Notes),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	orderCount, err := s.repo.CountOrders(ctx, id)
	if err != nil {
		return apperror.Internal(err)
	}
	if orderCount > 0 {
		return apperror.Conflict("customer has orders and cannot be deleted")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
