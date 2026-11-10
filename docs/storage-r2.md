# Storage R2

## Konfigurasi

```env
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET=jiandco
R2_ENDPOINT=
R2_PUBLIC_BASE_URL=https://cdn.jiandco.id
```

- `R2_ENDPOINT` opsional. Jika kosong, diturunkan menjadi `https://{R2_ACCOUNT_ID}.r2.cloudflarestorage.com`.
- `R2_PUBLIC_BASE_URL` dipakai untuk membentuk URL response. Jika kosong, response mengembalikan object key.
- R2 bersifat opsional di development. Upload mengembalikan `503` jika konfigurasi belum lengkap.

Client memakai AWS SDK v2 S3 dengan region `auto` dan static credentials.

## Object key

Database hanya menyimpan object key.

```text
user/avatars/{userID}/{uuid}.{ext}
catalog/products/{productID}/{uuid}.{ext}
catalog/products/{productID}/colors/{colorID}/{uuid}.{ext}
accounting/receipts/{expenseID}/{uuid}.{ext}
```

## Batasan upload

| Jenis | Format | Maksimal |
|---|---|---|
| Avatar | JPEG, PNG, WebP | 10 MB |
| Gambar produk | JPEG, PNG, WebP | 10 MB per file |
| Struk | JPEG, PNG, PDF | 10 MB |

MIME dideteksi dari isi file (`http.DetectContentType`), bukan hanya ekstensi.

## Alur upload

1. Client mengirim `multipart/form-data` ke backend.
2. Handler membaca file header dan ukuran.
3. Service memvalidasi MIME dan ukuran.
4. Service membentuk key dengan UUID.
5. Service upload ke R2.
6. Service menyimpan key ke database.
7. Jika update database gagal, object baru dihapus best-effort.
8. Saat file diganti, object lama dihapus best-effort.

## URL

- Bucket publik: response memakai `R2_PUBLIC_BASE_URL + "/" + object_key`.
- Untuk Ji&Co production, `R2_PUBLIC_BASE_URL=https://cdn.jiandco.id`.
- Bucket privat: gunakan presigned GET di fase berikutnya; jangan menyimpan URL sementara di database.

## Aturan

- Jangan menyimpan full URL di database.
- Jangan menaruh credentials R2 di repository.
- Jangan menimpa object lama dengan key yang sama; selalu buat UUID baru.
- Hapus object hanya setelah operasi database berhasil.
