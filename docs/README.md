# Jiandco API Documentation

Dokumentasi internal dashboard brand clothing Jiandco.

## Status implementasi

| Area | Schema | Dokumentasi | Handler |
|---|---|---|---|
| Auth & session | Selesai | Selesai | Selesai |
| RBAC | Selesai | Selesai | Selesai |
| User management | Selesai | Selesai | Selesai |
| Products | Selesai | Selesai | Selesai |
| Accounting | Selesai | Selesai | Selesai |
| Orders | Selesai | Selesai | Selesai |
| Revenue overview | Selesai | Selesai | Selesai |
| R2 storage | Selesai | Selesai | Selesai |
| Inventory audit | Selesai | Selesai | Selesai |

## Daftar dokumen

### Referensi API

- [API Reference](api/README.md)
- [Health](api/health.md)
- [Auth API](api/auth.md)
- [RBAC API](api/rbac.md)
- [Users API](api/users.md)
- [Catalog API](api/catalog.md)
- [Products API](api/products.md)
- [Public Catalog API](api/catalog-public.md)
- [Catalog Settings API](api/catalog-settings.md)
- [Customers API](api/customers.md)
- [Orders API](api/orders.md)
- [Accounting API](api/accounting.md)
- [Overview API](api/overview.md)

### Dokumen konsep dan operasional

- [Architecture](architecture.md)
- [Database Schema](database-schema.md)
- [API Conventions](api-conventions.md)
- [Auth](auth.md)
- [RBAC](rbac.md)
- [Users](users.md)
- [Products](products.md)
- [Accounting](accounting.md)
- [Orders](orders.md)
- [Overview](overview.md)
- [Storage R2](storage-r2.md)
- [Migrations](migrations.md)
- [Development](development.md)
- [Deployment](deployment.md)
- [Deployment VPS rjm](deployment-vps-rjm.md)
- [Changelog](changelog.md)

## Prinsip utama

- Semua endpoint berada di `/api/v1`.
- Semua response memakai envelope `success`/`data` atau `success`/`code`/`message`.
- Atlas adalah satu-satunya sumber perubahan schema. `AutoMigrate` dilarang.
- Uang disimpan sebagai `bigint` rupiah, waktu sebagai `timestamptz`.
- Permission memakai format `page.action`, contoh `products.create`.
- File disimpan di R2 sebagai object key, bukan URL penuh.
