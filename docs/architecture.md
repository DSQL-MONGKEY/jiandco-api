# Architecture

## Ringkasan

API ini adalah backend dashboard internal. Arsitektur dibuat sederhana: handler -> service -> repository -> PostgreSQL, dengan Atlas sebagai pengelola schema.

```text
HTTP request
  -> Fiber middleware (recover, request id, logger, CORS)
  -> auth middleware (JWT + session + user active)
  -> RBAC middleware (page.action)
  -> handler
  -> service
  -> repository (GORM)
  -> PostgreSQL
```

## Struktur folder

```text
cmd/api              entrypoint, config, DB, graceful shutdown
cmd/atlas            mencetak schema GORM untuk Atlas
cmd/seed             seed role, permission, developer user
config               load .env dan validasi env
internal/app         wiring Fiber, route, middleware, error handler
internal/auth        JWT, bcrypt, session, login/refresh/logout/profile/avatar
internal/rbac        katalog permission, middleware, role permission service
internal/user        user CRUD, role assignment, reset password
internal/product     kategori, produk, varian, gambar, stok
internal/customer    customer CRUD
internal/order       order, item, status, pembayaran, pengiriman, stok
internal/accounting  expense, item, struk, accounting overview
internal/overview    revenue overview
internal/models      model GORM + konstanta status
internal/storage     client R2 (S3-compatible)
pkg/database         koneksi GORM + connection pool
pkg/response         envelope sukses/error + pagination
pkg/apperror         error terstruktur dengan HTTP status
pkg/utils            pagination, validasi, rentang bulan
migrations           migration Atlas + atlas.sum
docs                 dokumentasi
```

## Dependency wiring

`internal/app.New` menerima config, `*gorm.DB`, dan storage client. Semua repository, service, dan handler dibuat di sana. Tidak ada global singleton untuk database.

## Error handling

- Service mengembalikan `*apperror.Error` untuk error yang sudah diperkirakan.
- `internal/app.errorHandler` memetakan error tersebut ke envelope error.
- Error tak terduga dikembalikan sebagai `500 internal server error`; detail hanya ditulis ke log.
- 404 route juga memakai envelope yang sama.

## Transaksi

Transaksi database dipakai untuk operasi multi-tabel:

- replace role permissions (delete + insert),
- user delete (soft delete + revoke sessions),
- user password reset (update password + revoke sessions).

Order memakai transaksi untuk membuat header + item, mengurangi/mengembalikan stok, dan mengubah status. Expense memakai transaksi untuk header + item.

## Waktu dan uang

- Timestamp disimpan `timestamptz` (UTC).
- Tanggal bisnis seperti `expense_date` memakai `date`.
- Filter bulan memakai `APP_TIMEZONE`, default `Asia/Jakarta`.
- Uang selalu `bigint` rupiah. Tidak ada `float` untuk nominal uang.

## Dokumentasi API

Referensi endpoint lengkap beserta contoh payload ada di [docs/api/README.md](api/README.md).
