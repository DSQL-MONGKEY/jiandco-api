# Accounting

Status: selesai. Referensi endpoint dan contoh payload ada di [docs/api/accounting.md](api/accounting.md).

## Model

Satu ledger pengeluaran:

```text
expenses 1─N expense_items
```

- `expenses` menyimpan header: judul, kategori, vendor, tanggal, metode bayar, total, struk.
- `expense_items` menyimpan rincian belanja.
- Jika item dikirim, `total_amount` dihitung dari jumlah `subtotal` item.
- Jika tidak ada item, `total_amount` diisi manual.
- `subtotal = round(quantity * unit_price)` dalam rupiah.
- Expense memakai soft delete; soft delete tidak menghapus item. Hard delete expense akan menghapus item lewat FK cascade.

## Kategori

`operational`, `production`, `marketing`, `restock`, `other`.

## Struk

- Upload `multipart/form-data`, field `file`.
- Format: JPEG, PNG, PDF.
- Maksimal 10 MB.
- Object key: `accounting/receipts/{expenseID}/{uuid}.{ext}`.
- Struk lama dihapus best-effort saat diganti.
- `DELETE /expenses/:id/receipt` menghapus struk dari database dan R2.

## Endpoint

```text
GET    /api/v1/expenses
POST   /api/v1/expenses
GET    /api/v1/expenses/:id
PUT    /api/v1/expenses/:id
DELETE /api/v1/expenses/:id
POST   /api/v1/expenses/:id/receipt

GET    /api/v1/accounting/overview?month=YYYY-MM
```

Semua endpoint memakai permission `accounting.*`.

## Accounting overview

`GET /api/v1/accounting/overview?month=YYYY-MM` (default bulan berjalan di `APP_TIMEZONE`) mengembalikan:

```json
{
  "month": "2026-11",
  "total_expenses": 1250000,
  "by_category": [
    { "category": "operational", "total": 500000 },
    { "category": "restock", "total": 750000 }
  ],
  "recent_expenses": []
}
```

Filter tanggal memakai `expense_date` (tanggal bisnis), bukan timestamp.
