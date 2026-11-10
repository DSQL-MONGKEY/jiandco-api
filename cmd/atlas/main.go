package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"jianco-api/internal/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.AuthSession{},
		&models.Category{},
		&models.Size{},
		&models.Color{},
		&models.ProductColor{},
		&models.Product{},
		&models.ProductImage{},
		&models.ProductVariant{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.Expense{},
		&models.ExpenseItem{},
		&models.StockMovement{},
		&models.CatalogSetting{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(stmts)
}
