# API Conventions



Referensi endpoint lengkap beserta contoh payload ada di [docs/api/README.md](api/README.md).



## Base URL



```text

/api/v1

```



Health check berada di luar prefix: `GET /healthz`, `GET /readyz`.



## Envelope sukses



```json

{

  "success": true,

  "data": {}

}

```



Untuk list:



```json

{

  "success": true,

  "data": {

    "items": [],

    "pagination": {

      "page": 1,

      "per_page": 20,

      "total": 0,

      "total_pages": 0

    }

  }

}

```



## Envelope error



```json

{

  "success": false,

  "code": 400,

  "message": "validation failed"

}

```



`code` selalu sama dengan HTTP status. Pesan error tidak boleh memuat stack trace, password, token, atau detail internal database.



Kode yang dipakai:



| Code | Arti |

|---|---|

| 400 | request tidak valid / validasi gagal |

| 401 | belum login, token invalid, session revoked/expired |

| 403 | tidak punya permission atau user inactive |

| 404 | resource/route tidak ditemukan |

| 409 | konflik unique (username/email/slug/SKU/order number) |

| 422 | validasi bisnis khusus (jika diperlukan) |

| 429 | rate limit |

| 500 | error internal |

| 503 | storage/R2 belum dikonfigurasi |



## Pagination dan filter



Query umum:



```text

?page=1&per_page=20&search=keyword

```



- `page` default 1.

- `per_page` default 20, maksimum 100.

- `search` opsional dan dibatasi per modul.

- Filter lain mengikuti dokumentasi modul.



## Autentikasi



```http

Authorization: Bearer <access_token>

```



Token akses berlaku 1 jam. Refresh memakai endpoint terpisah dan tidak memakai header Authorization.



## Format data



- Tanggal/waktu response: RFC3339 dengan timezone, contoh `2026-11-04T15:00:00+07:00`.

- Database menyimpan UTC.

- Filter bulan memakai format `YYYY-MM` dan timezone `APP_TIMEZONE` (default `Asia/Jakarta`).

- Uang: integer rupiah, contoh `150000`.

- ID: UUID string.

- Field kosong yang opsional dikirim sebagai `""` atau `null` sesuai DTO.



## Upload file



Endpoint upload memakai `multipart/form-data` dengan field `file` (avatar) atau `images` (katalog, fase berikutnya).



Batasan:



- Avatar: JPEG/PNG/WebP, maksimal 10 MB.

- Gambar produk: JPEG/PNG/WebP, maksimal 10 MB per file.
- Gambar preview warna: JPEG/PNG/WebP, maksimal 10 MB, satu gambar per warna.

- Struk: JPEG/PNG/PDF, maksimal 10 MB.



Response upload mengembalikan object key dan URL publik (jika R2 dikonfigurasi).



## Versioning



Perubahan breaking harus menaikkan versi path (`/api/v2`). Perubahan non-breaking (field baru opsional, filter baru) tetap di `/api/v1`.
