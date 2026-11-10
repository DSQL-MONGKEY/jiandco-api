# AGENTS.md

## Project

Jiandco internal dashboard API. Go + Fiber + GORM + PostgreSQL, with Atlas as the only schema migration tool.

Current phase: schema, documentation, auth, RBAC, user management, catalog/products, customers/orders, accounting, revenue overview, inventory audit/stock movements, and production packaging/CI are implemented.

## Stack

- Go 1.26
- Fiber v3
- GORM v1.31 (runtime only, never AutoMigrate)
- PostgreSQL 16+
- Atlas v1.x with GORM external schema provider
- Cloudflare R2 via AWS SDK v2 S3 client
- JWT HS256 (`github.com/golang-jwt/jwt/v5`)
- bcrypt cost 12

## Structure

```text
cmd/api        API entrypoint
cmd/atlas      GORM schema provider for Atlas
cmd/seed       idempotent roles/permissions/developer seed
config         env loading and validation
internal/app   Fiber wiring, middleware, error handler
internal/auth  JWT, sessions, login/refresh/logout/profile/avatar
internal/rbac  permission catalog, middleware, role permission management
internal/user  user CRUD and role assignment
internal/product catalog, products, variants, images, stock
internal/customer customer CRUD
internal/order orders, items, status transitions, stock
internal/accounting expenses, items, receipts, accounting overview
internal/overview revenue overview
internal/models all GORM models and constants
internal/storage R2 client
pkg/database   GORM connection
pkg/response   success/error envelope
pkg/apperror   typed application errors
pkg/utils      pagination, validation, month range
migrations     Atlas migrations and atlas.sum (committed)
docs           project documentation
```

## Commands

```bash
go mod tidy
gofmt -w cmd config internal pkg
go vet ./...
go test ./...

go run ./cmd/api
go run ./cmd/seed
docker build -t jiandco-api .

atlas migrate diff <name> --env local
atlas migrate validate --env local
atlas migrate apply --env local
atlas migrate status --env local
```

If the compiler runs out of memory on Windows, use `GOFLAGS=-p=1` for build/test commands.

## Database rules

- Atlas is the source of truth. Never call `AutoMigrate` in application code.
- Every schema change must include:
  1. model update in `internal/models`
  2. generated Atlas migration committed in `migrations/`
  3. `docs/database-schema.md` update
- Use UUID primary keys with `gen_random_uuid()`.
- Use `timestamptz` for timestamps and `date` for business dates.
- Store money as `bigint` rupiah. Never use floating point for money.
- Use soft delete only for master/history tables. Child rows are hard deleted.
- Unique indexes remain reserved for soft-deleted rows. Do not reuse usernames, emails, slugs, SKUs, or order numbers.
- Add indexes for every foreign key used in filters or joins.
- Use transactions for multi-table writes (role permissions, user delete/password reset, order/expense writes).
- Every stock change must write a `stock_movements` row with snapshot SKU/label, before/after stock, and reference.

## API rules

- Base path: `/api/v1`.
- Success envelope: `{"success": true, "data": ...}`.
- Error envelope: `{"success": false, "code": <http status>, "message": "..."}`.
- Lists use `data: {items: [], pagination: {...}}`.
- `code` equals the HTTP status code.
- Use `pkg/response` helpers; do not hand-roll envelopes.
- Use `pkg/apperror` for expected errors.
- Validate request bodies with `pkg/utils.Validate`.
- Pagination defaults: `page=1`, `per_page=20`, max `100`.
- Dates are RFC3339 UTC in responses. Business month filters use `APP_TIMEZONE` (default `Asia/Jakarta`).

## Auth and RBAC rules

- Access token: JWT HS256, 1 hour. Refresh token: 7 days, rotated on every refresh.
- Store only SHA-256 hashes of refresh tokens.
- Every authenticated request validates JWT, session state, and user active state.
- Level 0 Developer and level 1 Super Admin bypass permission checks.
- Level 2/3 permissions come from `role_permissions`.
- Only level 0/1 may change role permissions.
- Actors may only manage users with a lower privilege (higher numeric level), except Developer.
- Never return `password_hash` or raw refresh tokens.

## R2 rules

- Store object keys in the database, not full URLs.
- Key conventions:
  - `user/avatars/{userID}/{uuid}.{ext}`
  - `catalog/products/{productID}/{uuid}.{ext}`
  - `catalog/products/{productID}/colors/{colorID}/{uuid}.{ext}`
  - `accounting/receipts/{expenseID}/{uuid}.{ext}`
- Validate MIME type and file size before upload.
- Max upload: avatar 10 MB, product image 10 MB per file, receipt 10 MB.
- Delete replaced avatar and removed product image objects on a best-effort basis.

## Documentation rules

Every change must update the relevant file under `docs/`:

- schema change -> `docs/database-schema.md`
- API contract change -> `docs/api-conventions.md` and `docs/api/<module>.md`
- auth/session change -> `docs/auth.md`
- permission change -> `docs/rbac.md`
- storage key/limit change -> `docs/storage-r2.md`
- migration workflow change -> `docs/migrations.md`
- deployment/packaging change -> `docs/deployment.md`
- notable phase change -> `docs/changelog.md`

API request/response examples live in `docs/api/`. Every endpoint must have a payload or query example there.

## Environment

See `.env.example`. Required:

- `DATABASE_URL`
- `JWT_SECRET` (minimum 32 characters)

R2 variables are optional in development. Upload endpoints return 503 when R2 is not configured.

## Definition of done

- Code compiles and is gofmt'd.
- `go vet ./...` passes.
- `go test ./...` passes.
- Atlas migration generated and reviewed if schema changed.
- Docs updated.
- No secrets committed.
