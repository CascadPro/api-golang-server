CREATE SCHEMA IF NOT EXISTS infrastructure;

CREATE TABLE infrastructure.outbox_events (
  id UUID PRIMARY KEY,

  type         TEXT  NOT NULL,
  aggregate_id UUID            DEFAULT NULL,
  payload      JSONB NOT NULL,

  attempts   INTEGER NOT NULL DEFAULT 0,
  last_error TEXT             DEFAULT NULL,

  processed_at TIMESTAMPTZ          DEFAULT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_outbox_events_pending
ON infrastructure.outbox_events (created_at, id)
WHERE processed_at IS NULL;
