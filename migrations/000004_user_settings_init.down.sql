DROP TRIGGER IF EXISTS trg_update_user_settings_timestamp ON base.user_settings;

DROP INDEX IF EXISTS idx_user_settings_user_id;

DROP TABLE IF EXISTS base.user_settings;

DROP TYPE IF EXISTS session_expire_term_enum;
