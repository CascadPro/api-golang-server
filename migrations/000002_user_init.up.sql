CREATE TYPE user_role AS ENUM ('foreman', 'project_manager', 'clerk', 'engineer', 'director', 'regular', 'admin');

CREATE TABLE base.users (
  id         UUID             PRIMARY KEY,
  version    BIGINT  NOT NULL DEFAULT 1,
  activated  BOOLEAN NOT NULL DEFAULT false,

  email         VARCHAR(255)            CHECK(email ~ '^[A-Za-z0-9._-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
  password_hash VARCHAR(255),
  role          user_role     NOT NULL,

  name      VARCHAR(50)  NOT NULL CHECK (length(name) BETWEEN 2 AND 50),
  surname   VARCHAR(50)  NOT NULL CHECK (length(surname) BETWEEN 2 AND 50),
  last_name VARCHAR(50),          CHECK (length(last_name) BETWEEN 2 AND 50),

  last_active_at  TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
  created_at      TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMPTZ    NOT NULL   DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT chk_name_format CHECK (name ~* '^[[:alpha:]]+(?:[ -][[:alpha:]]+)*$'),
  CONSTRAINT chk_surname_format CHECK (surname ~* '^[[:alpha:]]+(?:[ -][[:alpha:]]+)*$'),
  CONSTRAINT chk_last_name_format CHECK (last_name IS NULL OR last_name ~* '^[[:alpha:]\s-]+$'),

  CONSTRAINT chk_activation CHECK (
    (activated AND email IS NOT NULL AND password_hash IS NOT NULL) OR
    (NOT activated AND email IS NULL AND password_hash IS NULL)
  ),

  CONSTRAINT chk_timestamps CHECK (
    (updated_at >= created_at AND last_active_at >= created_at)
  )
);

CREATE UNIQUE INDEX idx_user_email ON base.users (lower(email)) WHERE email IS NOT NULL;

CREATE INDEX idx_user_role ON base.users USING HASH (role);

CREATE TRIGGER trg_update_user_timestamp
BEFORE UPDATE ON base.users
FOR EACH ROW EXECUTE FUNCTION base.update_table_timestamp();
