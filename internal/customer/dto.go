package customer

import (
	"time"

	"github.com/google/uuid"

	"jianco-api/internal/models"
)

type Request struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Phone      string `json:"phone" validate:"required,max=30"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Address    string `json:"address"`
	City       string `json:"city" validate:"max=100"`
	Province   string `json:"province" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"max=10"`
	Notes      string `json:"notes"`
}

type ListFilter struct {
	Page    int
	PerPage int
	Search  string
}

type Response struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToResponse(customer *models.Customer) Response {
	return Response{
		ID:         customer.ID,
		Name:       customer.Name,
		Phone:      customer.Phone,
		Email:      customer.Email,
		Address:    customer.Address,
		City:       customer.City,
		Province:   customer.Province,
		PostalCode: customer.PostalCode,
		Notes:      customer.Notes,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}
}

func ToResponses(customers []models.Customer) []Response {
	result := make([]Response, 0, len(customers))
	for index := range customers {
		result = append(result, ToResponse(&customers[index]))
	}
	return result
}
