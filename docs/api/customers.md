# API Customers

Semua endpoint memakai permission `customers.*`.

## GET /customers

Query:

```text
?page=1&per_page=20&search=budi
```

```bash
curl "$BASE/customers?page=1&per_page=20&search=budi" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "88888888-8888-4888-8888-888888888888",
        "name": "Budi",
        "phone": "081298765432",
        "email": "budi@example.com",
        "address": "Jl. Test No. 1",
        "city": "Jakarta",
        "province": "DKI Jakarta",
        "postal_code": "12345",
        "notes": "",
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

## POST /customers

Request body:

```json
{
  "name": "Budi",
  "phone": "081298765432",
  "email": "budi@example.com",
  "address": "Jl. Test No. 1",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "notes": "Pelanggan langganan"
}
```

```bash
curl -X POST "$BASE/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Budi",
    "phone":"081298765432",
    "email":"budi@example.com",
    "address":"Jl. Test No. 1",
    "city":"Jakarta",
    "province":"DKI Jakarta",
    "postal_code":"12345",
    "notes":"Pelanggan langganan"
  }'
```

Response `201`: object customer.

## GET /customers/:id

```bash
curl "$BASE/customers/<customer-id>" \
  -H "Authorization: Bearer $TOKEN"
```

## PUT /customers/:id

Request body:

```json
{
  "name": "Budi Updated",
  "phone": "081298765432",
  "email": "budi@example.com",
  "address": "Jl. Test No. 2",
  "city": "Bandung",
  "province": "Jawa Barat",
  "postal_code": "40123",
  "notes": ""
}
```

```bash
curl -X PUT "$BASE/customers/<customer-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Budi Updated",
    "phone":"081298765432",
    "email":"budi@example.com",
    "address":"Jl. Test No. 2",
    "city":"Bandung",
    "province":"Jawa Barat",
    "postal_code":"40123",
    "notes":""
  }'
```

## DELETE /customers/:id

```bash
curl -X DELETE "$BASE/customers/<customer-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "customer deleted"
  }
}
```

Customer yang sudah memiliki order tidak dapat dihapus (`409`).
