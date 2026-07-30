CREATE TYPE user_role AS ENUM ('foreman', 'project_manager', 'clerk', 'engineer', 'director', 'regular', 'admin');

CREATE TABLE base.users (
  id              UUID                      PRIMARY KEY,
  version         BIGINT         NOT NULL   DEFAULT 1,
  activated       BOOLEAN        NOT NULL   DEFAULT false,
  email           VARCHAR(255)              UNIQUE,
  password_hash   VARCHAR(255),
  role            user_role      NOT NULL,
  name            VARCHAR(100)   NOT NULL   CHECK (char_length(name) BETWEEN 2 AND 100),
  surname         VARCHAR(100)   NOT NULL   CHECK (char_length(surname) BETWEEN 2 AND 100),
  last_name       VARCHAR(100),             CHECK (char_length(last_name) BETWEEN 2 AND 100),
  last_active_at  TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
  created_at      TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,

  CHECK (
    (activated AND email IS NOT NULL AND password_hash IS NOT NULL)
    OR
    (NOT activated AND email IS NULL AND password_hash IS NULL)
    AND
    (updated_at >= created_at AND last_active_at >= created_at)
  )
);

CREATE UNIQUE INDEX idx_user_email ON base.users (email);

CREATE TRIGGER trg_update_user_timestamp
BEFORE UPDATE ON base.users
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
