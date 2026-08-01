CREATE TYPE token_type AS ENUM ('register', 'email_verify');

CREATE TABLE base.tokens (
  id      UUID             PRIMARY KEY,
  user_id UUID   NOT NULL,
  version BIGINT NOT NULL  DEFAULT 1,

  token VARCHAR(255) NOT NULL  CHECK (token ~ '^[A-Fa-f0-9-]+$'),
  type  token_type   NOT NULL,

  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_token_user_id FOREIGN KEY (user_id) REFERENCES base.users(id) ON DELETE CASCADE,

  CONSTRAINT chk_token_timestamps CHECK (expires_at >= created_at)
);

CREATE UNIQUE INDEX idx_user_token ON base.tokens (token, user_id);

CREATE UNIQUE INDEX idx_user_token_type ON base.tokens (user_id, type);
