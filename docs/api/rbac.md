# API RBAC

## GET /roles

```bash
curl "$BASE/roles" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "9a1b2c3d-4e5f-6789-abcd-222222222222",
      "name": "Developer",
      "level": 0,
      "description": "Full access to every feature and system setting",
      "is_system": true
    }
  ]
}
```

## GET /permissions

```bash
curl "$BASE/permissions" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "1b2c3d4e-5f67-8901-abcd-333333333333",
      "key": "products.create",
      "page": "products",
      "action": "create",
      "label": "Create products"
    }
  ]
}
```

## GET /roles/:id/permissions

```bash
curl "$BASE/roles/<role-id>/permissions" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "role": {
      "id": "<role-id>",
      "name": "Staff",
      "level": 3,
      "description": "Daily operational access",
      "is_system": true
    },
    "permissions": [
      {
        "id": "<permission-id>",
        "key": "orders.view",
        "page": "orders",
        "action": "view",
        "label": "View orders"
      }
    ]
  }
}
```

## PUT /roles/:id/permissions

Mengganti seluruh permission role. Kirim array kosong untuk mencabut semua.

Request body:

```json
{
  "permission_ids": [
    "1b2c3d4e-5f67-8901-abcd-333333333333",
    "2c3d4e5f-6789-0123-abcd-444444444444"
  ]
}
```

```bash
curl -X PUT "$BASE/roles/<role-id>/permissions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"permission_ids":["1b2c3d4e-5f67-8901-abcd-333333333333"]}'
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "role": {
      "id": "<role-id>",
      "name": "Staff",
      "level": 3,
      "description": "Daily operational access",
      "is_system": true
    },
    "permissions": [
      {
        "id": "1b2c3d4e-5f67-8901-abcd-333333333333",
        "key": "products.create",
        "page": "products",
        "action": "create",
        "label": "Create products"
      }
    ]
  }
}
```

Aturan:

- Hanya level 0/1 yang boleh mengubah permission role.
- Role level 0/1 selalu full access dan tidak dapat diubah permission-nya.
- Permission tidak dikenal -> `400`.
