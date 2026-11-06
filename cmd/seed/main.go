package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jianco-api/config"
	"jianco-api/internal/auth"
	"jianco-api/internal/catalog"
	"jianco-api/internal/models"
	"jianco-api/internal/rbac"
	"jianco-api/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL, cfg.IsProduction(), cfg.DBMaxOpenConn, cfg.DBMaxIdleConn)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	if err := seedRoles(db); err != nil {
		slog.Error("failed to seed roles", "error", err)
		os.Exit(1)
	}
	permissions, err := seedPermissions(db)
	if err != nil {
		slog.Error("failed to seed permissions", "error", err)
		os.Exit(1)
	}
	if err := seedRolePermissions(db, permissions); err != nil {
		slog.Error("failed to seed role permissions", "error", err)
		os.Exit(1)
	}
	if err := seedDeveloper(db, cfg); err != nil {
		slog.Error("failed to seed developer user", "error", err)
		os.Exit(1)
	}
	if err := seedCatalog(db); err != nil {
		slog.Error("failed to seed catalog defaults", "error", err)
		os.Exit(1)
	}

	slog.Info("seed completed")
}

func seedRoles(db *gorm.DB) error {
	for _, role := range rbac.RoleCatalog {
		record := models.Role{
			ID:          uuid.New(),
			Name:        role.Name,
			Level:       role.Level,
			Description: role.Description,
			IsSystem:    true,
		}
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "level"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_system", "updated_at"}),
		}).Create(&record).Error
		if err != nil {
			return fmt.Errorf("upsert role level %d: %w", role.Level, err)
		}
	}
	return nil
}

func seedPermissions(db *gorm.DB) (map[string]uuid.UUID, error) {
	for _, permission := range rbac.PermissionCatalog {
		record := models.Permission{
			ID:     uuid.New(),
			Key:    permission.Key,
			Page:   permission.Page,
			Action: permission.Action,
			Label:  permission.Label,
		}
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"page", "action", "label", "updated_at"}),
		}).Create(&record).Error
		if err != nil {
			return nil, fmt.Errorf("upsert permission %s: %w", permission.Key, err)
		}
	}

	var records []models.Permission
	if err := db.Find(&records).Error; err != nil {
		return nil, err
	}

	result := make(map[string]uuid.UUID, len(records))
	for _, record := range records {
		result[record.Key] = record.ID
	}
	return result, nil
}

func seedRolePermissions(db *gorm.DB, permissions map[string]uuid.UUID) error {
	for level, keys := range rbac.DefaultRolePermissions {
		var role models.Role
		if err := db.First(&role, "level = ?", level).Error; err != nil {
			return fmt.Errorf("find role level %d: %w", level, err)
		}

		var existing int64
		if err := db.Model(&models.RolePermission{}).Where("role_id = ?", role.ID).Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			slog.Info("role permissions already customized, skipping defaults", "role", role.Name)
			continue
		}

		rows := make([]models.RolePermission, 0, len(keys))
		for _, key := range keys {
			permissionID, ok := permissions[key]
			if !ok {
				return fmt.Errorf("permission %s not found after seeding", key)
			}
			rows = append(rows, models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permissionID,
			})
		}
		if len(rows) == 0 {
			continue
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error; err != nil {
			return fmt.Errorf("seed permissions for role %s: %w", role.Name, err)
		}
	}
	return nil
}

func seedDeveloper(db *gorm.DB, cfg *config.Config) error {
	var role models.Role
	if err := db.First(&role, "level = ?", models.RoleLevelDeveloper).Error; err != nil {
		return fmt.Errorf("find developer role: %w", err)
	}

	username := strings.ToLower(strings.TrimSpace(cfg.Seed.DeveloperUsername))
	var existing models.User
	err := db.Unscoped().Where("username = ?", username).First(&existing).Error
	if err == nil {
		slog.Info("developer user already exists, skipping", "username", username)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if strings.TrimSpace(cfg.Seed.DeveloperPassword) == "" {
		return fmt.Errorf("SEED_DEVELOPER_PASSWORD is required when creating the developer user")
	}

	hash, err := auth.HashPassword(cfg.Seed.DeveloperPassword)
	if err != nil {
		return err
	}

	user := models.User{
		RoleID:       role.ID,
		Name:         cfg.Seed.DeveloperName,
		Username:     username,
		Email:        strings.ToLower(strings.TrimSpace(cfg.Seed.DeveloperEmail)),
		PasswordHash: hash,
		IsActive:     true,
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	slog.Info("developer user created", "username", username)
	return nil
}

func seedCatalog(db *gorm.DB) error {
	categories := []models.Category{
		{Name: "Chongsam", Slug: "chongsam", Description: "Lace vest dengan detail romantis, dibuat untuk dipadukan dengan apa pun di lemari kamu."},
		{Name: "Blush", Slug: "blush", Description: "Satin halter dengan drape lembut dan warna hangat untuk hari yang ingin terasa istimewa."},
	}
	for index := range categories {
		if err := db.Where("slug = ?", categories[index].Slug).FirstOrCreate(&categories[index]).Error; err != nil {
			return err
		}
	}

	sizes := []models.Size{
		{Name: "S", SortOrder: 1},
		{Name: "M", SortOrder: 2},
		{Name: "L", SortOrder: 3},
		{Name: "XL", SortOrder: 4},
	}
	for index := range sizes {
		if err := db.Where("name = ?", sizes[index].Name).FirstOrCreate(&sizes[index]).Error; err != nil {
			return err
		}
	}

	colors := []models.Color{
		{Name: "Cocoa", HexCode: "#6B4A3A"},
		{Name: "Ivory", HexCode: "#F3EBDD"},
		{Name: "Rose", HexCode: "#C8A49F"},
		{Name: "Peach", HexCode: "#E8B9A6"},
		{Name: "Cream", HexCode: "#EEE4DA"},
		{Name: "Onyx", HexCode: "#1A1A1A"},
		{Name: "Sand", HexCode: "#D8C4AC"},
		{Name: "Navy", HexCode: "#1F2A44"},
	}
	for index := range colors {
		if err := db.Where("name = ?", colors[index].Name).FirstOrCreate(&colors[index]).Error; err != nil {
			return err
		}
	}

	var settingsCount int64
	if err := db.Model(&models.CatalogSetting{}).Where("key = ?", "default").Count(&settingsCount).Error; err != nil {
		return err
	}
	if settingsCount == 0 {
		if err := db.Create(catalog.DefaultSettings()).Error; err != nil {
			return err
		}
	}

	slog.Info("catalog defaults seeded")
	return nil
}
