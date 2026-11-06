package order

import "jianco-api/internal/models"

var orderTransitions = map[string]map[string]bool{
	models.OrderStatusPending: {
		models.OrderStatusConfirmed: true,
		models.OrderStatusCancelled: true,
	},
	models.OrderStatusConfirmed: {
		models.OrderStatusProcessing: true,
		models.OrderStatusCancelled:  true,
	},
	models.OrderStatusProcessing: {
		models.OrderStatusShipped:   true,
		models.OrderStatusCancelled: true,
	},
	models.OrderStatusShipped: {
		models.OrderStatusCompleted: true,
		models.OrderStatusCancelled: true,
	},
}

var paymentTransitions = map[string]map[string]bool{
	models.PaymentStatusUnpaid: {
		models.PaymentStatusPaid: true,
	},
	models.PaymentStatusPaid: {
		models.PaymentStatusRefunded: true,
		models.PaymentStatusUnpaid:   true,
	},
	models.PaymentStatusRefunded: {
		models.PaymentStatusPaid: true,
	},
}

var shippingTransitions = map[string]map[string]bool{
	models.ShippingStatusPending: {
		models.ShippingStatusPickedUp: true,
	},
	models.ShippingStatusPickedUp: {
		models.ShippingStatusInTransit: true,
	},
	models.ShippingStatusInTransit: {
		models.ShippingStatusDelivered: true,
		models.ShippingStatusReturned:  true,
	},
	models.ShippingStatusDelivered: {
		models.ShippingStatusReturned: true,
	},
}

func canTransitionOrder(from, to string) bool {
	if from == to {
		return true
	}
	return orderTransitions[from][to]
}

func canTransitionPayment(from, to string) bool {
	if from == to {
		return true
	}
	return paymentTransitions[from][to]
}

func canTransitionShipping(from, to string) bool {
	if from == to {
		return true
	}
	return shippingTransitions[from][to]
}
