CREATE SCHEMA media;

CREATE TYPE file_tag AS ENUM ('avatars', 'images', 'videos', 'docs');

CREATE TABLE media.files (
  id           VARCHAR(48)            PRIMARY KEY,
  version      BIGINT       NOT NULL  DEFAULT 1,
  tag          file_tag     NOT NULL,
  filename     VARCHAR(255) NOT NULL  CHECK(char_length(filename) BETWEEN 1 AND 255),
  content_type VARCHAR(255) NOT NULL,
  size         BIGINT       NOT NULL  CHECK(size > 0),
  deleted      BOOLEAN      NOT NULL  DEFAULT false,
  deleted_at   TIMESTAMPTZ            DEFAULT NULL,
  created_at   TIMESTAMPTZ  NOT NULL  DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMPTZ  NOT NULL  DEFAULT CURRENT_TIMESTAMP,

  CHECK (
    (deleted AND deleted_at IS NOT NULL)
    OR
    NOT deleted
    AND
    (updated_at >= created_at AND deleted_at >= created_at)
  )
);

ALTER TABLE base.users ADD COLUMN IF NOT EXISTS avatar_file_id VARCHAR(48) UNIQUE DEFAULT NULL;

ALTER TABLE base.users
ADD CONSTRAINT fk_user_file_id
FOREIGN KEY (avatar_file_id) REFERENCES media.files(id)
ON DELETE SET NULL ON UPDATE CASCADE;

CREATE UNIQUE INDEX idx_user_avatar_file_id ON base.users (avatar_file_id);

CREATE INDEX idx_file_filename ON media.files (filename);

CREaTE INDEX idx_file_content_type ON media.files (content_type);

CREATE TRIGGER trg_update_file_timestamp
BEFORE UPDATE ON media.files
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
