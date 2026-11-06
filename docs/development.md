# Development

## Prasyarat

- Go 1.26+
- PostgreSQL 16+
- Atlas CLI v1.x
- Docker (opsional, untuk dev database Atlas)
- R2 credentials (opsional)

## Setup

```bash
cp .env.example .env
```

Isi minimal:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/jiandco_db?sslmode=disable
JWT_SECRET=random-secret-at-least-32-characters
SEED_DEVELOPER_PASSWORD=Developer123!
```

Jalankan PostgreSQL:

```bash
docker compose up -d postgres
```

Atau gunakan PostgreSQL lokal yang sudah ada.

## Migration dan seed

```bash
atlas migrate apply --env local
go run ./cmd/seed
```

Jika Docker tidak berjalan, lihat [Migrations](migrations.md) untuk memakai dev database lokal.

## Menjalankan API

```bash
go run ./cmd/api
```

Default port `8080`. Override dengan `APP_PORT`.

Health check:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## Login awal

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"developer","password":"Developer123!"}'
```

## Perintah kualitas

```bash
go mod tidy
gofmt -w cmd config internal pkg
go vet ./...
go test ./...
```

Jika compiler kehabisan memori di Windows:

```bash
GOFLAGS=-p=1 go build ./...
GOFLAGS=-p=1 go test ./...
```

## Referensi API

Contoh request lengkap untuk semua endpoint ada di [docs/api/README.md](api/README.md).

## Smoke test

1. Login developer.
2. `GET /api/v1/auth/me` dengan Bearer token.
3. Refresh token, pastikan refresh lama ditolak.
4. Buat user staff melalui developer.
5. Login staff, pastikan endpoint `roles`/`users` mengembalikan `403`.
6. Logout, pastikan token yang sama mengembalikan `401`.
7. Buat kategori, size, color, produk, dan varian.
8. Buat customer dan order, lalu uji transisi status serta pengurangan stok.
9. Buat expense dan cek accounting overview.
10. Cek revenue overview dengan filter bulan.

## Troubleshooting

- `bind: An attempt was made to access a socket...` -> port sudah dipakai; set `APP_PORT` lain.
- `atlas: docker daemon not running` -> jalankan Docker Desktop atau gunakan `--dev-url` lokal.
- `JWT_SECRET must be at least 32 characters` -> isi secret yang cukup panjang.
- Upload `503` -> R2 belum dikonfigurasi; isi variabel R2.
- `record not found` pada seed developer -> normal saat pengecekan awal; user akan dibuat.
