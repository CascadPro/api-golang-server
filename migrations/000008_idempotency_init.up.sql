CREATE TYPE idempotency_status AS ENUM ('in_progress', 'completed');

CREATE TABLE infrastructure.idempotency_records (
    id UUID PRIMARY KEY,

    scope        TEXT  NOT NULL,
    operation    TEXT  NOT NULL,
    key_hash     BYTEA NOT NULL,
    request_hash BYTEA NOT NULL,

    status idempotency_status NOT NULL DEFAULT 'in_progress',

    response_status  SMALLINT,
    response_body    BYTEA,
    response_headers JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT pk_idempotency_records PRIMARY KEY (scope, operation, key_hash),

    CONSTRAINT chk_idempotency_response CHECK (
        (
            status = 'processing'
            AND response_status IS NULL
            AND response_body IS NULL
        )
        OR
        (
            status = 'completed'
            AND response_status IS NOT NULL
            AND response_body IS NOT NULL
        )
    )
);

CREATE INDEX idx_idempotency_records_expires_at ON infrastructure.idempotency_records (expires_at);

CREATE TRIGGER trg_update_idempotency_records_timestamp
BEFORE UPDATE ON infrastructure.idempotency_records
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
