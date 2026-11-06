# API Reference

Base URL:

```text
http://localhost:8080/api/v1
```

Semua endpoint kecuali health, login, dan refresh memerlukan header:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

Untuk upload file gunakan `multipart/form-data`.

## Envelope

Sukses:

```json
{
  "success": true,
  "data": {}
}
```

List:

```json
{
  "success": true,
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 0,
      "total_pages": 0
    }
  }
}
```

Error:

```json
{
  "success": false,
  "code": 400,
  "message": "validation failed"
}
```

`code` sama dengan HTTP status. Detail lengkap ada di [API Conventions](../api-conventions.md).

## Endpoint index

### Health

| Method | Path | Auth |
|---|---|---|
| GET | `/healthz` | tidak |
| GET | `/readyz` | tidak |

Dokumentasi: [health.md](health.md)

### Auth

| Method | Path | Permission |
|---|---|---|
| POST | `/auth/login` | publik |
| POST | `/auth/refresh` | publik |
| POST | `/auth/logout` | login |
| GET | `/auth/me` | login |
| PUT | `/auth/profile` | login |
| PUT | `/auth/password` | login |
| POST | `/auth/avatar` | login |

Dokumentasi: [auth.md](auth.md)

### RBAC & Users

| Method | Path | Permission |
|---|---|---|
| GET | `/roles` | `roles.view` |
| GET | `/permissions` | `roles.view` |
| GET | `/roles/:id/permissions` | `roles.view` |
| PUT | `/roles/:id/permissions` | `roles.update` |
| GET | `/users` | `users.view` |
| POST | `/users` | `users.create` |
| GET | `/users/:id` | `users.view` |
| PUT | `/users/:id` | `users.update` |
| DELETE | `/users/:id` | `users.delete` |
| PUT | `/users/:id/password` | `users.update` |

Dokumentasi: [rbac.md](rbac.md), [users.md](users.md)

### Catalog

| Method | Path | Permission |
|---|---|---|
| GET/POST | `/categories` | `products.view` / `products.create` |
| GET/PUT/DELETE | `/categories/:id` | `products.view` / `products.update` / `products.delete` |
| GET/POST | `/sizes` | `products.view` / `products.create` |
| GET/PUT/DELETE | `/sizes/:id` | `products.view` / `products.update` / `products.delete` |
| GET/POST | `/colors` | `products.view` / `products.create` |
| GET/PUT/DELETE | `/colors/:id` | `products.view` / `products.update` / `products.delete` |

Dokumentasi: [catalog.md](catalog.md)

### Public Catalog

| Method | Path | Auth |
|---|---|---|
| GET | `/public/catalog` | publik |
| GET | `/public/catalog/products/:slug` | publik |
| GET | `/public/media/*` | publik |
| GET | `/catalog/settings` | `products.view` |
| PUT | `/catalog/settings` | `products.update` |

Dokumentasi: [catalog-public.md](catalog-public.md), [catalog-settings.md](catalog-settings.md)

### Products

| Method | Path | Permission |
|---|---|---|
| GET/POST | `/products` | `products.view` / `products.create` |
| GET/PUT/DELETE | `/products/:id` | `products.view` / `products.update` / `products.delete` |
| POST | `/products/:id/images` | `products.update` |
| DELETE | `/products/:id/images/:imageId` | `products.update` |
| PUT | `/products/:id/images/:imageId/primary` | `products.update` |
| POST | `/products/:id/variants` | `products.create` |
| PUT | `/products/:id/variants/:variantId` | `products.update` |
| DELETE | `/products/:id/variants/:variantId` | `products.delete` |
| PATCH | `/products/:id/variants/:variantId/stock` | `products.update` |
| GET | `/products/:id/variants/:variantId/movements` | `products.view` |

Dokumentasi: [products.md](products.md)

### Customers

| Method | Path | Permission |
|---|---|---|
| GET/POST | `/customers` | `customers.view` / `customers.create` |
| GET/PUT/DELETE | `/customers/:id` | `customers.view` / `customers.update` / `customers.delete` |

Dokumentasi: [customers.md](customers.md)

### Orders

| Method | Path | Permission |
|---|---|---|
| GET/POST | `/orders` | `orders.view` / `orders.create` |
| GET/PUT/DELETE | `/orders/:id` | `orders.view` / `orders.update` / `orders.delete` |
| PUT | `/orders/:id/items` | `orders.update` |
| PATCH | `/orders/:id/status` | `orders.update` |
| PATCH | `/orders/:id/payment` | `orders.update` |
| PATCH | `/orders/:id/shipping` | `orders.update` |

Dokumentasi: [orders.md](orders.md)

### Accounting

| Method | Path | Permission |
|---|---|---|
| GET/POST | `/expenses` | `accounting.view` / `accounting.create` |
| GET/PUT/DELETE | `/expenses/:id` | `accounting.view` / `accounting.update` / `accounting.delete` |
| POST | `/expenses/:id/receipt` | `accounting.update` |
| DELETE | `/expenses/:id/receipt` | `accounting.update` |
| GET | `/accounting/overview` | `accounting.view` |

Dokumentasi: [accounting.md](accounting.md)

### Overview

| Method | Path | Permission |
|---|---|---|
| GET | `/overview/revenue` | `overview.view` |

Dokumentasi: [overview.md](overview.md)

## Pagination

Query umum:

```text
?page=1&per_page=20
```

- `page` default 1.
- `per_page` default 20, maksimum 100.
- `search` tersedia di list users, products, customers, orders, expenses.

## Contoh variabel shell

```bash
BASE=http://localhost:8080/api/v1
TOKEN=<access_token>
PRODUCT_ID=<uuid>
ORDER_ID=<uuid>
```

Semua contoh di dokumentasi memakai placeholder UUID agar mudah disalin.
