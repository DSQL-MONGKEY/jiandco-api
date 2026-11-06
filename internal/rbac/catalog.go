package rbac

import "jianco-api/internal/models"

const (
	PageOverview   = "overview"
	PageProducts   = "products"
	PageOrders     = "orders"
	PageCustomers  = "customers"
	PageAccounting = "accounting"
	PageUsers      = "users"
	PageRoles      = "roles"

	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

type PermissionDef struct {
	Key    string
	Page   string
	Action string
	Label  string
}

type RoleDef struct {
	Level       int
	Name        string
	Description string
}

var RoleCatalog = []RoleDef{
	{Level: models.RoleLevelDeveloper, Name: "Developer", Description: "Full access to every feature and system setting"},
	{Level: models.RoleLevelSuperAdmin, Name: "Super Admin", Description: "Full access to all operational features"},
	{Level: models.RoleLevelAdmin, Name: "Admin", Description: "Operational management with limited user administration"},
	{Level: models.RoleLevelStaff, Name: "Staff", Description: "Daily operational access"},
}

var PermissionCatalog = []PermissionDef{
	{Key: "overview.view", Page: PageOverview, Action: ActionView, Label: "View overview"},

	{Key: "products.view", Page: PageProducts, Action: ActionView, Label: "View products"},
	{Key: "products.create", Page: PageProducts, Action: ActionCreate, Label: "Create products"},
	{Key: "products.update", Page: PageProducts, Action: ActionUpdate, Label: "Update products"},
	{Key: "products.delete", Page: PageProducts, Action: ActionDelete, Label: "Delete products"},

	{Key: "orders.view", Page: PageOrders, Action: ActionView, Label: "View orders"},
	{Key: "orders.create", Page: PageOrders, Action: ActionCreate, Label: "Create orders"},
	{Key: "orders.update", Page: PageOrders, Action: ActionUpdate, Label: "Update orders"},
	{Key: "orders.delete", Page: PageOrders, Action: ActionDelete, Label: "Delete orders"},

	{Key: "customers.view", Page: PageCustomers, Action: ActionView, Label: "View customers"},
	{Key: "customers.create", Page: PageCustomers, Action: ActionCreate, Label: "Create customers"},
	{Key: "customers.update", Page: PageCustomers, Action: ActionUpdate, Label: "Update customers"},
	{Key: "customers.delete", Page: PageCustomers, Action: ActionDelete, Label: "Delete customers"},

	{Key: "accounting.view", Page: PageAccounting, Action: ActionView, Label: "View accounting"},
	{Key: "accounting.create", Page: PageAccounting, Action: ActionCreate, Label: "Create expenses"},
	{Key: "accounting.update", Page: PageAccounting, Action: ActionUpdate, Label: "Update expenses"},
	{Key: "accounting.delete", Page: PageAccounting, Action: ActionDelete, Label: "Delete expenses"},

	{Key: "users.view", Page: PageUsers, Action: ActionView, Label: "View users"},
	{Key: "users.create", Page: PageUsers, Action: ActionCreate, Label: "Create users"},
	{Key: "users.update", Page: PageUsers, Action: ActionUpdate, Label: "Update users"},
	{Key: "users.delete", Page: PageUsers, Action: ActionDelete, Label: "Delete users"},

	{Key: "roles.view", Page: PageRoles, Action: ActionView, Label: "View roles and permissions"},
	{Key: "roles.update", Page: PageRoles, Action: ActionUpdate, Label: "Update role permissions"},
}

var DefaultRolePermissions = map[int][]string{
	models.RoleLevelAdmin: {
		"overview.view",
		"products.view", "products.create", "products.update", "products.delete",
		"orders.view", "orders.create", "orders.update", "orders.delete",
		"customers.view", "customers.create", "customers.update", "customers.delete",
		"accounting.view", "accounting.create", "accounting.update", "accounting.delete",
		"users.view", "users.create", "users.update",
		"roles.view",
	},
	models.RoleLevelStaff: {
		"overview.view",
		"products.view",
		"orders.view", "orders.create", "orders.update",
		"customers.view", "customers.create", "customers.update",
		"accounting.view", "accounting.create",
	},
}
