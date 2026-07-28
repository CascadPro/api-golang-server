package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

type FileTag string

const (
	FileTagAvatars = FileTag("avatars")
	FileTagImages  = FileTag("images")
	FileTagVideos  = FileTag("videos")
	FileTagDocs    = FileTag("docs")
	FileTagNil     = FileTag("")
)

var (
	FileTags = []FileTag{FileTagAvatars, FileTagDocs, FileTagImages, FileTagVideos}
)

type File struct {
	ID      string
	Version int

	Tag         FileTag
	Filename    string
	ContentType string
	Size        int

	Deleted   bool
	DeletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFile() File {
	return File{}
}

func (f *File) Validate() error {
	if f.Tag == FileTagNil {
		return fmt.Errorf("`Tag` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := core_validation.ValidateStringLength(&f.Filename, "Filename", 1, 255); err != nil {
		return err
	}

	if f.Size <= 0 {
		return fmt.Errorf("`Size` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}

	if f.Deleted {
		if f.DeletedAt == nil {
			return fmt.Errorf("`DeletedAt` can't be NULL if `Deleted` is true: %w", core_errors.ErrInvalidArgument)
		}

		if f.CreatedAt.After(*f.DeletedAt) {
			return fmt.Errorf("`DeletedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if f.DeletedAt != nil {
			return fmt.Errorf("`DeletedAt` must be NULL if `Deleted` is false: %w", core_errors.ErrInvalidArgument)
		}
	}

	if f.CreatedAt.After(f.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
