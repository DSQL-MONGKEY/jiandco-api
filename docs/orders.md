# Orders

Status: selesai. Referensi endpoint dan contoh payload ada di [docs/api/orders.md](api/orders.md).

## Model

```text
customers 1─N orders 1─N order_items
products/variants 1─N order_items (nullable)
```

- `customers` adalah master pembeli.
- `orders` menyimpan snapshot penerima dan alamat agar histori tidak berubah saat master customer diedit.
- `order_items` menyimpan snapshot nama produk, label varian, SKU, dan harga saat order dibuat.
- Order memakai soft delete; soft delete tidak menghapus item. Hard delete order akan menghapus item lewat FK cascade.
- `orders.stock_deducted` menandai apakah stok sudah dikurangi, sehingga restore stok tidak terjadi dua kali.

## Status

### Order

`pending -> confirmed -> processing -> shipped -> completed`

Pembatalan dapat dilakukan dari status selain `completed` dan `cancelled`.

### Payment

`unpaid -> paid -> refunded`

`paid_at` diisi saat status berubah menjadi `paid`.

### Shipping

`pending -> picked_up -> in_transit -> delivered`

`returned` dapat terjadi setelah `shipped`/`in_transit`.

## Resi dan pengiriman

Field:

- `courier`
- `tracking_number` (resi)
- `shipped_at`
- `delivered_at`

Perubahan resi dan status pengiriman memakai endpoint update order.

## Edit item

- `PUT /orders/:id/items` mengganti seluruh item order.
- Hanya boleh saat status `pending` dan stok belum dikurangi.
- Item baru memakai harga varian saat ini.
- Setelah order `confirmed`, item tidak dapat diubah.

## Nomor order

Format: `JDC-YYYYMMDD-XXXXXX` dengan suffix acak. Unique index `order_number` menjaga duplikasi. Service melakukan retry saat terjadi konflik.

## Stok

- Stok dikurangi saat order berubah menjadi `confirmed`.
- Stok dikembalikan saat order `cancelled` atau shipping `returned`.
- Perubahan stok dan order harus berada dalam satu transaksi.
- Jangan mengizinkan stok negatif.

## Total

```text
subtotal      = sum(order_items.subtotal)
total_amount  = subtotal - discount_amount + shipping_cost
```

Semua nominal `bigint` rupiah.

## Endpoint

```text
GET    /api/v1/orders
POST   /api/v1/orders
GET    /api/v1/orders/:id
PUT    /api/v1/orders/:id
DELETE /api/v1/orders/:id

PATCH  /api/v1/orders/:id/status
PATCH  /api/v1/orders/:id/payment
PATCH  /api/v1/orders/:id/shipping

GET/POST/PUT/DELETE /api/v1/customers
```

Filter listing:

```text
?status=&payment_status=&shipping_status=&search=&from=&to=&page=&per_page=
```

`search` mencakup nomor order, nama pembeli, nomor telepon, dan resi.

Semua endpoint memakai permission `orders.*` atau `customers.*`.
