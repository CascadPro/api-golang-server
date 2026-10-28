package settings_postgres_repository

import (
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

// type SettingsRepository interface {
// 	CreateUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
// 	GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
// 	PatchUserSettings(ctx context.Context, id uuid.UUID, patch domain.UserSettingsPatch) (domain.UserSettings, error)
// }

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{pool: pool}
}
