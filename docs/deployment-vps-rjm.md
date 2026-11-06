# Deployment VPS rjm — Jiandco API, Dashboard, dan Landing

Runbook ini menjelaskan cara menjalankan tiga service Jiandco di VPS rjm:

| Service | Repo | Domain | Image |
|---|---|---|---|
| API | `jiandco-api` | `api.jiandco.id` | `ghcr.io/dsql-mongkey/jiandco-api` |
| Dashboard | `jiandco-dashboard` | `dash.jiandco.id` | `ghcr.io/dsql-mongkey/jiandco-dashboard` |
| Landing / catalog | `oinsh-page` | `www.jiandco.id` | `ghcr.io/dsql-mongkey/oinsh-page` |

Arsitektur yang dipakai adalah **stack terpisah, shared infra**:

- Compose project baru di `/opt/jiandco`.
- Memakai network Docker, PostgreSQL, dan Cloudflare Tunnel milik stack rjm.
- Service Jiandco diberi nama unik (`jiandco-*`) agar tidak bentrok dengan service rjm.
- Database baru `jiandco` dibuat di container PostgreSQL rjm.

Referensi pola VPS, Cloudflare Tunnel, GHCR, dan backup ada di repo terpisah
`rjm-go-api/docs/deployment.md`.

## 1. Prasyarat

- VPS rjm sudah menjalankan Docker, stack rjm, PostgreSQL, dan Cloudflare Tunnel.
- Akses SSH sebagai user `rjm`.
- Docker sudah login ke GHCR untuk pull image private.
- DNS `jiandco.id` dikelola di Cloudflare.
- R2 bucket dan kredensial tersedia jika upload gambar/struk dipakai.

Cek network dan container rjm:

```bash
docker network ls --format '{{.Name}}' | grep rjm
docker ps --format '{{.Names}}\t{{.Image}}' | grep -E 'rjm|cloudflared|postgres'
```

Nama network biasanya `rjm_rjm`. Nilai aktual dipakai sebagai `RJM_NETWORK`.

## 2. Cloudflare Tunnel, DNS, dan Access

Tambahkan Public Hostname pada tunnel rjm yang sudah ada:

| Subdomain | Domain       | Service                         |
| --------- | ------------ | ------------------------------- |
| `www`     | `jiandco.id` | `http://jiandco-landing:80`     |
| `dash`    | `jiandco.id` | `http://jiandco-dashboard:3000` |
| `api`     | `jiandco.id` | `http://jiandco-api:8080`       |

Langkah:

1. Cloudflare → Zero Trust → Networks → Tunnels → tunnel rjm → **Public Hostname**.
2. Tambahkan tiga hostname di atas satu per satu.
3. Hapus A record lama untuk hostname tersebut jika ada; Cloudflare akan membuat CNAME
   `<tunnel-id>.cfargotunnel.com` secara otomatis.
4. Tambahkan redirect rule `jiandco.id` → `https://www.jiandco.id` (301) jika apex belum
   dipakai.

Lindungi dashboard:

1. Zero Trust → Access → Applications → Add self-hosted.
2. Domain: `dash.jiandco.id`.
3. Policy: Allow → Emails → email admin.
4. Login method: One-time PIN.
5. Session duration: 24 jam.

Tambahkan WAF rate limit untuk login API:

- Name: `jiandco-api-auth-login`
- Expression: `http.host eq "api.jiandco.id" and http.request.uri.path contains "/auth/login"`
- Rate: 5 requests / 1 minute / IP
- Action: Block, timeout 60 detik

API dan landing tetap publik karena landing memanggil API dari browser.

## 3. Database di PostgreSQL rjm

Masuk ke stack rjm:

```bash
cd /opt/rjm
docker compose exec db psql -U postgres -c "SHOW max_connections;"
```

Buat role dan database Jiandco:

```bash
docker compose exec db psql -U postgres -c \
  "CREATE ROLE jiandco LOGIN PASSWORD '<password-jiandco>';"
docker compose exec db psql -U postgres -c \
  "CREATE DATABASE jiandco OWNER jiandco;"
```

Catatan kapasitas koneksi:

- PostgreSQL rjm memakai `max_connections=20` secara default.
- API Jiandco memakai pool yang dapat diatur lewat `DB_MAX_OPEN_CONNS` dan
  `DB_MAX_IDLE_CONNS`.
- Untuk shared PostgreSQL, gunakan `DB_MAX_OPEN_CONNS=5` dan `DB_MAX_IDLE_CONNS=2`.
- Jika API rjm + Jiandco + migration bersamaan masih kehabisan koneksi, naikkan
  `max_connections` PostgreSQL rjm menjadi 40 dan restart service `db`.

## 4. Direktori dan environment VPS

```bash
sudo mkdir -p /opt/jiandco
sudo chown -R rjm:rjm /opt/jiandco
cd /opt/jiandco
```

Salin `deploy/vps/docker-compose.yml` dari repo `jiandco-api` ke
`/opt/jiandco/docker-compose.yml`.

Buat file berikut:

### `/opt/jiandco/.env` (chmod 600)

```env
RJM_NETWORK=rjm_rjm
IMAGE_PREFIX=ghcr.io/dsql-mongkey
API_TAG=latest
DASHBOARD_TAG=latest
LANDING_TAG=latest
```

### `/opt/jiandco/.env.api` (chmod 600)

```env
APP_ENV=production
APP_NAME=jiandco-api
APP_PORT=8080
APP_TIMEZONE=Asia/Jakarta

DATABASE_URL=postgres://jiandco:<password-jiandco>@db:5432/jiandco?sslmode=disable
DB_MAX_OPEN_CONNS=5
DB_MAX_IDLE_CONNS=2

JWT_SECRET=<random-minimal-32-karakter>
JWT_ISSUER=jiandco-api
JWT_ACCESS_TTL=1h
JWT_REFRESH_TTL=168h

CORS_ALLOWED_ORIGINS=https://www.jiandco.id,https://jiandco.id

R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET=jiandco
R2_ENDPOINT=
R2_PUBLIC_BASE_URL=

SEED_DEVELOPER_NAME=Developer
SEED_DEVELOPER_USERNAME=developer
SEED_DEVELOPER_EMAIL=developer@jiandco.local
SEED_DEVELOPER_PASSWORD=<password-developer>
```

### `/opt/jiandco/.env.dashboard` (chmod 600)

```env
BACKEND_URL=http://jiandco-api:8080/api/v1
NEXT_PUBLIC_APP_URL=https://dash.jiandco.id
BUILD_STANDALONE=true
```

Landing tidak memakai env runtime. `VITE_API_BASE_URL` di-inline saat build image.

```bash
chmod 600 /opt/jiandco/.env /opt/jiandco/.env.api /opt/jiandco/.env.dashboard
```

## 5. Compose stack Jiandco

Compose menggunakan network external milik rjm dan tidak membuka port host:

- `jiandco-api` → `jiandco-api:8080`
- `jiandco-dashboard` → `jiandco-dashboard:3000`
- `jiandco-landing` → `jiandco-landing:80`
- `jiandco-migrate` → Atlas, profile `tools`
- `jiandco-seed` → seed binary, profile `tools`

Cek resolusi network:

```bash
cd /opt/jiandco
docker compose config
docker network inspect "${RJM_NETWORK:-rjm_rjm}" --format '{{range .Containers}}{{.Name}} {{end}}'
```

## 6. Login GHCR di VPS

Buat GitHub Personal Access Token (classic) dengan scope `read:packages`, lalu:

```bash
echo "<PAT>" | docker login ghcr.io -u DSQL-MONGKEY --password-stdin
```

Lakukan sekali. Jika PAT expired, pull akan gagal dengan `denied`.

## 7. Deploy pertama

Urutan wajib: database → migration → seed → API → dashboard → landing.

```bash
cd /opt/jiandco

docker compose --profile tools pull jiandco-api jiandco-migrate
docker compose --profile tools run --rm jiandco-migrate
docker compose --profile tools run --rm jiandco-seed

docker compose pull jiandco-dashboard jiandco-landing
docker compose up -d jiandco-api jiandco-dashboard jiandco-landing
```

Verifikasi:

```bash
docker compose ps

# API dari dalam network
docker run --rm --network "${RJM_NETWORK:-rjm_rjm}" \
  curlimages/curl:8.10.1 -fsS http://jiandco-api:8080/readyz

# API publik
curl -fsS https://api.jiandco.id/readyz
curl -fsS https://api.jiandco.id/api/v1/public/catalog | head -c 300

# Landing dan dashboard
curl -I https://www.jiandco.id
curl -I https://dash.jiandco.id
```

Buka `https://dash.jiandco.id`, login Access Cloudflare, lalu login dengan akun developer
yang di-seed.

## 8. CI/CD GitHub Actions

Setiap repo memiliki workflow `.github/workflows/ci.yml`:

| Repo | Verify | Image | Deploy |
|---|---|---|---|
| `jiandco-api` | gofmt, vet, test, build, Atlas, seed | API + migrations | migrate lalu API |
| `jiandco-dashboard` | typecheck, lint, build | dashboard | update `DASHBOARD_TAG` |
| `oinsh-page` | lint, build | landing | update `LANDING_TAG` |

Workflow berjalan pada push `main`/`master`, PR, dan `workflow_dispatch`.

### GitHub Secrets

Set di ketiga repo (atau Organization secrets):

| Secret | Isi |
|---|---|
| `VPS_HOST` | IP publik VPS rjm |
| `VPS_USER` | `rjm` |
| `VPS_SSH_KEY` | private key deploy |

### GitHub Variables

| Variable | Nilai |
|---|---|
| `VPS_APP_DIR` | `/opt/jiandco` |
| `IMAGE_PREFIX` | `ghcr.io/dsql-mongkey` (opsional, jika owner GHCR berbeda) |

### Alur deploy

```text
push main/master
  -> verify
  -> build & push image ke GHCR (tag sha + latest)
  -> SSH ke VPS
  -> update tag di /opt/jiandco/.env
  -> API: pull image, jalankan migration, up -d --wait
  -> dashboard/landing: pull image, up -d --wait
```

Semua job deploy memakai `flock /tmp/jiandco-deploy.lock` agar tiga repo tidak saling
menimpa `.env` saat deploy bersamaan.

### Build-time values

- Landing: `VITE_API_BASE_URL=https://api.jiandco.id/api/v1` di-inline saat build.
- Dashboard: `NEXT_PUBLIC_APP_URL=https://dash.jiandco.id` di-inline saat build.
- Jika domain berubah, rebuild image, bukan hanya mengganti env runtime.

## 9. Rollback

Cari SHA image sebelumnya di GHCR → Versions atau `docker images` di VPS.

```bash
cd /opt/jiandco

# API
sed -i 's/^API_TAG=.*/API_TAG=<sha-lama>/' .env
docker compose --profile tools pull jiandco-api jiandco-migrate
docker compose up -d --wait jiandco-api

# Dashboard
sed -i 's/^DASHBOARD_TAG=.*/DASHBOARD_TAG=<sha-lama>/' .env
docker compose pull jiandco-dashboard
docker compose up -d --wait jiandco-dashboard

# Landing
sed -i 's/^LANDING_TAG=.*/LANDING_TAG=<sha-lama>/' .env
docker compose pull jiandco-landing
docker compose up -d --wait jiandco-landing
```

Migration bersifat forward-only. Jangan mengedit migration yang sudah diterapkan; buat
migration baru untuk memperbaiki schema.

## 10. Backup database

Tambahkan dump database Jiandco ke backup harian rjm:

```bash
cd /opt/rjm
STAMP=$(date -u +%Y%m%d-%H%M%S)
docker compose exec -T db pg_dump -U jiandco -d jiandco -Fc \
  | gzip > "backups/jiandco-$STAMP.dump.gz"
rclone copy "backups/jiandco-$STAMP.dump.gz" r2:rjm-backups/db/
find backups -name 'jiandco-*.dump.gz' -mtime +7 -delete
```

Uji restore ke database sementara sebelum benar-benar dibutuhkan.

## 11. Troubleshooting

| Gejala | Penyebab | Solusi |
|---|---|---|
| Cloudflare 502 | Service tidak ada di network / nama salah | Cek `docker network inspect` dan ingress `jiandco-api:8080` |
| Dashboard 500 ke API | `BACKEND_URL` salah | Harus `http://jiandco-api:8080/api/v1` |
| CORS error di landing | Origin belum diizinkan | Tambahkan `https://www.jiandco.id` ke `CORS_ALLOWED_ORIGINS` |
| Login dashboard loop | Cookie tidak tersimpan / Access | Akses via HTTPS, cek `NEXT_PUBLIC_APP_URL` dan Access policy |
| Migration checksum mismatch | Migration diubah setelah apply | Jangan edit migration lama; buat migration baru |
| `too many connections` | Pool Jiandco + rjm melebihi `max_connections` | Set `DB_MAX_OPEN_CONNS=5`, atau naikkan `max_connections` |
| `denied` saat pull GHCR | PAT VPS expired | `docker login ghcr.io` ulang |
| Image lama masih jalan | Tag tidak berubah | Cek `.env`, `docker compose pull`, `docker compose up -d` |
| Landing memanggil API lama | `VITE_API_BASE_URL` build-time | Rebuild image landing dengan URL yang benar |
| RAM VPS penuh | Dashboard + monitoring + API | Cek `docker stats`, turunkan `mem_limit`, pertimbangkan pisahkan DB |

## 12. File contoh

- Compose: `deploy/vps/docker-compose.yml`
- Env root: `deploy/vps/.env.example`
- Env API: `deploy/vps/.env.api.example`
- Env dashboard: `deploy/vps/.env.dashboard.example`
