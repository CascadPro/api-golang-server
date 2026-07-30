package media_postgres_repository

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

type FileModel struct {
	ID      string
	Version int

	Tag      domain.FileTag
	Filename string
	MimeType domain.FileMimeType
	Size     int

	Deleted   bool
	DeletedAt *time.Time
	CreatedAt time.Time
}

func domainFromModel(model FileModel) domain.File {
	return domain.File{
		ID:        model.ID,
		Version:   model.Version,
		Tag:       model.Tag,
		Filename:  model.Filename,
		MimeType:  model.MimeType,
		Size:      model.Size,
		Deleted:   model.Deleted,
		DeletedAt: model.DeletedAt,
		CreatedAt: model.CreatedAt,
	}
}
