package core_postgres_token

const createQuery = `
	INSERT INTO base.tokens (id, token, type, expires_at, user_id)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, version, token, type, expires_at;
`

const deleteQuery = `
	DELETE FROM base.tokens
	WHERE id = $1;
`
