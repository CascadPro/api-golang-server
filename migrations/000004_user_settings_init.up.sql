CREATE TYPE session_expire_term_enum AS ENUM ('7d', '30d', '90d');

CREATE TABLE base.user_settings (
  id                  UUID                               PRIMARY KEY,
  version             BIGINT                   NOT NULL  DEFAULT 1,
  session_expire_term session_expire_term_enum NOT NULL  DEFAULT '30d',
  created_at          TIMESTAMPTZ              NOT NULL  DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMPTZ              NOT NULL  DEFAULT CURRENT_TIMESTAMP,

  user_id             UUID                     NOT NULL  UNIQUE REFERENCES base.users (id) ON DELETE CASCADE,

  CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX idx_user_settings_user_id ON base.user_settings (user_id);

CREATE TRIGGER trg_update_user_settings_timestamp
BEFORE UPDATE ON base.user_settings
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
