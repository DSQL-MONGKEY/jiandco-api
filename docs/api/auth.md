# API Auth

## POST /auth/login

Login memakai username dan password.

Request body:

```json
{
  "username": "developer",
  "password": "Developer123!"
}
```

```bash
curl -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"developer","password":"Developer123!"}'
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "token": "<jwt>",
    "refreshToken": "<refresh-token>",
    "expiresIn": 3600,
    "user": {
      "id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
      "role_id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
      "role": {
        "id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
        "name": "Developer",
        "level": 0
      },
      "name": "Developer",
      "username": "developer",
      "email": "developer@jiandco.local",
      "phone": "",
      "avatar_key": "",
      "avatar_url": "",
      "is_active": true,
      "last_login_at": "2026-11-04T16:00:00+07:00",
      "created_at": "2026-11-04T15:00:00+07:00",
      "updated_at": "2026-11-04T16:00:00+07:00"
    }
  }
}
```

Error:

- `401` username/password salah.
- `403` user inactive.
- `429` terlalu banyak request.

## POST /auth/refresh

Rotasi refresh token. Refresh token lama langsung tidak berlaku.

Request body:

```json
{
  "refreshToken": "<refresh-token>"
}
```

```bash
curl -X POST "$BASE/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refreshToken":"<refresh-token>"}'
```

Response `200` sama dengan login (`token`, `refreshToken`, `expiresIn`, `user`).

## POST /auth/logout

Memerlukan access token. Session aktif dicabut.

```bash
curl -X POST "$BASE/auth/logout" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "logged out"
  }
}
```

## GET /auth/me

```bash
curl "$BASE/auth/me" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
      "role_id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
      "role": {
        "id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
        "name": "Developer",
        "level": 0
      },
      "name": "Developer",
      "username": "developer",
      "email": "developer@jiandco.local",
      "phone": "",
      "avatar_key": "",
      "avatar_url": "",
      "is_active": true,
      "last_login_at": "2026-11-04T16:00:00+07:00",
      "created_at": "2026-11-04T15:00:00+07:00",
      "updated_at": "2026-11-04T16:00:00+07:00"
    },
    "permissions": [
      "products.view",
      "orders.create"
    ]
  }
}
```

Level 0/1 menerima seluruh permission key. Level 2/3 menerima permission dari `role_permissions`.

## PUT /auth/profile

Request body:

```json
{
  "name": "Developer Updated",
  "email": "developer@jiandco.local",
  "phone": "0811111111"
}
```

```bash
curl -X PUT "$BASE/auth/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Developer Updated","email":"developer@jiandco.local","phone":"0811111111"}'
```

Response `200`: object user terbaru.

## PUT /auth/password

Request body:

```json
{
  "currentPassword": "Developer123!",
  "newPassword": "NewPassword123!"
}
```

```bash
curl -X PUT "$BASE/auth/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"currentPassword":"Developer123!","newPassword":"NewPassword123!"}'
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "password updated"
  }
}
```

Session lain milik user dicabut, session saat ini tetap aktif.

## POST /auth/avatar

Upload avatar. Format JPEG/PNG/WebP, maksimal 10 MB.

```bash
curl -X POST "$BASE/auth/avatar" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@avatar.jpg"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "avatar_key": "user/avatars/<user-id>/<uuid>.jpg",
    "avatar_url": "https://cdn.example.com/user/avatars/<user-id>/<uuid>.jpg"
  }
}
```

Jika R2 belum dikonfigurasi, response `503 file storage is not configured`.
