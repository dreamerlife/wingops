# CMDB Asset Sync API

`POST /api/v1/cmdb/assets/sync` accepts signed asset sync requests.

Development credentials:

- `X-Api-Key`: `dev-sync-key`
- secret: `dev-sync-secret`

`X-Signature` is the lowercase hex HMAC-SHA256 of the raw JSON request body using the secret.

```json
{
  "model_id": "1",
  "sync_mode": "incremental",
  "assets": [
    {
      "unique_key": "sn:ABC123456",
      "attributes": {
        "hostname": "web-server-01",
        "management_ip": "10.0.1.101"
      }
    }
  ]
}
```

`incremental` and `full` currently both upsert incoming assets by `(model_id, unique_key)`.
