# Deployment

## VPS rjm

Untuk deployment tiga service (API, dashboard, landing) di VPS rjm, gunakan runbook
[deployment-vps-rjm.md](deployment-vps-rjm.md).

## Prasyarat

- Docker untuk build image.
- PostgreSQL 16+ yang dapat diakses dari container API.
- Atlas CLI untuk menjalankan migration.
- Bucket R2 dan kredensial (jika upload file dipakai).
- Reverse proxy/TLS di depan API (disarankan).

## Environment produksi

Minimal:

```env
APP_ENV=production
APP_NAME=jiandco-api
APP_PORT=8080
APP_TIMEZONE=Asia/Jakarta

DATABASE_URL=postgres://user:password@db-host:5432/jiandco_db?sslmode=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10

JWT_SECRET=<random-secret-minimal-32-karakter>
JWT_ISSUER=jiandco-api
JWT_ACCESS_TTL=1h
JWT_REFRESH_TTL=168h

CORS_ALLOWED_ORIGINS=https://www.jiandco.id,https://jiandco.id

R2_ACCOUNT_ID=<account-id>
R2_ACCESS_KEY_ID=<access-key>
R2_SECRET_ACCESS_KEY=<secret-key>
R2_BUCKET=jiandco
R2_PUBLIC_BASE_URL=https://cdn.example.com
```

Catatan:

- `JWT_SECRET` wajib panjang minimal 32 karakter dan tidak boleh memakai placeholder (`change-me`, `replace-with`, `dev-only`, `ci-only`).
- `APP_ENV` valid: `development`, `test`, `staging`, `production`.
- `R2_PUBLIC_BASE_URL` opsional; jika kosong, response berisi object key.
- Jangan menyimpan `.env` produksi di repository.

## Build image

```bash
docker build -t jiandco-api:latest .
```

Image berisi dua binary:

- `/app/api` sebagai entrypoint default.
- `/app/seed` untuk seed role/permission/developer.

Image migration tersedia lewat `Dockerfile.migrations` untuk menjalankan Atlas di VPS
tanpa memasang Atlas CLI di host.

## Migration

Jalankan Atlas dari mesin CI/CD atau job terpisah, bukan dari container API:

```bash
atlas migrate apply --env local
```

Untuk production dari image migration:

```bash
docker run --rm --env-file .env.production \
  --entrypoint /bin/sh \
  ghcr.io/dsql-mongkey/jiandco-api-migrations:<sha> \
  -c 'exec /atlas migrate apply --dir file:///workspace/migrations --url "$DATABASE_URL"'
```

Untuk environment lain, sesuaikan `atlas.hcl` atau gunakan URL database yang sesuai. Migration harus selesai sebelum API versi baru dijalankan.

## Seed

Seed idempotent. Jalankan sekali setelah migration pertama, atau saat menambah permission baru:

```bash
docker run --rm \
  --env-file .env.production \
  jiandco-api:latest /app/seed
```

`SEED_DEVELOPER_PASSWORD` wajib diisi saat developer user belum ada. Seed tidak menimpa password user yang sudah ada.

## Menjalankan API

```bash
docker run -d \
  --name jiandco-api \
  --restart unless-stopped \
  --env-file .env.production \
  -p 8080:8080 \
  jiandco-api:latest
```

## Health check

- `GET /healthz` untuk liveness.
- `GET /readyz` untuk readiness (ping database).

Gunakan endpoint tersebut di orchestrator/load balancer. Jangan arahkan traffic sebelum `/readyz` mengembalikan `200`.

## Reverse proxy

API tidak menangani TLS. Letakkan di belakang Nginx, Caddy, Traefik, atau load balancer yang menangani TLS.

Pastikan:

- Header `X-Forwarded-For`/`X-Forwarded-Proto` diteruskan.
- Body limit proxy minimal 64 MB (batas upload API).
- Timeout proxy minimal 30 detik.
- CORS origin memuat domain landing/catalog yang memanggil API dari browser.

## R2

- Bucket publik: isi `R2_PUBLIC_BASE_URL`.
- Bucket privat: gunakan presigned GET (belum diimplementasikan).
- Object key disimpan di database, bukan URL.
- Backup metadata database; object R2 dikelola terpisah.

## Database

- Gunakan PostgreSQL 16+.
- Aktifkan backup otomatis dan uji restore.
- Gunakan koneksi TLS (`sslmode=require` atau lebih ketat) untuk database remote.
- Pantau koneksi dan ukuran tabel `auth_sessions`/`stock_movements`.

## CI

Workflow `.github/workflows/ci.yml` menjalankan:

- format check,
- `go vet`,
- `go test`,
- `go build`,
- validasi migration,
- apply migration ke PostgreSQL service,
- seed,
- build & push image API + migration ke GHCR,
- deploy ke VPS rjm lewat SSH saat push ke `main`/`master`.

## Rollback

- Rollback aplikasi dengan image sebelumnya.
- Migration bersifat forward-only; buat migration baru untuk memperbaiki schema.
- Jangan mengedit migration yang sudah diterapkan di environment bersama.
