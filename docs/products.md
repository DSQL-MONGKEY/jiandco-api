# Products



Status: selesai. Referensi endpoint dan contoh payload ada di [docs/api/products.md](api/products.md) dan [docs/api/catalog.md](api/catalog.md).



## Model



```text

categories 1─N products 1─N product_images

products 1─N product_variants N─1 sizes

products 1─N product_colors 1─N product_variants

```



- Harga dan stok berada di level varian.

- Satu varian adalah kombinasi unik `product + size + product_color`. Warna bersifat custom per produk.

- SKU unik per varian.

- Produk memiliki banyak gambar; satu gambar ditandai primary di service.

- Produk memakai soft delete. Soft delete tidak menghapus child rows; hard delete produk akan menghapus gambar/varian lewat FK cascade.

- Varian yang sudah dipakai order tidak dihapus; nonaktifkan dengan `is_active=false`.



## Status produk



`draft`, `active`, `archived`.

Produk baru dari dashboard default `active` agar langsung tampil di catalog publik. API
`POST /products` tanpa field `status` tetap default `draft`; kirim `status: "active"` untuk
langsung mempublikasikan produk.



## Endpoint



```text

GET    /api/v1/products

POST   /api/v1/products

GET    /api/v1/products/:id

PUT    /api/v1/products/:id

DELETE /api/v1/products/:id



POST   /api/v1/products/:id/colors
PUT    /api/v1/products/:id/colors/:colorId
DELETE /api/v1/products/:id/colors/:colorId
POST   /api/v1/products/:id/colors/:colorId/image
DELETE /api/v1/products/:id/colors/:colorId/image
POST   /api/v1/products/:id/images

DELETE /api/v1/products/:id/images/:imageId

PUT    /api/v1/products/:id/images/:imageId/primary



POST   /api/v1/products/:id/variants

PUT    /api/v1/products/:id/variants/:variantId

DELETE /api/v1/products/:id/variants/:variantId

PATCH  /api/v1/products/:id/variants/:variantId/stock



GET/POST/PUT/DELETE /api/v1/categories

GET/POST/PUT/DELETE /api/v1/sizes

GET/POST/PUT/DELETE /api/v1/colors (legacy, tidak dipakai UI produk)

```



Semua endpoint memakai permission `products.*`.



## Upload gambar



- `multipart/form-data`, field `images` (bisa lebih dari satu).

- Format: JPEG, PNG, WebP.

- Maksimal 10 MB per file.

- Object key: `catalog/products/{productID}/{uuid}.{ext}`.
- Gambar warna: `catalog/products/{productID}/colors/{colorID}/{uuid}.{ext}`.
- URL publik memakai `R2_PUBLIC_BASE_URL`, production `https://cdn.jiandco.id`.

- Database menyimpan `object_key`, bukan URL.

- Hapus object R2 best-effort saat gambar dihapus dari produk.



## Catalog publik



Field yang dipakai halaman catalog `oinsh-page`:



- `shopee_url`, `shopee_price`

- `details` (jsonb array)

- `badge`, `badge_tone`

- `sort_order`

- `status = active`



Satu produk dashboard mewakili satu kartu warna catalog; varian mewakili ukuran. Endpoint publik dan settings dijelaskan di `docs/api/catalog-public.md` dan `docs/api/catalog-settings.md`.



## Stok dan audit



- Stok disimpan di `product_variants.stock`.

- Setiap perubahan stok dicatat di `stock_movements`:

  - `initial` saat produk/varian dibuat dengan stok awal,

  - `adjustment` saat update stok manual,

  - `order` saat order dikonfirmasi,

  - `cancel` saat order dibatalkan/dihapus,

  - `return` saat pengiriman returned.

- Riwayat dapat dilihat di `GET /products/:id/variants/:variantId/movements`.

- Stok tidak boleh negatif; database memiliki check `stock >= 0`.



## Listing



Filter:



```text

?search=&category_id=&status=&page=&per_page=

```



Response menyertakan kategori, gambar (urut `sort_order`), dan varian beserta size/color.

