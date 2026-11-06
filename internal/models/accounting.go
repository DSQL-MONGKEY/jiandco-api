package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	Base
	Title         string        `gorm:"type:varchar(150);not null"`
	Category      string        `gorm:"type:varchar(50);not null;default:'operational';index:idx_expenses_category"`
	Vendor        string        `gorm:"type:varchar(150)"`
	ExpenseDate   time.Time     `gorm:"type:date;not null;index:idx_expenses_expense_date"`
	PaymentMethod string        `gorm:"type:varchar(30)"`
	TotalAmount   int64         `gorm:"type:bigint;not null;default:0;check:chk_expenses_total_amount,total_amount >= 0"`
	ReceiptKey    string        `gorm:"type:varchar(255)"`
	Notes         string        `gorm:"type:text"`
	CreatedByID   *uuid.UUID    `gorm:"type:uuid;index:idx_expenses_created_by"`
	CreatedBy     *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items         []ExpenseItem `gorm:"foreignKey:ExpenseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SoftDelete
}

type ExpenseItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ExpenseID uuid.UUID `gorm:"type:uuid;not null;index:idx_expense_items_expense_id"`
	Expense   Expense   `gorm:"foreignKey:ExpenseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      string    `gorm:"type:varchar(150);not null"`
	Quantity  float64   `gorm:"type:numeric(12,2);not null;default:1;check:chk_expense_items_quantity,quantity > 0"`
	Unit      string    `gorm:"type:varchar(20);default:'pcs'"`
	UnitPrice int64     `gorm:"type:bigint;not null;default:0;check:chk_expense_items_unit_price,unit_price >= 0"`
	Subtotal  int64     `gorm:"type:bigint;not null;default:0;check:chk_expense_items_subtotal,subtotal >= 0"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}
