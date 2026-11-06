# API Catalog

Mencakup kategori, size, dan color. Semua endpoint memakai permission `products.*`.

## Categories

### GET /categories

Query opsional:

```text
?search=kaos
```

```bash
curl "$BASE/categories?search=kaos" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "11111111-1111-4111-8111-111111111111",
      "name": "Kaos",
      "slug": "kaos",
      "description": "Kaos dan t-shirt",
      "created_at": "2026-11-04T15:00:00+07:00",
      "updated_at": "2026-11-04T15:00:00+07:00"
    }
  ]
}
```

### POST /categories

Request body:

```json
{
  "name": "Kaos",
  "slug": "kaos",
  "description": "Kaos dan t-shirt"
}
```

`slug` opsional; jika kosong dibuat dari `name`.

```bash
curl -X POST "$BASE/categories" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Kaos","slug":"kaos","description":"Kaos dan t-shirt"}'
```

Response `201`: object kategori.

### GET /categories/:id

```bash
curl "$BASE/categories/<category-id>" \
  -H "Authorization: Bearer $TOKEN"
```

### PUT /categories/:id

Request body:

```json
{
  "name": "Kaos Premium",
  "slug": "kaos-premium",
  "description": "Kaos premium"
}
```

```bash
curl -X PUT "$BASE/categories/<category-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Kaos Premium","slug":"kaos-premium","description":"Kaos premium"}'
```

### DELETE /categories/:id

```bash
curl -X DELETE "$BASE/categories/<category-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "category deleted"
  }
}
```

Soft delete. Produk yang memakai kategori ini akan dilepas kategorinya.

## Sizes

### GET /sizes

```bash
curl "$BASE/sizes" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "22222222-2222-4222-8222-222222222222",
      "name": "M",
      "sort_order": 1,
      "created_at": "2026-11-04T15:00:00+07:00",
      "updated_at": "2026-11-04T15:00:00+07:00"
    }
  ]
}
```

### POST /sizes

Request body:

```json
{
  "name": "M",
  "sort_order": 1
}
```

```bash
curl -X POST "$BASE/sizes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"M","sort_order":1}'
```

### GET /sizes/:id

```bash
curl "$BASE/sizes/<size-id>" \
  -H "Authorization: Bearer $TOKEN"
```

### PUT /sizes/:id

Request body:

```json
{
  "name": "Medium",
  "sort_order": 2
}
```

```bash
curl -X PUT "$BASE/sizes/<size-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Medium","sort_order":2}'
```

### DELETE /sizes/:id

```bash
curl -X DELETE "$BASE/sizes/<size-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Size yang dipakai varian produk tidak dapat dihapus (`409`).

## Colors

### GET /colors

```bash
curl "$BASE/colors" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "33333333-3333-4333-8333-333333333333",
      "name": "Black",
      "hex_code": "#111111",
      "created_at": "2026-11-04T15:00:00+07:00",
      "updated_at": "2026-11-04T15:00:00+07:00"
    }
  ]
}
```

### POST /colors

Request body:

```json
{
  "name": "Black",
  "hex_code": "#111111"
}
```

```bash
curl -X POST "$BASE/colors" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Black","hex_code":"#111111"}'
```

`hex_code` harus format `#RRGGBB` jika diisi.

### GET /colors/:id

```bash
curl "$BASE/colors/<color-id>" \
  -H "Authorization: Bearer $TOKEN"
```

### PUT /colors/:id

Request body:

```json
{
  "name": "Jet Black",
  "hex_code": "#000000"
}
```

```bash
curl -X PUT "$BASE/colors/<color-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Jet Black","hex_code":"#000000"}'
```

### DELETE /colors/:id

```bash
curl -X DELETE "$BASE/colors/<color-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Color yang dipakai varian produk tidak dapat dihapus (`409`).
