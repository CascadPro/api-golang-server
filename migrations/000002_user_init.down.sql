DROP TRIGGER IF EXISTS trg_update_user_timestamp ON base.users;

DROP INDEX IF EXISTS idx_user_avatar;

DROP INDEX IF EXISTS idx_user_email;

DROP TABLE IF EXISTS base.users;

DROP TYPE IF EXISTS user_role;
