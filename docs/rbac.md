# RBAC

## Level

| Level | Role | Perilaku |
|---|---|---|
| 0 | Developer | bypass semua permission; boleh mengelola user level 0/1 |
| 1 | Super Admin | bypass semua permission |
| 2 | Admin | permission default operasional + user view/create/update + role view |
| 3 | Staff | permission default harian |

Level 0 dan 1 tidak memakai `role_permissions`. Middleware langsung meloloskan request.

## Permission

Format key: `page.action`.

| Page | Action | Key |
|---|---|---|
| overview | view | `overview.view` |
| products | view/create/update/delete | `products.view`, `products.create`, `products.update`, `products.delete` |
| orders | view/create/update/delete | `orders.view`, `orders.create`, `orders.update`, `orders.delete` |
| customers | view/create/update/delete | `customers.view`, `customers.create`, `customers.update`, `customers.delete` |
| accounting | view/create/update/delete | `accounting.view`, `accounting.create`, `accounting.update`, `accounting.delete` |
| users | view/create/update/delete | `users.view`, `users.create`, `users.update`, `users.delete` |
| roles | view/update | `roles.view`, `roles.update` |

Halaman `products` mencakup kategori, size, color, varian, gambar, dan stok.

## Middleware

`rbac.RequirePermission(repo, page, action)`:

1. Ambil user dari context auth.
2. Jika `role.level <= 1`, lanjut.
3. Cek `role_permissions` join `permissions` untuk `(page, action)`.
4. Jika tidak ada, response `403`.

Permission dicek dari database setiap request. Perubahan permission langsung berlaku tanpa menunggu token baru. `GET /auth/me` mengembalikan daftar permission efektif user untuk kebutuhan navigasi frontend; level 0/1 menerima seluruh permission key.

## Endpoint

```text
GET /api/v1/roles
GET /api/v1/permissions
GET /api/v1/roles/:id/permissions
PUT /api/v1/roles/:id/permissions
```

`PUT` memerlukan `roles.update`. Hanya level 0/1 yang boleh mengubah permission role, dan hanya role level 2/3 yang dapat diubah. Role Developer/Super Admin selalu full access.

Request:

```json
{
  "permission_ids": ["uuid", "uuid"]
}
```

Mengirim array kosong akan menghapus semua permission role tersebut.

## Aturan user management

- Aktor hanya boleh mengelola user dengan level lebih rendah (angka lebih besar).
- Developer (level 0) dapat mengelola semua level.
- Super Admin (level 1) dapat mengelola level 2/3.
- Admin (level 2) dapat mengelola level 3.
- User tidak dapat mengubah role sendiri, menonaktifkan diri sendiri, atau menghapus diri sendiri.
- Role assignment mengikuti aturan yang sama: target role harus lebih rendah dari aktor, kecuali Developer.
- Reset password dan nonaktif user mencabut semua session user tersebut.

## Seed default

`cmd/seed` mengisi role dan permission. Default permission role 2/3 hanya diisi jika role belum memiliki permission, sehingga kustomisasi tidak tertimpa.

| Permission | Admin (2) | Staff (3) |
|---|---|---|
| overview.view | ya | ya |
| products.view | ya | ya |
| products.create/update/delete | ya | tidak |
| orders.view/create/update/delete | ya | view/create/update |
| customers.view/create/update/delete | ya | view/create/update |
| accounting.view/create/update/delete | ya | view/create |
| users.view/create/update | ya | tidak |
| users.delete | tidak | tidak |
| roles.view | ya | tidak |
| roles.update | tidak | tidak |
