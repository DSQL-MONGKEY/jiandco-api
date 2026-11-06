package order

import (
	"testing"

	"jianco-api/internal/models"
)

func TestCanTransitionOrder(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{models.OrderStatusPending, models.OrderStatusConfirmed, true},
		{models.OrderStatusPending, models.OrderStatusCancelled, true},
		{models.OrderStatusPending, models.OrderStatusShipped, false},
		{models.OrderStatusConfirmed, models.OrderStatusProcessing, true},
		{models.OrderStatusProcessing, models.OrderStatusShipped, true},
		{models.OrderStatusShipped, models.OrderStatusCompleted, true},
		{models.OrderStatusShipped, models.OrderStatusCancelled, true},
		{models.OrderStatusCompleted, models.OrderStatusCancelled, false},
		{models.OrderStatusCancelled, models.OrderStatusPending, false},
		{models.OrderStatusPending, models.OrderStatusPending, true},
	}

	for _, test := range tests {
		if got := canTransitionOrder(test.from, test.to); got != test.want {
			t.Fatalf("canTransitionOrder(%s, %s) = %v, want %v", test.from, test.to, got, test.want)
		}
	}
}

func TestCanTransitionPayment(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{models.PaymentStatusUnpaid, models.PaymentStatusPaid, true},
		{models.PaymentStatusUnpaid, models.PaymentStatusRefunded, false},
		{models.PaymentStatusPaid, models.PaymentStatusRefunded, true},
		{models.PaymentStatusPaid, models.PaymentStatusUnpaid, true},
		{models.PaymentStatusRefunded, models.PaymentStatusPaid, true},
	}

	for _, test := range tests {
		if got := canTransitionPayment(test.from, test.to); got != test.want {
			t.Fatalf("canTransitionPayment(%s, %s) = %v, want %v", test.from, test.to, got, test.want)
		}
	}
}

func TestCanTransitionShipping(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{models.ShippingStatusPending, models.ShippingStatusPickedUp, true},
		{models.ShippingStatusPickedUp, models.ShippingStatusInTransit, true},
		{models.ShippingStatusInTransit, models.ShippingStatusDelivered, true},
		{models.ShippingStatusInTransit, models.ShippingStatusReturned, true},
		{models.ShippingStatusDelivered, models.ShippingStatusReturned, true},
		{models.ShippingStatusPending, models.ShippingStatusDelivered, false},
	}

	for _, test := range tests {
		if got := canTransitionShipping(test.from, test.to); got != test.want {
			t.Fatalf("canTransitionShipping(%s, %s) = %v, want %v", test.from, test.to, got, test.want)
		}
	}
}
