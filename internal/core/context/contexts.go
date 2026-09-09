package core_context

import (
	"context"
	"net"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID.String())
}

func WithUserRole(ctx context.Context, role domain.UserRole) context.Context {
	return context.WithValue(ctx, ctxKeyUserRole, role)
}

func WithSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ctxKeySessionID, sid)
}

func WithFileMimeType(ctx context.Context, mimeType domain.FileMimeType) context.Context {
	return context.WithValue(ctx, ctxKeyMimeType, mimeType)
}

func WithFileTag(ctx context.Context, tag domain.FileTag) context.Context {
	return context.WithValue(ctx, ctxKeyTag, tag)
}

func WithClientIP(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, ctxKeyIP, ip)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, requestID)
}

func WithLocale(ctx context.Context, locale language.Tag) context.Context {
	return context.WithValue(ctx, ctxKeyLocale, locale)
}
