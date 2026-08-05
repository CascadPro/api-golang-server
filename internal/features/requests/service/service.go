package requests_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	requests_mongo_repository "github.com/CascadePro/api-golang-server/internal/features/requests/repository/mongo"
	"github.com/google/uuid"
)

type Service struct {
	mediaService      media_service.ServiceMethods
	requestsMongoRepo requests_mongo_repository.RepositoryMethods
}

type ServiceMethods interface {
}

func NewService(
	mediaService media_service.ServiceMethods,
	requestsMongoRepo requests_mongo_repository.RepositoryMethods,
) *Service {
	return &Service{
		mediaService:      mediaService,
		requestsMongoRepo: requestsMongoRepo,
	}
}
