DROP TRIGGER IF EXISTS trg_update_idempotency_records_timestamp ON infrastructure.idempotency_records;

DROP INDEX IF EXISTS idx_idempotency_records_expires_at;

DROP TABLE IF EXISTS infrastructure.idempotency_records;

DROP TYPE IF EXISTS idempotency_status;
