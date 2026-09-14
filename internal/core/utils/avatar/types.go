package core_avatar_utils

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
)

type ProcessedAvatar struct {
	Data      []byte
	Format    core_media_utils.Format
	MimeType  domain.FileMimeType
	Extension string
}
