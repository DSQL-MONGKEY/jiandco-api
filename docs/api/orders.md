# API Orders

Semua endpoint memakai permission `orders.*`.

## GET /orders

Query:

```text
?page=1&per_page=20&search=JDC&status=confirmed&payment_status=paid&shipping_status=pending&customer_id=<customer-id>&from=2026-11-01&to=2026-11-30
```

```bash
curl "$BASE/orders?page=1&per_page=20&status=confirmed&payment_status=paid" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "99999999-9999-4999-8999-999999999999",
        "order_number": "JDC-20261104-ABC123",
        "customer_id": "88888888-8888-4888-8888-888888888888",
        "customer": {
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
        },
        "status": "confirmed",
        "payment_status": "paid",
        "payment_method": "transfer",
        "paid_at": "2026-11-04T15:30:00+07:00",
        "shipping_status": "pending",
        "courier": "JNE",
        "tracking_number": "",
        "shipped_at": null,
        "delivered_at": null,
        "recipient_name": "Budi",
        "recipient_phone": "081298765432",
        "shipping_address": "Jl. Test No. 1",
        "shipping_city": "Jakarta",
        "shipping_province": "DKI Jakarta",
        "shipping_postal_code": "12345",
        "subtotal": 300000,
        "discount_amount": 10000,
        "shipping_cost": 20000,
        "total_amount": 310000,
        "notes": "",
        "stock_deducted": true,
        "created_by_id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
        "items": [
          {
            "id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "product_id": "44444444-4444-4444-8444-444444444444",
            "variant_id": "66666666-6666-4666-8666-666666666666",
            "product_name": "Kaos Oversize",
            "variant_label": "M / Black",
            "sku": "KAOS-OS-M-BLK",
            "price": 150000,
            "quantity": 2,
            "subtotal": 300000
          }
        ],
        "created_at": "2026-11-04T15:20:00+07:00",
        "updated_at": "2026-11-04T15:30:00+07:00"
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

## POST /orders

Request body:

```json
{
  "customer_id": "88888888-8888-4888-8888-888888888888",
  "items": [
    {
      "variant_id": "66666666-6666-4666-8666-666666666666",
      "quantity": 2
    }
  ],
  "recipient_name": "Budi",
  "recipient_phone": "081298765432",
  "shipping_address": "Jl. Test No. 1",
  "shipping_city": "Jakarta",
  "shipping_province": "DKI Jakarta",
  "shipping_postal_code": "12345",
  "payment_method": "transfer",
  "courier": "JNE",
  "discount_amount": 10000,
  "shipping_cost": 20000,
  "notes": "Kirim via JNE"
}
```

```bash
curl -X POST "$BASE/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"88888888-8888-4888-8888-888888888888",
    "items":[{"variant_id":"66666666-6666-4666-8666-666666666666","quantity":2}],
    "recipient_name":"Budi",
    "recipient_phone":"081298765432",
    "shipping_address":"Jl. Test No. 1",
    "shipping_city":"Jakarta",
    "shipping_province":"DKI Jakarta",
    "shipping_postal_code":"12345",
    "payment_method":"transfer",
    "courier":"JNE",
    "discount_amount":10000,
    "shipping_cost":20000,
    "notes":"Kirim via JNE"
  }'
```

Catatan:

- Harga item diambil dari harga varian saat order dibuat (snapshot).
- `recipient_name`, `recipient_phone`, `shipping_address` boleh kosong; service mengambil dari data customer. Jika tetap kosong -> `400`.
- Stok harus mencukupi saat order dibuat. Stok baru dikurangi saat status menjadi `confirmed`.
- `total_amount = subtotal - discount_amount + shipping_cost`.
- `discount_amount` tidak boleh melebihi subtotal.

Response `201`: object order.

## GET /orders/:id

```bash
curl "$BASE/orders/<order-id>" \
  -H "Authorization: Bearer $TOKEN"
```

## PUT /orders/:id

Mengubah data header order. Item dan total dihitung ulang dari item yang sudah ada.

Request body:

```json
{
  "customer_id": "88888888-8888-4888-8888-888888888888",
  "recipient_name": "Budi Updated",
  "recipient_phone": "081298765432",
  "shipping_address": "Jl. Test No. 2",
  "shipping_city": "Bandung",
  "shipping_province": "Jawa Barat",
  "shipping_postal_code": "40123",
  "payment_method": "transfer",
  "courier": "JNE",
  "discount_amount": 0,
  "shipping_cost": 25000,
  "notes": "Alamat diubah"
}
```

```bash
curl -X PUT "$BASE/orders/<order-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"88888888-8888-4888-8888-888888888888",
    "recipient_name":"Budi Updated",
    "recipient_phone":"081298765432",
    "shipping_address":"Jl. Test No. 2",
    "shipping_city":"Bandung",
    "shipping_province":"Jawa Barat",
    "shipping_postal_code":"40123",
    "payment_method":"transfer",
    "courier":"JNE",
    "discount_amount":0,
    "shipping_cost":25000,
    "notes":"Alamat diubah"
  }'
```

## DELETE /orders/:id

```bash
curl -X DELETE "$BASE/orders/<order-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "order deleted"
  }
}
```

Jika stok sudah dikurangi, stok dikembalikan sebelum order di-soft-delete.

## PUT /orders/:id/items

Mengganti seluruh item order. Hanya boleh saat status order masih `pending` dan stok belum dikurangi.

Request body:

```json
{
  "items": [
    {
      "variant_id": "66666666-6666-4666-8666-666666666666",
      "quantity": 3
    }
  ]
}
```

```bash
curl -X PUT "$BASE/orders/<order-id>/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"variant_id":"66666666-6666-4666-8666-666666666666","quantity":3}]}'
```

Response `200`: object order dengan item dan total terbaru.

Catatan:

- Item lama dihapus dan diganti dengan daftar baru.
- Harga item diambil ulang dari varian saat ini.
- `discount_amount` dan `shipping_cost` yang sudah ada dipertahankan.
- Jika order sudah `confirmed` atau stok sudah dikurangi -> `409`.

## PATCH /orders/:id/status

Request body:

```json
{
  "status": "confirmed"
}
```

```bash
curl -X PATCH "$BASE/orders/<order-id>/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

Transisi valid:

```text
pending -> confirmed | cancelled
confirmed -> processing | cancelled
processing -> shipped | cancelled
shipped -> completed | cancelled
```

- `confirmed` mengurangi stok dan menulis `stock_movements` tipe `order`. Jika stok tidak cukup -> `409`.
- `cancelled` mengembalikan stok jika sudah dikurangi dan menulis `stock_movements` tipe `cancel`.
- `completed` dan `cancelled` tidak dapat diubah lagi.

## PATCH /orders/:id/payment

Request body:

```json
{
  "payment_status": "paid",
  "payment_method": "transfer"
}
```

```bash
curl -X PATCH "$BASE/orders/<order-id>/payment" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status":"paid","payment_method":"transfer"}'
```

Nilai valid: `unpaid`, `paid`, `refunded`.

Transisi payment:

```text
unpaid -> paid
paid -> refunded | unpaid
refunded -> paid
```

Saat menjadi `paid`, `paid_at` diisi otomatis jika belum ada.

## PATCH /orders/:id/shipping

Request body:

```json
{
  "shipping_status": "picked_up",
  "courier": "JNE",
  "tracking_number": "JNE123456789"
}
```

```bash
curl -X PATCH "$BASE/orders/<order-id>/shipping" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"shipping_status":"picked_up","courier":"JNE","tracking_number":"JNE123456789"}'
```

Transisi valid:

```text
pending -> picked_up
picked_up -> in_transit
in_transit -> delivered | returned
delivered -> returned
```

- `picked_up`/`in_transit` mengisi `shipped_at` jika belum ada.
- `delivered` mengisi `delivered_at`.
- `returned` mengembalikan stok jika stok sudah dikurangi dan menulis `stock_movements` tipe `return`.
