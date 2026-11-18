package client_postgres_repository

const getManyQuery = `
	SELECT id, version, company, contacts, created_at, COUNT(*) OVER() AS total_count
	FROM base.clients
	WHERE (company ILIKE '%' || $1 || '%' OR $1 = ANY(contacts))
	ORDER BY company ASC, created_at DESC
	LIMIT $2
	OFFSET $3;
`

const getCountQuery = `
	SELECT COUNT(*)
	FROM base.clients
	WHERE (company ILIKE '%' || $1 || '%' OR $1 = ANY(contacts));
`

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
