# API Products

Semua endpoint memakai permission `products.*`.

## GET /products

Query:

```text
?page=1&per_page=20&search=kaos&category_id=<category-id>&status=active
```

```bash
curl "$BASE/products?page=1&per_page=20&search=kaos&status=active" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "44444444-4444-4444-8444-444444444444",
        "category_id": "11111111-1111-4111-8111-111111111111",
        "category": {
          "id": "11111111-1111-4111-8111-111111111111",
          "name": "Kaos",
          "slug": "kaos",
          "description": "Kaos dan t-shirt",
          "created_at": "2026-11-04T15:00:00+07:00",
          "updated_at": "2026-11-04T15:00:00+07:00"
        },
        "name": "Kaos Oversize",
        "slug": "kaos-oversize",
        "description": "Kaos oversize cotton combed 30s",
        "status": "active",
        "created_by_id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
        "images": [
          {
            "id": "55555555-5555-4555-8555-555555555555",
            "object_key": "catalog/products/44444444-4444-4444-8444-444444444444/<uuid>.jpg",
            "url": "https://cdn.example.com/catalog/products/44444444-4444-4444-8444-444444444444/<uuid>.jpg",
            "alt_text": "",
            "sort_order": 0,
            "is_primary": true,
            "created_at": "2026-11-04T15:10:00+07:00"
          }
        ],
        "variants": [
          {
            "id": "66666666-6666-4666-8666-666666666666",
            "product_id": "44444444-4444-4444-8444-444444444444",
            "size": {
              "id": "22222222-2222-4222-8222-222222222222",
              "name": "M",
              "sort_order": 1,
              "created_at": "2026-11-04T15:00:00+07:00",
              "updated_at": "2026-11-04T15:00:00+07:00"
            },
            "color": {
              "id": "33333333-3333-4333-8333-333333333333",
              "name": "Black",
              "hex_code": "#111111",
              "created_at": "2026-11-04T15:00:00+07:00",
              "updated_at": "2026-11-04T15:00:00+07:00"
            },
            "sku": "KAOS-OS-M-BLK",
            "price": 150000,
            "stock": 10,
            "is_active": true,
            "created_at": "2026-11-04T15:10:00+07:00",
            "updated_at": "2026-11-04T15:10:00+07:00"
          }
        ],
        "created_at": "2026-11-04T15:10:00+07:00",
        "updated_at": "2026-11-04T15:10:00+07:00"
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

## Field catalog tambahan

Produk dapat menyimpan field khusus halaman catalog:

- `shopee_url` — URL produk di Shopee.
- `shopee_price` — harga Shopee (bigint rupiah).
- `details` — array string detail produk.
- `badge` — label badge, contoh `New`.
- `badge_tone` — `light` atau `accent`.
- `sort_order` — urutan tampil di catalog.

Field ini opsional pada create/update. Produk hanya tampil di catalog publik jika `status = active`.

## POST /products

Membuat produk beserta varian dalam satu transaksi. Gambar diupload terpisah.

Request body:

```json
{
  "category_id": "11111111-1111-4111-8111-111111111111",
  "name": "Kaos Oversize",
  "slug": "kaos-oversize",
  "description": "Kaos oversize cotton combed 30s",
  "status": "active",
  "variants": [
    {
      "size_id": "22222222-2222-4222-8222-222222222222",
      "color_id": "33333333-3333-4333-8333-333333333333",
      "sku": "KAOS-OS-M-BLK",
      "price": 150000,
      "stock": 10,
      "is_active": true
    },
    {
      "size_id": "77777777-7777-4777-8777-777777777777",
      "color_id": "33333333-3333-4333-8333-333333333333",
      "sku": "KAOS-OS-L-BLK",
      "price": 160000,
      "stock": 5,
      "is_active": true
    }
  ]
}
```

```bash
curl -X POST "$BASE/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id":"11111111-1111-4111-8111-111111111111",
    "name":"Kaos Oversize",
    "slug":"kaos-oversize",
    "description":"Kaos oversize cotton combed 30s",
    "status":"active",
    "variants":[
      {"size_id":"22222222-2222-4222-8222-222222222222","color_id":"33333333-3333-4333-8333-333333333333","sku":"KAOS-OS-M-BLK","price":150000,"stock":10,"is_active":true}
    ]
  }'
```

Catatan:

- `category_id`, `slug`, `status`, dan `variants` opsional.
- `status` default `draft`. Nilai valid: `draft`, `active`, `archived`.
- Kombinasi `size_id + color_id` harus unik per produk.
- `sku` unik.
- `price` dan `stock` tidak boleh negatif.

Response `201`: object produk lengkap.

## GET /products/:id

```bash
curl "$BASE/products/<product-id>" \
  -H "Authorization: Bearer $TOKEN"
```

## PUT /products/:id

Request body:

```json
{
  "category_id": "11111111-1111-4111-8111-111111111111",
  "name": "Kaos Oversize Premium",
  "slug": "kaos-oversize-premium",
  "description": "Kaos oversize premium",
  "status": "active"
}
```

```bash
curl -X PUT "$BASE/products/<product-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id":"11111111-1111-4111-8111-111111111111",
    "name":"Kaos Oversize Premium",
    "slug":"kaos-oversize-premium",
    "description":"Kaos oversize premium",
    "status":"active"
  }'
```

## DELETE /products/:id

```bash
curl -X DELETE "$BASE/products/<product-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "product deleted"
  }
}
```

Soft delete. Child rows tidak dihapus.

## Field catalog tambahan

Produk dapat menyimpan field khusus halaman catalog:

- `shopee_url` — URL produk di Shopee.
- `shopee_price` — harga Shopee (bigint rupiah).
- `details` — array string detail produk.
- `badge` — label badge, contoh `New`.
- `badge_tone` — `light` atau `accent`.
- `sort_order` — urutan tampil di catalog.

Field ini opsional pada create/update. Produk hanya tampil di catalog publik jika `status = active`.

## POST /products/:id/images

Upload satu atau beberapa gambar. Field multipart: `images`.

Format JPEG/PNG/WebP, maksimal 10 MB per file.

```bash
curl -X POST "$BASE/products/<product-id>/images" \
  -H "Authorization: Bearer $TOKEN" \
  -F "images=@kaos-1.jpg" \
  -F "images=@kaos-2.jpg"
```

Response `201`:

```json
{
  "success": true,
  "data": [
    {
      "id": "55555555-5555-4555-8555-555555555555",
      "object_key": "catalog/products/<product-id>/<uuid>.jpg",
      "url": "https://cdn.example.com/catalog/products/<product-id>/<uuid>.jpg",
      "alt_text": "",
      "sort_order": 0,
      "is_primary": true,
      "created_at": "2026-11-04T15:10:00+07:00"
    }
  ]
}
```

Gambar pertama otomatis menjadi primary jika produk belum punya primary.

## DELETE /products/:id/images/:imageId

```bash
curl -X DELETE "$BASE/products/<product-id>/images/<image-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "image deleted"
  }
}
```

Object R2 dihapus best-effort. Jika gambar primary dihapus, gambar berikutnya otomatis menjadi primary.

## PUT /products/:id/images/:imageId/primary

```bash
curl -X PUT "$BASE/products/<product-id>/images/<image-id>/primary" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "primary image updated"
  }
}
```

## Field catalog tambahan

Produk dapat menyimpan field khusus halaman catalog:

- `shopee_url` — URL produk di Shopee.
- `shopee_price` — harga Shopee (bigint rupiah).
- `details` — array string detail produk.
- `badge` — label badge, contoh `New`.
- `badge_tone` — `light` atau `accent`.
- `sort_order` — urutan tampil di catalog.

Field ini opsional pada create/update. Produk hanya tampil di catalog publik jika `status = active`.

## POST /products/:id/variants

Request body:

```json
{
  "size_id": "22222222-2222-4222-8222-222222222222",
  "color_id": "33333333-3333-4333-8333-333333333333",
  "sku": "KAOS-OS-XL-BLK",
  "price": 170000,
  "stock": 8,
  "is_active": true
}
```

```bash
curl -X POST "$BASE/products/<product-id>/variants" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "size_id":"22222222-2222-4222-8222-222222222222",
    "color_id":"33333333-3333-4333-8333-333333333333",
    "sku":"KAOS-OS-XL-BLK",
    "price":170000,
    "stock":8,
    "is_active":true
  }'
```

Response `201`: object varian.

## PUT /products/:id/variants/:variantId

Request body:

```json
{
  "size_id": "22222222-2222-4222-8222-222222222222",
  "color_id": "33333333-3333-4333-8333-333333333333",
  "sku": "KAOS-OS-XL-BLK",
  "price": 175000,
  "is_active": true
}
```

```bash
curl -X PUT "$BASE/products/<product-id>/variants/<variant-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "size_id":"22222222-2222-4222-8222-222222222222",
    "color_id":"33333333-3333-4333-8333-333333333333",
    "sku":"KAOS-OS-XL-BLK",
    "price":175000,
    "is_active":true
  }'
```

Stok tidak diubah lewat endpoint ini; gunakan endpoint stock.

## DELETE /products/:id/variants/:variantId

```bash
curl -X DELETE "$BASE/products/<product-id>/variants/<variant-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "variant deleted"
  }
}
```

## PATCH /products/:id/variants/:variantId/stock

Request body:

```json
{
  "stock": 20,
  "note": "Stock opname"
}
```

`note` opsional. Setiap perubahan stok membuat row `stock_movements` bertipe `adjustment`.

```bash
curl -X PATCH "$BASE/products/<product-id>/variants/<variant-id>/stock" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"stock":20,"note":"Stock opname"}'
```

Response `200`: object varian dengan stok terbaru.

## GET /products/:id/variants/:variantId/movements

Riwayat pergerakan stok varian, terbaru lebih dulu.

Query:

```text
?page=1&per_page=20
```

```bash
curl "$BASE/products/<product-id>/variants/<variant-id>/movements?page=1&per_page=20" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "dddddddd-dddd-4ddd-8ddd-dddddddddddd",
        "variant_id": "66666666-6666-4666-8666-666666666666",
        "product_id": "44444444-4444-4444-8444-444444444444",
        "sku": "KAOS-OS-M-BLK",
        "variant_label": "M / Black",
        "type": "order",
        "quantity": -2,
        "stock_before": 10,
        "stock_after": 8,
        "reference_type": "order",
        "reference_id": "99999999-9999-4999-8999-999999999999",
        "note": "order JDC-20261104-ABC123",
        "created_by_id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
        "created_by_name": "Developer",
        "created_at": "2026-11-04T15:30:00+07:00"
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

Tipe movement:

| Type | Arti |
|---|---|
| `initial` | stok awal saat produk/varian dibuat |
| `adjustment` | perubahan stok manual |
| `order` | stok keluar karena order dikonfirmasi |
| `cancel` | stok kembali karena order dibatalkan/dihapus |
| `return` | stok kembali karena pengiriman returned |
