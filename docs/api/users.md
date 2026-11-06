# API Users

## GET /users

Query:

```text
?page=1&per_page=20&search=budi&role_id=<role-id>&is_active=true
```

```bash
curl "$BASE/users?page=1&per_page=20&search=budi" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
        "role_id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
        "role": {
          "id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
          "name": "Staff",
          "level": 3
        },
        "name": "Staff Satu",
        "username": "staff1",
        "email": "staff1@jiandco.local",
        "phone": "08123456789",
        "avatar_key": "",
        "avatar_url": "",
        "is_active": true,
        "last_login_at": null,
        "created_at": "2026-11-04T15:00:00+07:00",
        "updated_at": "2026-11-04T15:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

## POST /users

Request body:

```json
{
  "name": "Staff Satu",
  "username": "staff1",
  "email": "staff1@jiandco.local",
  "password": "Staff123!",
  "phone": "08123456789",
  "role_id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
  "is_active": true
}
```

```bash
curl -X POST "$BASE/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Staff Satu",
    "username":"staff1",
    "email":"staff1@jiandco.local",
    "password":"Staff123!",
    "phone":"08123456789",
    "role_id":"9a1b2c3d-4e5f-6789-abcd-222222222222",
    "is_active":true
  }'
```

Response `201`: object user.

Catatan:

- `is_active` default `true` jika tidak dikirim.
- Username/email duplikat -> `409`.
- Role yang diberikan harus lebih rendah dari level aktor (kecuali Developer).

## GET /users/:id

```bash
curl "$BASE/users/<user-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`: object user.

## PUT /users/:id

Request body:

```json
{
  "name": "Staff Satu Updated",
  "email": "staff1@jiandco.local",
  "phone": "08123456789",
  "role_id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
  "is_active": true
}
```

```bash
curl -X PUT "$BASE/users/<user-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Staff Satu Updated",
    "email":"staff1@jiandco.local",
    "phone":"08123456789",
    "role_id":"9a1b2c3d-4e5f-6789-abcd-222222222222",
    "is_active":true
  }'
```

Username tidak dapat diubah. Menonaktifkan user mencabut semua session-nya.

## DELETE /users/:id

```bash
curl -X DELETE "$BASE/users/<user-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "user deleted"
  }
}
```

User tidak dapat menghapus dirinya sendiri.

## PUT /users/:id/password

Request body:

```json
{
  "new_password": "NewPassword123!"
}
```

```bash
curl -X PUT "$BASE/users/<user-id>/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"new_password":"NewPassword123!"}'
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "password reset"
  }
}
```

Semua session user tersebut dicabut setelah reset.
