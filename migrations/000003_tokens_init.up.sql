CREATE TYPE token_type AS ENUM ('register', 'email_verify');

CREATE TABLE base.tokens (
  id          UUID                    PRIMARY KEY,
  version     BIGINT        NOT NULL  DEFAULT 1,
  token       VARCHAR(255)  NOT NULL,
  type        token_type    NOT NULL,
  expires_at  TIMESTAMPTZ   NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL  DEFAULT CURRENT_TIMESTAMP,

  user_id     UUID          NOT NULL  REFERENCES base.users (id) ON DELETE CASCADE,

  CHECK (expires_at >= created_at)
);

CREATE UNIQUE INDEX idx_user_token ON base.tokens (token, user_id);

CREATE UNIQUE INDEX idx_user_token_type_active ON base.tokens (user_id, type);
