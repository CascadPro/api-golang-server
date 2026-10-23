package core_postgres_outbox

const getPendingQuery = `
	WITH events AS (
		SELECT id
		FROM infrastructure.outbox_events
		WHERE (
			processed_at IS NULL
			AND attempts < 4
			AND (
				locked_at IS NULL
				OR locked_at < NOW() - ($3 * INTERVAL '1 second')
			)
		)
		ORDER BY created_at, id
		FOR UPDATE SKIP LOCKED
		LIMIT $1
	)
	UPDATE infrastructure.outbox_events AS e
	SET locked_at = NOW(), locked_by = $2, attempts = attempts + 1
	FROM events
	WHERE e.id = events.id
	RETURNING
		e.id,
		e.type,
		e.aggregate_id,
		e.payload,
		e.attempts,
		e.last_error,
		e.locked_at,
		e.locked_by,
		e.processed_at,
		e.created_at;
`

const createQuery = `
	INSERT INTO infrastructure.outbox_events (id, type, aggregate_id, payload)
	VALUES ($1, $2, $3, $4)
	RETURNING id, type, aggregate_id, payload, created_at;
`

const markFailedQuery = `
	UPDATE infrastructure.outbox_events
	SET last_error = $3
	WHERE (id = $1 AND locked_by = $2 AND processed_at IS NULL);
`

const markProcessedQuery = `
	UPDATE infrastructure.outbox_events
	SET processed_at = $3, locked_at = NULL, locked_by = NULL, last_error = NULL
	WHERE (id = $1 AND locked_by = $2 AND processed_at IS NULL);
`

const cleanupQuery = `
	DELETE FROM infrastructure.outbox_events
	WHERE (
		processed_at IS NOT NULL
		AND last_error IS NULL
	);
`
