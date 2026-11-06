# API Accounting

Semua endpoint memakai permission `accounting.*`.

## GET /expenses

Query:

```text
?page=1&per_page=20&month=2026-11&category=restock&search=kain
```

```bash
curl "$BASE/expenses?page=1&per_page=20&month=2026-11&category=restock" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
        "title": "Belanja Bahan",
        "category": "restock",
        "vendor": "Toko Kain",
        "expense_date": "2026-11-04",
        "payment_method": "cash",
        "total_amount": 155000,
        "receipt_key": "accounting/receipts/<expense-id>/<uuid>.pdf",
        "receipt_url": "https://cdn.example.com/accounting/receipts/<expense-id>/<uuid>.pdf",
        "notes": "",
        "created_by_id": "3f0b8c1e-9f0d-4f5a-9f0d-111111111111",
        "items": [
          {
            "id": "cccccccc-cccc-4ccc-8ccc-cccccccccccc",
            "name": "Kain",
            "quantity": 2.5,
            "unit": "meter",
            "unit_price": 50000,
            "subtotal": 125000
          }
        ],
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

## POST /expenses

Request body dengan item:

```json
{
  "title": "Belanja Bahan",
  "category": "restock",
  "vendor": "Toko Kain",
  "expense_date": "2026-11-04",
  "payment_method": "cash",
  "notes": "Belanja bulanan",
  "items": [
    {
      "name": "Kain",
      "quantity": 2.5,
      "unit": "meter",
      "unit_price": 50000
    },
    {
      "name": "Benang",
      "quantity": 3,
      "unit": "pcs",
      "unit_price": 10000
    }
  ]
}
```

```bash
curl -X POST "$BASE/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Belanja Bahan",
    "category":"restock",
    "vendor":"Toko Kain",
    "expense_date":"2026-11-04",
    "payment_method":"cash",
    "notes":"Belanja bulanan",
    "items":[
      {"name":"Kain","quantity":2.5,"unit":"meter","unit_price":50000},
      {"name":"Benang","quantity":3,"unit":"pcs","unit_price":10000}
    ]
  }'
```

Request body tanpa item (total manual):

```json
{
  "title": "Biaya Listrik",
  "category": "operational",
  "expense_date": "2026-11-04",
  "payment_method": "transfer",
  "total_amount": 350000,
  "notes": "Tagihan listrik November"
}
```

```bash
curl -X POST "$BASE/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Biaya Listrik",
    "category":"operational",
    "expense_date":"2026-11-04",
    "payment_method":"transfer",
    "total_amount":350000,
    "notes":"Tagihan listrik November"
  }'
```

Catatan:

- Kategori valid: `operational`, `production`, `marketing`, `restock`, `other`.
- Jika `items` diisi, `total_amount` dihitung dari item dan nilai `total_amount` diabaikan.
- Jika `items` kosong, `total_amount` wajib > 0.
- `expense_date` format `YYYY-MM-DD`.

Response `201`: object expense.

## GET /expenses/:id

```bash
curl "$BASE/expenses/<expense-id>" \
  -H "Authorization: Bearer $TOKEN"
```

## PUT /expenses/:id

Request body:

```json
{
  "title": "Belanja Bahan Updated",
  "category": "restock",
  "vendor": "Toko Kain",
  "expense_date": "2026-11-05",
  "payment_method": "transfer",
  "notes": "Revisi",
  "items": [
    {
      "name": "Kain",
      "quantity": 3,
      "unit": "meter",
      "unit_price": 50000
    }
  ]
}
```

```bash
curl -X PUT "$BASE/expenses/<expense-id>" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Belanja Bahan Updated",
    "category":"restock",
    "vendor":"Toko Kain",
    "expense_date":"2026-11-05",
    "payment_method":"transfer",
    "notes":"Revisi",
    "items":[{"name":"Kain","quantity":3,"unit":"meter","unit_price":50000}]
  }'
```

## DELETE /expenses/:id

```bash
curl -X DELETE "$BASE/expenses/<expense-id>" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "expense deleted"
  }
}
```

## POST /expenses/:id/receipt

Upload bukti struk. Field multipart: `file`.

Format JPEG/PNG/PDF, maksimal 10 MB.

```bash
curl -X POST "$BASE/expenses/<expense-id>/receipt" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@struk.pdf"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "receipt_key": "accounting/receipts/<expense-id>/<uuid>.pdf",
    "receipt_url": "https://cdn.example.com/accounting/receipts/<expense-id>/<uuid>.pdf"
  }
}
```

Struk lama dihapus best-effort saat diganti.

## DELETE /expenses/:id/receipt

Menghapus bukti struk dari database dan R2.

```bash
curl -X DELETE "$BASE/expenses/<expense-id>/receipt" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "message": "receipt deleted"
  }
}
```

Jika expense belum memiliki struk -> `404`.

## GET /accounting/overview

Query:

```text
?month=2026-11
```

`month` opsional; default bulan berjalan di `APP_TIMEZONE`.

```bash
curl "$BASE/accounting/overview?month=2026-11" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "month": "2026-11",
    "total_expenses": 155000,
    "by_category": [
      {
        "category": "restock",
        "total": 155000
      }
    ],
    "recent_expenses": []
  }
}
```
