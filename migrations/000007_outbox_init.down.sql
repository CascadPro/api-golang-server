DROP INDEX IF EXISTS idx_outbox_events_locked;

DROP INDEX IF EXISTS idx_outbox_events_pending;

DROP TABLE IF EXISTS infrastructure.outbox_events;

DROP TYPE IF EXISTS outbox_events_type;

DROP SCHEMA IF EXISTS infrastructure;
