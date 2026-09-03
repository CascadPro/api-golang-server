CREATE SCHEMA IF NOT EXISTS infrastructure;

CREATE TYPE outbox_events_type AS ENUM ('media.delete_file');

CREATE TABLE infrastructure.outbox_events (
  id UUID PRIMARY KEY,

  type         outbox_events_type NOT NULL,
  aggregate_id UUID                         DEFAULT NULL,
  payload      JSONB              NOT NULL,

  attempts   INTEGER NOT NULL DEFAULT 0,
  last_error TEXT             DEFAULT NULL,

  processed_at TIMESTAMPTZ          DEFAULT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_outbox_events_pending
ON infrastructure.outbox_events (created_at, id)
WHERE processed_at IS NULL;
