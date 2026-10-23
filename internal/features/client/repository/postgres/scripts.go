package client_postgres_repository

const getQuery = `
	SELECT id, version, company, contacts, created_at
	FROM base.clients
	WHERE (id = $1)
	LIMIT 1;
`

const createQuery = `
	INSERT INTO base.clients (id, company, contacts)
	VALUES($1, $2, $3)
	RETURNING id, version, company, contacts, created_at;
`

const deleteQuery = `
	DELETE FROM base.clients
	WHERE (id = $1);
`
