# Changelog

## 2026-11-05 - Fase 7

- Menaikkan batas upload gambar produk dan avatar menjadi 10 MB.
- Menaikkan body limit API menjadi 64 MB untuk upload beberapa gambar.
- Memperbarui dokumentasi batas upload dan body limit reverse proxy.

## 2026-11-05 - Fase 6

- Menambahkan workflow CI/CD API: verify, build/push image GHCR, migration, dan deploy SSH.
- Menambahkan `Dockerfile.migrations` untuk Atlas di VPS.
- Menambahkan konfigurasi pool database `DB_MAX_OPEN_CONNS` dan `DB_MAX_IDLE_CONNS`.
- Menambahkan runbook deployment tiga service di VPS rjm (shared network, PostgreSQL, tunnel).
- Menambahkan contoh compose dan env VPS di `deploy/vps/`.

## 2026-11-05 - Fase 5

- Menambahkan field catalog pada produk: Shopee URL/harga, details, badge, badge tone, sort order.
- Menambahkan `catalog_settings` dan endpoint admin/public untuk halaman catalog.
- Menambahkan endpoint publik catalog, product detail, dan media proxy R2.
- Menambahkan `order_status_breakdown` pada overview.
- Menambahkan seed kategori Chongsam/Blush, size S–XL, color catalog, dan settings default.
- Dashboard: admin catalog settings, money visibility, bar/pie chart, invoice A4, struk 80mm, tema Ji&Co, logo.
- oinsh-page: catalog/product/cart/checkout menggunakan API publik dan settings backend.

## 2026-11-04 - Fase 4

- Menambahkan Dockerfile multi-stage (API + seed, distroless, timezone embedded).
- Menambahkan `.dockerignore` dan workflow CI GitHub Actions.
- Menambahkan helmet, compression, dan timeout HTTP.
- Menambahkan validasi `APP_ENV` dan penolakan JWT secret placeholder di production.
- Menambahkan dokumentasi deployment.

## 2026-11-04 - Fase 3

- Menambahkan tabel `stock_movements` dan migration `20261104131154_add_stock_movements.sql`.
- Mencatat stok awal, adjustment manual, order, cancel, dan return ke ledger.
- Menambahkan endpoint riwayat pergerakan stok per varian.
- Menambahkan edit item order selama status masih `pending`.
- Menambahkan endpoint hapus struk expense.

## 2026-11-04 - Fase 2

- Menambahkan CRUD kategori, size, color, produk, varian, gambar produk, dan update stok.
- Menambahkan upload gambar produk ke R2 (`catalog/products/...`) dan hapus object best-effort.
- Menambahkan CRUD customer.
- Menambahkan CRUD order, item snapshot, status order, pembayaran, pengiriman, resi, dan transaksi stok.
- Menambahkan `orders.stock_deducted` dan migration `20261104084102_add_order_stock_deducted.sql`.
- Menambahkan CRUD expense + item, upload struk R2 (`accounting/receipts/...`), dan accounting overview.
- Menambahkan revenue overview dengan filter bulan WIB.
- Menambahkan referensi API lengkap dengan contoh request untuk setiap endpoint di `docs/api/`.

## 2026-11-04 - Fase 1

- Menambahkan schema lengkap untuk auth/session, RBAC, katalog produk, customer/order, dan accounting.
- Menambahkan migration Atlas awal `20261104075754_init.sql`.
- Menambahkan seed idempotent untuk role, permission, default role permission, dan developer user.
- Menambahkan response envelope sukses/error dan pagination.
- Menambahkan auth: login, refresh rotation, logout, me, profile, password, avatar.
- Menambahkan RBAC middleware dan endpoint role/permission.
- Menambahkan user management: list, create, detail, update, delete, reset password.
- Menambahkan R2 client dan aturan object key.
- Menambahkan dokumentasi `/docs` dan `AGENTS.md`.
- Modul products, customers/orders, accounting, dan overview masih berupa schema + kontrak; handler menyusul.
