DROP TRIGGER IF EXISTS trg_update_file_timestamp ON media.files;

DROP INDEX IF EXISTS idx_file_created_at;

DROP INDEX IF EXISTS idx_file_deleted;

DROP INDEX IF EXISTS idx_file_tag;

DROP INDEX IF EXISTS idx_user_avatar_file_id;

ALTER TABLE base.users DROP CONSTRAINT IF EXISTS fk_user_file_id;

ALTER TABLE base.users DROP COLUMN IF EXISTS avatar_file_id;

DROP TABLE IF EXISTS media.files;

DROP TYPE IF EXISTS file_tag;

DROP SCHEMA IF EXISTS media;
