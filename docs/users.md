# Users

Status: selesai untuk fase 1.

## Endpoint

```text
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
PUT    /api/v1/users/:id/password
```

Permission:

- list/detail: `users.view`
- create: `users.create`
- update/reset password: `users.update`
- delete: `users.delete`

## List

```text
?page=1&per_page=20&search=&role_id=&is_active=
```

Response:

```json
{
  "success": true,
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 0,
      "total_pages": 0
    }
  }
}
```

`search` mencari nama, username, dan email.

## Create

```json
{
  "name": "Staff Satu",
  "username": "staff1",
  "email": "staff1@jiandco.local",
  "password": "Staff123!",
  "phone": "08123456789",
  "role_id": "uuid-role-staff",
  "is_active": true
}
```

- Username dan email dinormalisasi ke lowercase.
- Password di-hash bcrypt cost 12.
- `is_active` default true jika tidak dikirim.
- Username/email duplikat -> `409`.

## Update

```json
{
  "name": "Staff Satu Updated",
  "email": "staff1@jiandco.local",
  "phone": "08123456789",
  "role_id": "uuid-role-staff",
  "is_active": true
}
```

Username tidak dapat diubah. Menonaktifkan user mencabut semua session-nya.

## Delete

Soft delete user dan revoke semua session. User tidak dapat menghapus dirinya sendiri.

## Reset password

```json
{
  "new_password": "NewPassword123!"
}
```

Password baru minimal 8 karakter, maksimal 72 byte. Semua session user dicabut setelah reset.

## Aturan privilege

- Aktor hanya dapat mengelola user dengan level lebih rendah.
- Developer dapat mengelola semua level.
- Super Admin dapat mengelola level 2/3.
- Admin dapat mengelola level 3.
- User tidak dapat mengubah role sendiri atau menonaktifkan diri sendiri.
