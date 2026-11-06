package accounting

import (
	"time"

	"github.com/google/uuid"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
)

type ItemRequest struct {
	Name      string  `json:"name" validate:"required,max=150"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	Unit      string  `json:"unit" validate:"max=20"`
	UnitPrice int64   `json:"unit_price" validate:"min=0"`
}

type CreateRequest struct {
	Title         string        `json:"title" validate:"required,min=2,max=150"`
	Category      string        `json:"category" validate:"required,oneof=operational production marketing restock other"`
	Vendor        string        `json:"vendor" validate:"max=150"`
	ExpenseDate   string        `json:"expense_date" validate:"required"`
	PaymentMethod string        `json:"payment_method" validate:"max=30"`
	Notes         string        `json:"notes"`
	Items         []ItemRequest `json:"items" validate:"omitempty,dive"`
	TotalAmount   int64         `json:"total_amount" validate:"min=0"`
}

type UpdateRequest struct {
	Title         string        `json:"title" validate:"required,min=2,max=150"`
	Category      string        `json:"category" validate:"required,oneof=operational production marketing restock other"`
	Vendor        string        `json:"vendor" validate:"max=150"`
	ExpenseDate   string        `json:"expense_date" validate:"required"`
	PaymentMethod string        `json:"payment_method" validate:"max=30"`
	Notes         string        `json:"notes"`
	Items         []ItemRequest `json:"items" validate:"omitempty,dive"`
	TotalAmount   int64         `json:"total_amount" validate:"min=0"`
}

type ListFilter struct {
	Page     int
	PerPage  int
	Search   string
	Category string
	Month    string
}

type ItemResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Quantity  float64   `json:"quantity"`
	Unit      string    `json:"unit"`
	UnitPrice int64     `json:"unit_price"`
	Subtotal  int64     `json:"subtotal"`
}

type Response struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	Category      string         `json:"category"`
	Vendor        string         `json:"vendor"`
	ExpenseDate   string         `json:"expense_date"`
	PaymentMethod string         `json:"payment_method"`
	TotalAmount   int64          `json:"total_amount"`
	ReceiptKey    string         `json:"receipt_key"`
	ReceiptURL    string         `json:"receipt_url"`
	Notes         string         `json:"notes"`
	CreatedByID   *uuid.UUID     `json:"created_by_id"`
	Items         []ItemResponse `json:"items"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type CategoryTotal struct {
	Category string `json:"category"`
	Total    int64  `json:"total"`
}

type OverviewResponse struct {
	Month          string          `json:"month"`
	TotalExpenses  int64           `json:"total_expenses"`
	ByCategory     []CategoryTotal `json:"by_category"`
	RecentExpenses []Response      `json:"recent_expenses"`
}

type ReceiptResponse struct {
	ReceiptKey string `json:"receipt_key"`
	ReceiptURL string `json:"receipt_url"`
}

func ToResponse(expense *models.Expense, storageClient *storage.Client) Response {
	items := make([]ItemResponse, 0, len(expense.Items))
	for index := range expense.Items {
		item := expense.Items[index]
		items = append(items, ItemResponse{
			ID:        item.ID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Unit:      item.Unit,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return Response{
		ID:            expense.ID,
		Title:         expense.Title,
		Category:      expense.Category,
		Vendor:        expense.Vendor,
		ExpenseDate:   expense.ExpenseDate.Format("2006-01-02"),
		PaymentMethod: expense.PaymentMethod,
		TotalAmount:   expense.TotalAmount,
		ReceiptKey:    expense.ReceiptKey,
		ReceiptURL:    storageClient.URL(expense.ReceiptKey),
		Notes:         expense.Notes,
		CreatedByID:   expense.CreatedByID,
		Items:         items,
		CreatedAt:     expense.CreatedAt,
		UpdatedAt:     expense.UpdatedAt,
	}
}

func ToResponses(expenses []models.Expense, storageClient *storage.Client) []Response {
	result := make([]Response, 0, len(expenses))
	for index := range expenses {
		result = append(result, ToResponse(&expenses[index], storageClient))
	}
	return result
}
