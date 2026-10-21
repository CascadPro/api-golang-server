package requests_service

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

const (
	cacheKeyStats = "stats"
)

var (
	cacheTTLStats   = 1 * time.Hour
	cacheTTLRequest = 10 * time.Minute
)

type cacheRequest struct {
	Request domain.Request      `json:"request"`
	User    domain.User         `json:"status_by"`
	Client  domain.Client       `json:"client"`
	Files   map[int]domain.File `json:"files"`
}

func cacheKeyRequest(id uuid.UUID) string {
	return fmt.Sprintf("uuid:%s", id.String())
}
