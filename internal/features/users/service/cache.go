package users_service

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

const (
	cacheKeyUUID = "uuid"
)

var (
	cacheTTLUser = 5 * time.Minute
)

type cacheUser struct {
	User        domain.User `json:"user"`
	Placeholder []byte      `json:"placeholder"`
}

func cacheKeyUser(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", cacheKeyUUID, id.String())
}
