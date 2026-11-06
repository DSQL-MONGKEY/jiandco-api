package models

const (
	RoleLevelDeveloper  = 0
	RoleLevelSuperAdmin = 1
	RoleLevelAdmin      = 2
	RoleLevelStaff      = 3
)

const (
	ProductStatusDraft    = "draft"
	ProductStatusActive   = "active"
	ProductStatusArchived = "archived"
)

const (
	OrderStatusPending    = "pending"
	OrderStatusConfirmed  = "confirmed"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

const (
	PaymentStatusUnpaid   = "unpaid"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

const (
	ShippingStatusPending   = "pending"
	ShippingStatusPickedUp  = "picked_up"
	ShippingStatusInTransit = "in_transit"
	ShippingStatusDelivered = "delivered"
	ShippingStatusReturned  = "returned"
)

const (
	ExpenseCategoryOperational = "operational"
	ExpenseCategoryProduction  = "production"
	ExpenseCategoryMarketing   = "marketing"
	ExpenseCategoryRestock     = "restock"
	ExpenseCategoryOther       = "other"
)

const (
	StockMovementTypeInitial    = "initial"
	StockMovementTypeAdjustment = "adjustment"
	StockMovementTypeOrder      = "order"
	StockMovementTypeCancel     = "cancel"
	StockMovementTypeReturn     = "return"
)

const (
	StockReferenceTypeProduct = "product"
	StockReferenceTypeOrder   = "order"
	StockReferenceTypeManual  = "manual"
)
