# Health

## GET /healthz

Cek proses API hidup. Tidak memerlukan auth.

```bash
curl http://localhost:8080/healthz
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

## GET /readyz

Cek koneksi database. Tidak memerlukan auth.

```bash
curl http://localhost:8080/readyz
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "status": "ready"
  }
}
```

Response `503` jika database tidak dapat dihubungi:

```json
{
  "success": false,
  "code": 503,
  "message": "database is not ready"
}
```
