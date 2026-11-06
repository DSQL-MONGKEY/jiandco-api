package accounting

import (
	"context"
	"errors"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/utils"
)

type Service struct {
	repo     *Repository
	storage  *storage.Client
	location *time.Location
}

func NewService(repo *Repository, storageClient *storage.Client, location *time.Location) *Service {
	return &Service{repo: repo, storage: storageClient, location: location}
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Expense, int64, error) {
	var start, end *time.Time
	if strings.TrimSpace(filter.Month) != "" {
		monthStart, monthEnd, err := utils.ParseMonth(filter.Month, s.location)
		if err != nil {
			return nil, 0, apperror.BadRequest(err.Error())
		}
		start = &monthStart
		end = &monthEnd
	}

	expenses, total, err := s.repo.List(ctx, filter, start, end)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return expenses, total, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	expense, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("expense not found")
		}
		return nil, apperror.Internal(err)
	}
	return expense, nil
}

func (s *Service) Create(ctx context.Context, actor *models.User, request CreateRequest) (*models.Expense, error) {
	expenseDate, err := s.parseExpenseDate(request.ExpenseDate)
	if err != nil {
		return nil, err
	}
	items, total, err := buildItems(request.Items, request.TotalAmount)
	if err != nil {
		return nil, err
	}

	expense := &models.Expense{
		Title:         strings.TrimSpace(request.Title),
		Category:      request.Category,
		Vendor:        strings.TrimSpace(request.Vendor),
		ExpenseDate:   expenseDate,
		PaymentMethod: strings.TrimSpace(request.PaymentMethod),
		TotalAmount:   total,
		Notes:         strings.TrimSpace(request.Notes),
	}
	if actor != nil {
		expense.CreatedByID = &actor.ID
	}

	if err := s.repo.Create(ctx, expense, items); err != nil {
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, expense.ID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, request UpdateRequest) (*models.Expense, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	expenseDate, err := s.parseExpenseDate(request.ExpenseDate)
	if err != nil {
		return nil, err
	}
	items, total, err := buildItems(request.Items, request.TotalAmount)
	if err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, id, map[string]any{
		"title":          strings.TrimSpace(request.Title),
		"category":       request.Category,
		"vendor":         strings.TrimSpace(request.Vendor),
		"expense_date":   expenseDate,
		"payment_method": strings.TrimSpace(request.PaymentMethod),
		"total_amount":   total,
		"notes":          strings.TrimSpace(request.Notes),
	}, items)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) UploadReceipt(ctx context.Context, id uuid.UUID, fileHeader *multipart.FileHeader) (*ReceiptResponse, error) {
	if !s.storage.Enabled() {
		return nil, apperror.New(503, "file storage is not configured")
	}
	expense, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	upload, err := storage.OpenValidated(fileHeader, storage.ReceiptTypes, storage.MaxReceiptSize)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyFile):
			return nil, apperror.BadRequest("receipt file is required")
		case errors.Is(err, storage.ErrFileTooLarge):
			return nil, apperror.BadRequest("receipt must be at most 10 MB")
		case errors.Is(err, storage.ErrUnsupportedType):
			return nil, apperror.BadRequest("receipt must be JPEG, PNG, or PDF")
		default:
			return nil, apperror.BadRequest("failed to read uploaded receipt")
		}
	}
	defer upload.File.Close()

	key := storage.BuildKey("accounting/receipts", id.String(), "receipt"+upload.Extension)
	url, err := s.storage.Upload(ctx, key, upload.File, upload.ContentType, upload.Size)
	if err != nil {
		return nil, apperror.Wrap(502, "failed to upload receipt", err)
	}

	if err := s.repo.UpdateReceipt(ctx, id, key); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, apperror.Internal(err)
	}
	if expense.ReceiptKey != "" {
		_ = s.storage.Delete(ctx, expense.ReceiptKey)
	}

	return &ReceiptResponse{ReceiptKey: key, ReceiptURL: url}, nil
}

func (s *Service) DeleteReceipt(ctx context.Context, id uuid.UUID) error {
	expense, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if expense.ReceiptKey == "" {
		return apperror.NotFound("receipt not found")
	}
	if err := s.repo.UpdateReceipt(ctx, id, ""); err != nil {
		return apperror.Internal(err)
	}
	_ = s.storage.Delete(ctx, expense.ReceiptKey)
	return nil
}

func (s *Service) Overview(ctx context.Context, month string) (*OverviewResponse, error) {
	start, end, err := utils.ParseMonth(month, s.location)
	if err != nil {
		return nil, apperror.BadRequest(err.Error())
	}

	total, err := s.repo.SumExpenses(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	byCategory, err := s.repo.SumByCategory(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	recent, err := s.repo.Recent(ctx, start, end, 5)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &OverviewResponse{
		Month:          start.Format(utils.MonthLayout),
		TotalExpenses:  total,
		ByCategory:     byCategory,
		RecentExpenses: ToResponses(recent, s.storage),
	}, nil
}

func (s *Service) parseExpenseDate(value string) (time.Time, error) {
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(value), s.location)
	if err != nil {
		return time.Time{}, apperror.BadRequest("expense_date must use YYYY-MM-DD format")
	}
	return parsed, nil
}

func buildItems(requests []ItemRequest, manualTotal int64) ([]models.ExpenseItem, int64, error) {
	items := make([]models.ExpenseItem, 0, len(requests))
	var total int64

	for _, request := range requests {
		unit := strings.TrimSpace(request.Unit)
		if unit == "" {
			unit = "pcs"
		}
		subtotal := int64(math.Round(float64(request.UnitPrice) * request.Quantity))
		if subtotal < 0 {
			return nil, 0, apperror.BadRequest("expense item subtotal cannot be negative")
		}
		items = append(items, models.ExpenseItem{
			Name:      strings.TrimSpace(request.Name),
			Quantity:  request.Quantity,
			Unit:      unit,
			UnitPrice: request.UnitPrice,
			Subtotal:  subtotal,
		})
		total += subtotal
	}

	if len(items) == 0 {
		total = manualTotal
	}
	if total <= 0 {
		return nil, 0, apperror.BadRequest("total_amount must be greater than zero")
	}

	return items, total, nil
}
