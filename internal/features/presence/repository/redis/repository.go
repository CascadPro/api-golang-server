package presence_redis_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"

	"github.com/google/uuid"
)

const (
	LeaseTTL          = 60 * time.Second
	HeartbeatInterval = 15 * time.Second
)

type Repository struct {
	pool core_redis_pool.Pool
}

type RepositoryMethods interface {
	Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (RegisterResult, error)
	Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error)

	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)

	Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error
}

func NewRepository(pool core_redis_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func presenceKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s:users:%s", core_redis_pool.PresenceFolder, userID)
}

func presenceMember(sessionID, connectionID string) string {
	return fmt.Sprintf("%s:members:%s:%s", core_redis_pool.PresenceFolder, sessionID, connectionID)
}

// cascade__presence:members:1f64181baebd284cb5786c9c:bf6fb6c0-8bcc-468e-849e-db9158572fcc

func parsePresenceMember(member string) (sessionID, connectionID string, ok bool) {
	folder := core_redis_pool.PresenceFolder + ":"

	if len(member) <= len(folder) {
		return "", "", false
	}

	if member[:len(folder)] != string(folder) {
		return "", "", false
	}

	value := member[len(folder):]

	sessionID, connectionID, ok = splitMember(value)

	if !ok {
		return "", "", false
	}

	return sessionID, connectionID, true
}

func splitMember(value string) (sessionID, connectionID string, ok bool) {
	value = strings.TrimPrefix(value, "members:")

	for i := 0; i < len(value); i++ {
		if value[i] != ':' {
			continue
		}

		sessionID = value[:i]
		connectionID = value[i+1:]

		if sessionID == "" || connectionID == "" {
			return "", "", false
		}

		return sessionID, connectionID, true
	}

	return "", "", false
}
