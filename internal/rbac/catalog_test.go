package rbac

import (
	"testing"

	"jianco-api/internal/models"
)

func TestPermissionCatalogIsUnique(t *testing.T) {
	keys := make(map[string]struct{}, len(PermissionCatalog))
	pages := make(map[string]struct{}, len(PermissionCatalog))

	for _, permission := range PermissionCatalog {
		if _, exists := keys[permission.Key]; exists {
			t.Fatalf("duplicate permission key: %s", permission.Key)
		}
		keys[permission.Key] = struct{}{}

		pair := permission.Page + "." + permission.Action
		if _, exists := pages[pair]; exists {
			t.Fatalf("duplicate page/action: %s", pair)
		}
		pages[pair] = struct{}{}

		if permission.Key != pair {
			t.Fatalf("permission key %s does not match page/action %s", permission.Key, pair)
		}
	}
}

func TestRoleCatalogCoversAllLevels(t *testing.T) {
	levels := make(map[int]struct{}, len(RoleCatalog))
	for _, role := range RoleCatalog {
		levels[role.Level] = struct{}{}
	}

	for level := models.RoleLevelDeveloper; level <= models.RoleLevelStaff; level++ {
		if _, exists := levels[level]; !exists {
			t.Fatalf("missing role for level %d", level)
		}
	}
}

func TestDefaultRolePermissionsExistInCatalog(t *testing.T) {
	keys := make(map[string]struct{}, len(PermissionCatalog))
	for _, permission := range PermissionCatalog {
		keys[permission.Key] = struct{}{}
	}

	for level, permissions := range DefaultRolePermissions {
		for _, key := range permissions {
			if _, exists := keys[key]; !exists {
				t.Fatalf("default permission %s for level %d is not in catalog", key, level)
			}
		}
	}
}
