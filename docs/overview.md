# Overview

Status: selesai. Referensi endpoint dan contoh payload ada di [docs/api/overview.md](api/overview.md).

## Endpoint

```text
GET /api/v1/overview/revenue?month=YYYY-MM
```

- `month` opsional.
- Default: bulan berjalan pada `APP_TIMEZONE` (default `Asia/Jakarta`).
- Format salah -> `400`.

## Definisi pendapatan

Pendapatan dihitung secara cash basis:

```text
revenue = SUM(orders.total_amount)
WHERE payment_status = 'paid'
  AND paid_at >= start_of_month
  AND paid_at < start_of_next_month
  AND deleted_at IS NULL
```

Batas bulan dihitung di timezone aplikasi, lalu dikonversi ke UTC untuk query `timestamptz`.

Pengeluaran:

```text
expenses = SUM(expenses.total_amount)
WHERE expense_date >= start_date
  AND expense_date <= end_date
  AND deleted_at IS NULL
```

Laba bersih:

```text
net_profit = revenue - expenses
```

## Response

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
    "order_status_breakdown": [
      { "status": "pending", "count": 4 },
      { "status": "completed", "count": 12 }
    ],
    "daily": [
      { "date": "2026-11-01", "revenue": 500000, "orders": 2 }
    ]
  }
}
```

## Catatan

- Order `unpaid`/`refunded` tidak masuk revenue.
- Order `cancelled` tidak masuk revenue.
- Jika nanti dibutuhkan basis accrual, gunakan `created_at` atau tanggal order dan dokumentasikan perubahannya.
- Permission: `overview.view`.
