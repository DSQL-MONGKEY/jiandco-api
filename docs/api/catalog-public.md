# Public Catalog API

Endpoint publik tanpa auth untuk halaman catalog `oinsh-page`.

## GET /api/v1/public/catalog

Response:

```json
{
  "success": true,
  "data": {
    "collections": [
      {
        "id": "chongsam",
        "label": "Chongsam",
        "description": "Lace vest..."
      }
    ],
    "products": [
      {
        "id": "uuid",
        "slug": "chongsam-cocoa",
        "name": "Cocoa Lace Vest",
        "collection": "chongsam",
        "color": "Cocoa",
        "whatsappPrice": 429000,
        "shopeePrice": 459000,
        "images": [
          {
            "src": "https://cdn.example.com/catalog/...",
            "alt": "Cocoa Lace Vest"
          }
        ],
        "description": "...",
        "details": ["Lace vest", "Button front"],
        "sizes": ["S", "M", "L", "XL"],
        "shopeeUrl": "https://shopee.co.id/...",
        "badge": "New",
        "badgeTone": "accent"
      }
    ],
    "promo": {
      "badge": "Harga Spesial WhatsApp",
      "title": "Lebih hemat lewat WhatsApp.",
      "description": "...",
      "note": "...",
      "ctaLabel": "Chat Admin"
    },
    "page": {
      "eyebrow": "Catalog — Chongsam & Blush",
      "title": "The",
      "titleAccent": "collection.",
      "description": "...",
      "footnote": "..."
    },
    "cta": {
      "eyebrow": "The Collection",
      "title": "Temukan koleksi",
      "titleAccent": "Ji&Co.",
      "description": "...",
      "label": "Jelajahi Katalog",
      "imageUrl": "",
      "secondaryImageUrl": ""
    },
    "shopeeStoreUrl": "https://shopee.co.id/"
  }
}
```

Hanya produk `status = active` yang dikembalikan.

## Checklist produk dashboard -> catalog

1. Status produk `active` (produk baru dari dashboard default `active`).
2. Pilih kategori agar produk punya collection.
3. Tambahkan minimal satu varian lengkap: size, color, SKU, dan harga > 0.
4. Upload minimal satu gambar dari halaman detail produk.
5. Pastikan origin `oinsh-page` ada di `CORS_ALLOWED_ORIGINS`.
6. Buka `/catalog` dan cek response `GET /api/v1/public/catalog`.

## GET /api/v1/public/catalog/products/:slug

Response berisi `product`, `collections`, `promo`, `page`, `cta`, dan `shopeeStoreUrl`.

## GET /api/v1/public/media/*

Menyajikan gambar dari R2 untuk bucket privat. Hanya key dengan prefix `catalog/`, `user/`, atau `accounting/` yang diizinkan.

## CORS

Tambahkan origin `oinsh-page` ke `CORS_ALLOWED_ORIGINS`, misalnya:

```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```
