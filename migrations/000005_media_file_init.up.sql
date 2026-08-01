CREATE SCHEMA media;

CREATE TYPE file_tag AS ENUM ('avatars', 'images', 'videos', 'docs');

CREATE TABLE media.files (
  id           VARCHAR(48)            PRIMARY KEY,
  version      BIGINT       NOT NULL  DEFAULT 1,
  tag          file_tag     NOT NULL,
  filename     VARCHAR(255) NOT NULL  CHECK(char_length(filename) BETWEEN 1 AND 255),
  mime_type    VARCHAR(255) NOT NULL  CHECK(mime_type ~ '^[a-zA-Z0-9]+\/[a-zA-Z0-9\.\-\+]+$'),
  placeholder  BYTEA                  DEFAULT NULL,
  size         BIGINT       NOT NULL  CHECK(size > 0),
  deleted      BOOLEAN      NOT NULL  DEFAULT false,
  deleted_at   TIMESTAMPTZ            DEFAULT NULL,
  created_at   TIMESTAMPTZ  NOT NULL  DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMPTZ  NOT NULL  DEFAULT CURRENT_TIMESTAMP,

  CHECK (
    (deleted AND deleted_at IS NOT NULL) OR NOT deleted
    AND
    ((tag IN ('avatars', 'images', 'videos') AND placeholder IS NOT NULL) OR
    (tag NOT IN ('avatars', 'images', 'videos') AND placeholder IS NULL))
    AND
    tag != 'avatars' OR (tag = 'avatars' AND size <= 10485760)
    AND
    (updated_at >= created_at AND deleted_at >= created_at)
  )
);

ALTER TABLE base.users ADD COLUMN IF NOT EXISTS avatar_file_id VARCHAR(48) DEFAULT NULL;

ALTER TABLE base.users
ADD CONSTRAINT fk_user_file_id
FOREIGN KEY (avatar_file_id) REFERENCES media.files(id)
ON DELETE SET NULL ON UPDATE CASCADE;

CREATE UNIQUE INDEX idx_user_avatar_file_id ON base.users (avatar_file_id);

CREATE INDEX idx_file_tag ON media.files USING HASH (tag);

CREATE INDEX idx_file_deleted ON media.files USING HASH (deleted);

CREATE INDEX idx_file_created_at ON media.files (created_at DESC);

CREATE TRIGGER trg_update_file_timestamp
BEFORE UPDATE ON media.files
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
