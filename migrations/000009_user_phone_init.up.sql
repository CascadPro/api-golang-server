ALTER TABLE base.users ADD COLUMN phone VARCHAR(20) CHECK (phone ~ '^\+\d{1,15}$');

CREATE UNIQUE INDEX idx_user_phone ON base.users (phone) WHERE phone IS NOT NULL;
