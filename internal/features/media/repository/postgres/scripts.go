package media_postgres_repository

const getQuery = `
	SELECT id, version, tag, filename, mime_type, size, deleted, deleted_at, created_at, placeholder
	FROM media.files
	WHERE (id = $1)
	LIMIT 1;
`

const createQuery = `
	INSERT INTO media.files (id, tag, filename, mime_type, size, placeholder)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id, version, tag, filename, mime_type, size, created_at;
`

const deleteQuery = `
	DELETE FROM media.files
	WHERE id = $1;
`

const patchQuery = `
	UPDATE media.files
	SET version = version + 1, filename = $1, size = $2, deleted = $3, deleted_at = $4
	WHERE (id = $5 AND version = $6)
	RETURNING id, version, tag, filename, mime_type, size, deleted, deleted_at, created_at;
`

const patchPlaceholderQuery = `
	UPDATE media.files
	SET version = version + 1, placeholder = $1
	WHERE (id = $2 AND version = $3);
`
