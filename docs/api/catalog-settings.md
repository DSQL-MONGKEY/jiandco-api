# Catalog Settings API

Mengatur konten halaman catalog publik. Memerlukan auth dan permission `products.view` / `products.update`.

## GET /api/v1/catalog/settings

```bash
curl "$BASE/catalog/settings" \
  -H "Authorization: Bearer $TOKEN"
```

Response:

```json
{
  "success": true,
  "data": {
    "promo_badge": "Harga Spesial WhatsApp",
    "promo_title": "Lebih hemat lewat WhatsApp.",
    "promo_description": "...",
    "promo_note": "...",
    "promo_cta_label": "Chat Admin",
    "page_eyebrow": "Catalog — Chongsam & Blush",
    "page_title": "The",
    "page_title_accent": "collection.",
    "page_description": "...",
    "page_footnote": "...",
    "cta_eyebrow": "The Collection",
    "cta_title": "Temukan koleksi",
    "cta_title_accent": "Ji&Co.",
    "cta_description": "...",
    "cta_label": "Jelajahi Katalog",
    "cta_image_url": "",
    "cta_secondary_image_url": "",
    "shopee_store_url": "https://shopee.co.id/"
  }
}
```

## PUT /api/v1/catalog/settings

Kirim seluruh field settings.

```bash
curl -X PUT "$BASE/catalog/settings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_badge": "Harga Spesial WhatsApp",
    "promo_title": "Lebih hemat lewat WhatsApp.",
    "promo_description": "...",
    "promo_note": "...",
    "promo_cta_label": "Chat Admin",
    "page_eyebrow": "Catalog — Chongsam & Blush",
    "page_title": "The",
    "page_title_accent": "collection.",
    "page_description": "...",
    "page_footnote": "...",
    "cta_eyebrow": "The Collection",
    "cta_title": "Temukan koleksi",
    "cta_title_accent": "Ji&Co.",
    "cta_description": "...",
    "cta_label": "Jelajahi Katalog",
    "cta_image_url": "",
    "cta_secondary_image_url": "",
    "shopee_store_url": "https://shopee.co.id/"
  }'
```
