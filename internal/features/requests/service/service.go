package requests_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	client_service "github.com/CascadePro/api-golang-server/internal/features/client/service"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	"github.com/google/uuid"
)

type Service struct {
	mediaService      media_service.ServiceMethods
	cache             core_redis_cache.Cache
	clientService     client_service.ServiceMethods
	usersPostgresRepo UsersRepository
	requestsMongoRepo RequestsRepository
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

type ServiceMethods interface {
	CreateRequest(context.Context, domain.Request) (domain.Request, error)

	GetRequest(context.Context, uuid.UUID) (domain.Request, domain.User, domain.Client, map[int]domain.File, error)
	GetRequests(context.Context, int, int, domain.RequestQueryFilter) ([]domain.Request, int64, error)
	GetRequestsStats(context.Context) (domain.RequestStats, error)

	PatchRequest(context.Context, uuid.UUID, domain.RequestPatch) (domain.Request, error)
	PatchRequestStatus(context.Context, uuid.UUID, domain.RequestStatus) error

	UploadDoc(context.Context, uuid.UUID, domain.File, []byte, int) error
	DeleteDoc(context.Context, uuid.UUID, int) error
}

func NewService(
	mediaService media_service.ServiceMethods,
	cache core_redis_cache.Cache,
	clientService client_service.ServiceMethods,
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
