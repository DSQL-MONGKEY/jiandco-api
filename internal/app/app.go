package app

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"gorm.io/gorm"

	"jianco-api/config"
	"jianco-api/internal/accounting"
	"jianco-api/internal/auth"
	"jianco-api/internal/catalog"
	"jianco-api/internal/customer"
	"jianco-api/internal/order"
	"jianco-api/internal/overview"
	"jianco-api/internal/product"
	"jianco-api/internal/rbac"
	"jianco-api/internal/storage"
	"jianco-api/internal/user"
	"jianco-api/pkg/response"
)

type App struct {
	Fiber *fiber.App
}

func New(cfg *config.Config, db *gorm.DB, storageClient *storage.Client) *App {
	fiberApp := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		BodyLimit:    64 * 1024 * 1024,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: errorHandler,
	})

	fiberApp.Use(recover.New())
	fiberApp.Use(requestid.New())
	fiberApp.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${latency} ${method} ${path} ${ip}\n",
	}))
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	fiberApp.Use(helmet.New())
	fiberApp.Use(compress.New())

	authRepo := auth.NewRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTTL)
	authService := auth.NewService(authRepo, jwtManager, storageClient, cfg.RefreshTTL)
	authHandler := auth.NewHandler(authService, storageClient)

	rbacRepo := rbac.NewRepository(db)
	rbacService := rbac.NewService(rbacRepo)
	rbacHandler := rbac.NewHandler(rbacService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, storageClient)

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo, storageClient)
	productHandler := product.NewHandler(productService, storageClient)

	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, customerRepo, cfg.Location)
	orderHandler := order.NewHandler(orderService, cfg.Location)

	accountingRepo := accounting.NewRepository(db)
	accountingService := accounting.NewService(accountingRepo, storageClient, cfg.Location)
	accountingHandler := accounting.NewHandler(accountingService, storageClient)

	overviewRepo := overview.NewRepository(db)
	overviewService := overview.NewService(overviewRepo, cfg.Location)
	overviewHandler := overview.NewHandler(overviewService)

	catalogRepo := catalog.NewRepository(db)
	catalogService := catalog.NewService(catalogRepo, storageClient)
	catalogHandler := catalog.NewHandler(catalogService, storageClient)

	requireAuth := auth.RequireAuth(authRepo, jwtManager)
	loginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return response.Fail(c, fiber.StatusTooManyRequests, "too many requests, please try again later")
		},
	})

	fiberApp.Get("/healthz", func(c fiber.Ctx) error {
		return response.OK(c, map[string]any{"status": "ok"})
	})
	fiberApp.Get("/readyz", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return response.Fail(c, fiber.StatusServiceUnavailable, "database is not ready")
		}
		if err := sqlDB.PingContext(c.Context()); err != nil {
			return response.Fail(c, fiber.StatusServiceUnavailable, "database is not ready")
		}
		return response.OK(c, map[string]any{"status": "ready"})
	})

	api := fiberApp.Group("/api/v1")

	api.Post("/auth/login", loginLimiter, authHandler.Login)
	api.Post("/auth/refresh", loginLimiter, authHandler.Refresh)

	api.Get("/public/catalog", catalogHandler.PublicCatalog)
	api.Get("/public/catalog/products/:slug", catalogHandler.PublicProduct)
	api.Get("/public/media/*", catalogHandler.PublicMedia)

	authGroup := api.Group("/auth", requireAuth)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/me", authHandler.Me)
	authGroup.Put("/profile", authHandler.UpdateProfile)
	authGroup.Put("/password", authHandler.ChangePassword)
	authGroup.Post("/avatar", authHandler.UploadAvatar)

	roleGroup := api.Group("/roles", requireAuth)
	roleGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageRoles, rbac.ActionView), rbacHandler.ListRoles)
	roleGroup.Get("/:id/permissions", rbac.RequirePermission(rbacRepo, rbac.PageRoles, rbac.ActionView), rbacHandler.GetRolePermissions)
	roleGroup.Put("/:id/permissions", rbac.RequirePermission(rbacRepo, rbac.PageRoles, rbac.ActionUpdate), rbacHandler.UpdateRolePermissions)

	permissionGroup := api.Group("/permissions", requireAuth)
	permissionGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageRoles, rbac.ActionView), rbacHandler.ListPermissions)

	userGroup := api.Group("/users", requireAuth)
	userGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionView), userHandler.List)
	userGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionCreate), userHandler.Create)
	userGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionView), userHandler.Get)
	userGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionUpdate), userHandler.Update)
	userGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionDelete), userHandler.Delete)
	userGroup.Put("/:id/password", rbac.RequirePermission(rbacRepo, rbac.PageUsers, rbac.ActionUpdate), userHandler.ResetPassword)

	categoryGroup := api.Group("/categories", requireAuth)
	categoryGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListCategories)
	categoryGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionCreate), productHandler.CreateCategory)
	categoryGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.GetCategory)
	categoryGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateCategory)
	categoryGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionDelete), productHandler.DeleteCategory)

	sizeGroup := api.Group("/sizes", requireAuth)
	sizeGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListSizes)
	sizeGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionCreate), productHandler.CreateSize)
	sizeGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.GetSize)
	sizeGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateSize)
	sizeGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionDelete), productHandler.DeleteSize)

	colorGroup := api.Group("/colors", requireAuth)
	colorGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListColors)
	colorGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionCreate), productHandler.CreateColor)
	colorGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.GetColor)
	colorGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateColor)
	colorGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionDelete), productHandler.DeleteColor)

	productGroup := api.Group("/products", requireAuth)
	productGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListProducts)
	productGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionCreate), productHandler.CreateProduct)
	productGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.GetProduct)
	productGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateProduct)
	productGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionDelete), productHandler.DeleteProduct)
	productGroup.Get("/:id/colors", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListProductColors)
	productGroup.Post("/:id/colors", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.CreateProductColor)
	productGroup.Put("/:id/colors/:colorId", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateProductColor)
	productGroup.Delete("/:id/colors/:colorId", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.DeleteProductColor)
	productGroup.Post("/:id/colors/:colorId/image", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UploadProductColorImage)
	productGroup.Delete("/:id/colors/:colorId/image", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.DeleteProductColorImage)
	productGroup.Post("/:id/images", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UploadImages)
	productGroup.Delete("/:id/images/:imageId", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.DeleteImage)
	productGroup.Put("/:id/images/:imageId/primary", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.SetPrimaryImage)
	productGroup.Post("/:id/variants", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionCreate), productHandler.CreateVariant)
	productGroup.Put("/:id/variants/:variantId", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateVariant)
	productGroup.Delete("/:id/variants/:variantId", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionDelete), productHandler.DeleteVariant)
	productGroup.Patch("/:id/variants/:variantId/stock", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), productHandler.UpdateStock)
	productGroup.Get("/:id/variants/:variantId/movements", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), productHandler.ListStockMovements)

	customerGroup := api.Group("/customers", requireAuth)
	customerGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageCustomers, rbac.ActionView), customerHandler.List)
	customerGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageCustomers, rbac.ActionCreate), customerHandler.Create)
	customerGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageCustomers, rbac.ActionView), customerHandler.Get)
	customerGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageCustomers, rbac.ActionUpdate), customerHandler.Update)
	customerGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageCustomers, rbac.ActionDelete), customerHandler.Delete)

	orderGroup := api.Group("/orders", requireAuth)
	orderGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionView), orderHandler.List)
	orderGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionCreate), orderHandler.Create)
	orderGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionView), orderHandler.Get)
	orderGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionUpdate), orderHandler.Update)
	orderGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionDelete), orderHandler.Delete)
	orderGroup.Put("/:id/items", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionUpdate), orderHandler.UpdateItems)
	orderGroup.Patch("/:id/status", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionUpdate), orderHandler.UpdateStatus)
	orderGroup.Patch("/:id/payment", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionUpdate), orderHandler.UpdatePayment)
	orderGroup.Patch("/:id/shipping", rbac.RequirePermission(rbacRepo, rbac.PageOrders, rbac.ActionUpdate), orderHandler.UpdateShipping)

	expenseGroup := api.Group("/expenses", requireAuth)
	expenseGroup.Get("", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionView), accountingHandler.List)
	expenseGroup.Post("", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionCreate), accountingHandler.Create)
	expenseGroup.Get("/:id", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionView), accountingHandler.Get)
	expenseGroup.Put("/:id", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionUpdate), accountingHandler.Update)
	expenseGroup.Delete("/:id", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionDelete), accountingHandler.Delete)
	expenseGroup.Post("/:id/receipt", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionUpdate), accountingHandler.UploadReceipt)
	expenseGroup.Delete("/:id/receipt", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionUpdate), accountingHandler.DeleteReceipt)

	accountingGroup := api.Group("/accounting", requireAuth)
	accountingGroup.Get("/overview", rbac.RequirePermission(rbacRepo, rbac.PageAccounting, rbac.ActionView), accountingHandler.Overview)

	overviewGroup := api.Group("/overview", requireAuth)
	overviewGroup.Get("/revenue", rbac.RequirePermission(rbacRepo, rbac.PageOverview, rbac.ActionView), overviewHandler.Revenue)

	catalogGroup := api.Group("/catalog", requireAuth)
	catalogGroup.Get("/settings", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionView), catalogHandler.GetSettings)
	catalogGroup.Put("/settings", rbac.RequirePermission(rbacRepo, rbac.PageProducts, rbac.ActionUpdate), catalogHandler.UpdateSettings)

	fiberApp.Use(func(c fiber.Ctx) error {
		return response.Fail(c, fiber.StatusNotFound, "route not found")
	})

	return &App{Fiber: fiberApp}
}
