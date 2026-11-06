# API Overview

## GET /overview/revenue

Permission: `overview.view`.

Query:

```text
?month=2026-11
```

`month` opsional; default bulan berjalan di `APP_TIMEZONE`.

```bash
curl "$BASE/overview/revenue?month=2026-11" \
  -H "Authorization: Bearer $TOKEN"
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "month": "2026-11",
    "revenue": 12500000,
    "expenses": 4200000,
    "net_profit": 8300000,
    "total_orders": 42,
    "paid_orders": 35,
    "average_order_value": 357142,
    "daily": [
      {
        "date": "2026-11-01",
        "revenue": 500000,
        "orders": 2
      }
    ]
  }
}
```

Definisi:

- `revenue`: total order `paid` berdasarkan `paid_at` dalam bulan tersebut.
- `expenses`: total expense berdasarkan `expense_date` dalam bulan tersebut.
- `net_profit`: `revenue - expenses`.
- `total_orders`: jumlah order yang dibuat dalam bulan tersebut.
- `paid_orders`: jumlah order yang lunas dalam bulan tersebut.
- `average_order_value`: `revenue / paid_orders`.
- `daily`: agregasi pendapatan per tanggal pembayaran (timezone aplikasi).
