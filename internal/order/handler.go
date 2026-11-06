package order

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"jianco-api/internal/auth"
	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
	"jianco-api/pkg/utils"
)

type Handler struct {
	service  *Service
	location *time.Location
}

func NewHandler(service *Service, location *time.Location) *Handler {
	return &Handler{service: service, location: location}
}

func (h *Handler) List(c fiber.Ctx) error {
	pagination := utils.ParsePagination(c)
	filter := ListFilter{
		Page:           pagination.Page,
		PerPage:        pagination.PerPage,
		Search:         strings.TrimSpace(c.Query("search")),
		Status:         strings.TrimSpace(c.Query("status")),
		PaymentStatus:  strings.TrimSpace(c.Query("payment_status")),
		ShippingStatus: strings.TrimSpace(c.Query("shipping_status")),
	}
	if filter.Status != "" && !validOrderStatus(filter.Status) {
		return apperror.BadRequest("invalid status")
	}
	if filter.PaymentStatus != "" && !validPaymentStatus(filter.PaymentStatus) {
		return apperror.BadRequest("invalid payment_status")
	}
	if filter.ShippingStatus != "" && !validShippingStatus(filter.ShippingStatus) {
		return apperror.BadRequest("invalid shipping_status")
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		parsed, err := uuid.Parse(customerID)
		if err != nil {
			return apperror.BadRequest("invalid customer_id")
		}
		filter.CustomerID = &parsed
	}
	if from := c.Query("from"); from != "" {
		parsed, err := time.ParseInLocation("2006-01-02", from, h.location)
		if err != nil {
			return apperror.BadRequest("from must use YYYY-MM-DD format")
		}
		filter.From = &parsed
	}
	if to := c.Query("to"); to != "" {
		parsed, err := time.ParseInLocation("2006-01-02", to, h.location)
		if err != nil {
			return apperror.BadRequest("to must use YYYY-MM-DD format")
		}
		end := parsed.AddDate(0, 0, 1)
		filter.To = &end
	}

	orders, total, err := h.service.List(c.Context(), filter)
	if err != nil {
		return err
	}

	return response.OK(c, response.List{
		Items:      ToResponses(orders),
		Pagination: pagination.Meta(total),
	})
}

func (h *Handler) Create(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	var request CreateRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.Create(c.Context(), actor, request)
	if err != nil {
		return err
	}
	return response.Created(c, ToResponse(order))
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}
	order, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}

	var request UpdateRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.Update(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}
	if err := h.service.Delete(c.Context(), actor, id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "order deleted"})
}

func (h *Handler) UpdateItems(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}

	var request UpdateItemsRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.UpdateItems(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func (h *Handler) UpdateStatus(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}

	var request UpdateStatusRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.UpdateStatus(c.Context(), actor, id, request.Status)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func (h *Handler) UpdatePayment(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}

	var request UpdatePaymentRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.UpdatePayment(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func (h *Handler) UpdateShipping(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid order id")
	}

	var request UpdateShippingRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	order, err := h.service.UpdateShipping(c.Context(), actor, id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(order))
}

func validOrderStatus(value string) bool {
	switch value {
	case models.OrderStatusPending, models.OrderStatusConfirmed, models.OrderStatusProcessing,
		models.OrderStatusShipped, models.OrderStatusCompleted, models.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func validPaymentStatus(value string) bool {
	switch value {
	case models.PaymentStatusUnpaid, models.PaymentStatusPaid, models.PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func validShippingStatus(value string) bool {
	switch value {
	case models.ShippingStatusPending, models.ShippingStatusPickedUp, models.ShippingStatusInTransit,
		models.ShippingStatusDelivered, models.ShippingStatusReturned:
		return true
	default:
		return false
	}
}
