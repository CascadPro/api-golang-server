DROP INDEX IF EXISTS base.idx_user_phone;

ALTER TABLE base.users DROP COLUMN IF EXISTS phone;
