package core_postgres_outbox

const cleanupQuery = `
DELETE FROM infrastructure.outbox_events
WHERE (
	processed_at IS NOT NULL
	AND last_error IS NULL
);
`
