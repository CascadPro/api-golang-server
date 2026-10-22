package idempotency_postgres_repository

const reserveQuery = `
INSERT INTO infrastructure.idempotency_records (
  scope,
  operation,
  key_hash,
  request_hash,
  status,
  expires_at
)
VALUES ($1, $2, $3, $4, 'processing', $5)
ON CONFLICT (scope, operation, key_hash)
DO UPDATE SET
  request_hash = EXCLUDED.request_hash,
  status = 'processing',
  response_status = NULL,
  response_body = NULL,
  response_headers = NULL,
  created_at = CURRENT_TIMESTAMP,
  updated_at = CURRENT_TIMESTAMP,
  expires_at = EXCLUDED.expires_at
WHERE infrastructure.idempotency_records.expires_at <= CURRENT_TIMESTAMP
RETURNING request_hash, status, response_status, response_body, response_headers
`

const completeQuery = `
UPDATE infrastructure.idempotency_records
SET status = 'completed', response_status = $4, response_body = $5, response_headers = $6, updated_at = CURRENT_TIMESTAMP
WHERE scope = $1 AND operation = $2 AND key_hash = $3 AND status = 'processing'
`

const scanRecordQuery = `
SELECT request_hash, status, response_status, response_body, response_headers
FROM infrastructure.idempotency_records
WHERE scope = $1 AND operation = $2 AND key_hash = $3
`

const cleanupExpiredQuery = `
DELETE FROM infrastructure.idempotency_records
WHERE expires_at <= CURRENT_TIMESTAMP
  AND (
    status = 'completed'
    OR (
      status = 'processing'
      AND updated_at <= CURRENT_TIMESTAMP - INTERVAL '1 hour'
    )
  )
`
