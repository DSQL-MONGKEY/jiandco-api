# Migrations

## Prinsip

- Atlas adalah satu-satunya alat perubahan schema.
- Aplikasi tidak pernah memanggil `AutoMigrate`.
- Model GORM di `internal/models` menjadi sumber schema deklaratif.
- Migration dan `atlas.sum` wajib di-commit.
- Setiap perubahan schema wajib memperbarui `docs/database-schema.md`.

## Alur perubahan schema

1. Ubah/tambah model di `internal/models`.
2. Pastikan model terdaftar di `cmd/atlas/main.go`.
3. Jalankan:

```bash
atlas migrate diff <nama_perubahan> --env local
```

4. Review SQL yang dihasilkan: FK, index, check constraint, default, nullable.
5. Validasi checksum dan SQL:

```bash
atlas migrate validate --env local
```

Catatan: `atlas migrate lint` memerlukan Atlas Pro sejak v0.38. Jika tersedia, jalankan juga:

```bash
atlas migrate lint --env local --latest 1
```

6. Apply ke database lokal:

```bash
atlas migrate apply --env local
```

7. Update dokumentasi schema.

## Konfigurasi Atlas

`atlas.hcl` memakai:

- `external_schema` yang menjalankan `go run -mod=mod ./cmd/atlas`.
- Target lokal: `postgres://postgres:postgres@localhost:5432/jiandco_db`.
- Dev database default: Docker PostgreSQL 16.

Jika Docker tidak tersedia, gunakan database dev lokal:

```bash
atlas migrate diff <nama> --env local \
  --dev-url "postgres://postgres:postgres@localhost:5432/jiandco_atlas_dev?sslmode=disable&search_path=public"
```

Buat database dev tersebut terlebih dahulu jika belum ada, misalnya dengan `createdb jiandco_atlas_dev`.

## Seed

Setelah migration:

```bash
go run ./cmd/seed
```

Seed bersifat idempotent:

- Upsert role berdasarkan `level`.
- Upsert permission berdasarkan `key`.
- Isi default permission role 2/3 hanya jika role belum punya permission.
- Buat developer user dari env `SEED_DEVELOPER_*` jika belum ada.

`SEED_DEVELOPER_PASSWORD` wajib diisi saat developer user belum ada.

## Production

CI membangun image migration dari `Dockerfile.migrations`. Image tersebut berisi folder
`migrations` dan Atlas CLI:

```bash
docker run --rm --env-file .env.api \
  --entrypoint /bin/sh \
  ghcr.io/dsql-mongkey/jiandco-api-migrations:<sha> \
  -c 'exec /atlas migrate apply --dir file:///workspace/migrations --url "$DATABASE_URL"'
```

Di VPS rjm, migration dijalankan melalui service `jiandco-migrate` sebelum API baru
di-swap. Lihat [deployment-vps-rjm.md](deployment-vps-rjm.md).

## Rollback

Atlas mendukung down migration jika file migration menyertakan statement down. Untuk fase awal, perubahan schema dilakukan maju (forward-only) dan review manual. Jangan mengedit migration yang sudah diterapkan di environment bersama; buat migration baru.

## Naming

Atlas memberi nama file dengan timestamp. Gunakan nama deskriptif, contoh:

```text
atlas migrate diff add_stock_movements --env local
```
