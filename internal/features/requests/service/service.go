package requests_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	"github.com/google/uuid"
)

type Service struct {
	mediaService      MediaService
	cache             core_redis_cache.Cache
	clientService     ClientService
	usersPostgresRepo UsersRepository
	requestsMongoRepo RequestsRepository
}

type MediaService interface {
	GetFile(ctx context.Context, fileID string) (file domain.File, err error)

	UploadFile(ctx context.Context, file *domain.File, content []byte) (created domain.File, err error)

	DeleteFile(ctx context.Context, tag domain.FileTag, fileID string) error
}

type ClientService interface {
	Get(context.Context, uuid.UUID) (domain.Client, error)
}

type UsersRepository interface {
	GetUser(ctx context.Context, user domain.User) (domain.User, error)
}

type RequestsRepository interface {
	Create(ctx context.Context, request domain.Request) (uuid.UUID, error)

	Get(ctx context.Context, id uuid.UUID) (domain.Request, error)
	GetMany(ctx context.Context, page, limit int, filter domain.RequestQueryFilter) ([]domain.Request, int64, error)
	GetStats(ctx context.Context) (domain.RequestStats, error)

	Patch(ctx context.Context, id uuid.UUID, request domain.Request) error
	PatchStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.RequestStatus) error
}

func NewService(
	mediaService MediaService,
	cache core_redis_cache.Cache,
	clientService ClientService,
	usersPostgresRepo UsersRepository,
	requestsMongoRepo RequestsRepository,
) *Service {
	return &Service{
		mediaService:      mediaService,
		cache:             cache,
		clientService:     clientService,
		usersPostgresRepo: usersPostgresRepo,
		requestsMongoRepo: requestsMongoRepo,
	}
}
