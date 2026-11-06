# Auth

## Ringkasan

Autentikasi dikelola sendiri di PostgreSQL. Login memakai `username` + `password`, menghasilkan access token JWT dan refresh token.

- Access token: JWT HS256, berlaku 1 jam.
- Refresh token: 7 hari, dirotasi setiap refresh.
- Satu login membuat satu row `auth_sessions`.
- Refresh token asli tidak disimpan; database hanya menyimpan SHA-256 hash.
- Setiap request terautentikasi memeriksa signature JWT, session aktif, dan user aktif.

## Endpoint

### POST /api/v1/auth/login

Request:

```json
{
  "username": "developer",
  "password": "Developer123!"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "token": "<jwt>",
    "refreshToken": "<opaque-token>",
    "expiresIn": 3600,
    "user": {
      "id": "uuid",
      "role_id": "uuid",
      "role": { "id": "uuid", "name": "Developer", "level": 0 },
      "name": "Developer",
      "username": "developer",
      "email": "developer@jiandco.local",
      "phone": "",
      "avatar_key": "",
      "avatar_url": "",
      "is_active": true,
      "last_login_at": "2026-11-04T15:00:00+07:00",
      "created_at": "2026-11-04T14:58:48+07:00",
      "updated_at": "2026-11-04T15:00:00+07:00"
    }
  }
}
```

Username/password yang salah memakai pesan generik `invalid username or password`. User yang ditemukan tetapi `is_active=false` mengembalikan `403 user is inactive`.

### POST /api/v1/auth/refresh

Request:

```json
{ "refreshToken": "<refresh-token>" }
```

Response sama dengan login. Refresh token lama langsung tidak berlaku setelah rotasi.

### POST /api/v1/auth/logout

Memerlukan access token. Mengubah `revoked_at` pada session aktif. Access token yang sama langsung ditolak karena middleware memeriksa session.

### GET /api/v1/auth/me

Mengembalikan user yang sedang login beserta permission efektif:

```json
{
  "success": true,
  "data": {
    "user": { "...": "..." },
    "permissions": ["products.view", "orders.create"]
  }
}
```

Level 0/1 menerima semua permission key. Level 2/3 menerima permission dari `role_permissions`.

### PUT /api/v1/auth/profile

Request:

```json
{
  "name": "Developer Updated",
  "email": "developer@jiandco.local",
  "phone": "0811111111"
}
```

Username tidak dapat diubah dari endpoint ini.

### PUT /api/v1/auth/password

Request:

```json
{
  "currentPassword": "Developer123!",
  "newPassword": "NewPassword123!"
}
```

Password baru minimal 8 karakter, maksimal 72 byte (batas bcrypt). Session lain milik user dicabut; session saat ini tetap aktif.

### POST /api/v1/auth/avatar

`multipart/form-data`, field `file`. Validasi:

- JPEG, PNG, atau WebP.
- Maksimal 10 MB.
- Object key: `user/avatars/{userID}/{uuid}.{ext}`.
- Avatar lama dihapus best-effort setelah update berhasil.
- Jika R2 belum dikonfigurasi, response `503`.

## JWT claims

```json
{
  "sid": "session-uuid",
  "role_level": 0,
  "username": "developer",
  "iss": "jiandco-api",
  "sub": "user-uuid",
  "exp": 1793782787,
  "nbf": 1793779187,
  "iat": 1793779187
}
```

`sid` dipakai middleware untuk memuat session. `role_level` hanya informasi; RBAC tetap memuat role terbaru dari database agar perubahan role langsung berlaku.

## Session lifecycle

1. Login berhasil -> buat `auth_sessions` dengan hash refresh token dan `expires_at = now + 7 hari`.
2. Refresh -> cek hash, expiry, revoked, user aktif -> rotasi hash -> access token baru.
3. Logout -> set `revoked_at`.
4. Ganti password -> cabut session lain.
5. User dinonaktifkan/dihapus -> session dicabut, middleware menolak.

Refresh token memakai expiry absolut 7 hari dari login. Refresh tidak memperpanjang umur session.

## Rate limit

Login dan refresh dibatasi 10 request per menit per IP. Response `429` memakai envelope error standar.

## Keamanan

- Password di-hash bcrypt cost 12.
- JWT secret minimal 32 karakter dan wajib diisi.
- Jangan mengembalikan `password_hash`, refresh token mentah, atau hash refresh token.
- Session diperiksa per request; ini menambah satu query database tetapi membuat logout/nonaktif user langsung berlaku.
