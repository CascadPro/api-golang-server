package requests_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	requests_mongo_repository "github.com/CascadePro/api-golang-server/internal/features/requests/repository/mongo"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	mediaService      media_service.ServiceMethods
	usersPostgresRepo users_postgres_repository.RepositoryMethods
	requestsMongoRepo requests_mongo_repository.RepositoryMethods
}

type ServiceMethods interface {
	GetRequest(context.Context, uuid.UUID) (domain.Request, domain.User, map[int]domain.File, error)
	GetRequests(context.Context, int, int, domain.RequestQueryFilter) ([]domain.Request, int64, error)
	PatchRequest(context.Context, uuid.UUID, domain.RequestPatch) (domain.Request, error)
	PatchRequestStatus(context.Context, uuid.UUID, domain.RequestStatus) error
	UploadDoc(context.Context, uuid.UUID, domain.File, []byte, int) error
}

func NewService(
	mediaService media_service.ServiceMethods,
	usersPostgresRepo users_postgres_repository.RepositoryMethods,
	requestsMongoRepo requests_mongo_repository.RepositoryMethods,
) *Service {
	return &Service{
		mediaService:      mediaService,
		usersPostgresRepo: usersPostgresRepo,
		requestsMongoRepo: requestsMongoRepo,
	}
}
