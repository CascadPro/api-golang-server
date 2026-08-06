CREATE TABLE base.clients (
  id      UUID            PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 1,

  company  VARCHAR(255)   NOT NULL CHECK(length(company) BETWEEN 1 AND 255),
  contacts VARCHAR(255)[] NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_client_company ON base.clients (company);
