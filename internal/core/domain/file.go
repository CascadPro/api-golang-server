package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

type FileTag string
type FileMimeType string

const (
	FileTagAvatars = FileTag("avatars")
	FileTagImages  = FileTag("images")
	FileTagVideos  = FileTag("videos")
	FileTagDocs    = FileTag("docs")
	FileTagNil     = FileTag("")
)

const (
	FileMimeTypeJpeg = FileMimeType("image/jpeg")
	FileMimeTypePng  = FileMimeType("image/png")
	FileMimeTypeGif  = FileMimeType("image/gif")
	FileMimeTypeWebp = FileMimeType("image/webp")
	FileMimeTypeMMp4 = FileMimeType("video/mp4")
	FileMimeTypeJson = FileMimeType("application/json")
	FileMimeTypePdf  = FileMimeType("application/pdf")
	FileMimeTypeDocx = FileMimeType("application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	FileMimeTypeXlsx = FileMimeType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	FileMimeTypePptx = FileMimeType("application/vnd.openxmlformats-officedocument.presentationml.presentation")
	FileMimeTypeNil  = FileMimeType("")
)

var (
	FileIDByteLength    int = 24
	FileAvatarS3Size    int = 512
	FileAvatarS3Quality int = 100
)

var (
	FileTags = []FileTag{FileTagAvatars, FileTagDocs, FileTagImages, FileTagVideos}
)

type File struct {
	ID      string
	Version int64

	Tag      FileTag
	Filename string
	MimeType FileMimeType
	Size     int64

	placeholder []byte

	Deleted   bool
	DeletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUploadFile(filename string, tag FileTag, mimeType FileMimeType, size int64) File {
	return File{
		ID:       UninitializedID,
		Version:  UninitializedVersion,
		Tag:      tag,
		Filename: filename,
		MimeType: mimeType,
		Size:     size,
	}
}

func (f *File) GeneratePlaceholder(content []byte) error {
	if f.Tag != FileTagAvatars && f.Tag != FileTagImages && f.Tag != FileTagVideos {
		return fmt.Errorf("generate placeholder: %w", core_errors.ErrInvalidArgument)
	}

	placeholder, err := core_media_utils.GeneratePlaceholder(content)
	if err != nil {
		return fmt.Errorf("generate placeholder: %w", err)
	}

	f.placeholder = placeholder
	return nil
}

func (f *File) GetPlaceholder() []byte {
	return f.placeholder
}

func (f *File) Validate() error {
	if f.Tag == FileTagNil {
		return fmt.Errorf("`Tag` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if f.MimeType == FileMimeTypeNil {
		return fmt.Errorf("`MimeType` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := core_validation.ValidateStringLength(&f.Filename, "Filename", 1, 255); err != nil {
		return err
	}

	if f.Size <= 0 {
		return fmt.Errorf("`Size` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}

	if f.Tag == FileTagAvatars || f.Tag == FileTagImages || f.Tag == FileTagVideos {
		if f.placeholder == nil {
			return fmt.Errorf(
				"`Placeholder` can't be NULL for images, avatars and videos: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if err := f.ValidateDeleted(); err != nil {
		return err
	}

	if f.CreatedAt.After(f.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (f *File) ValidateDeleted() error {
	if f.Deleted {
		if f.DeletedAt == nil {
			return fmt.Errorf("`DeletedAt` can't be NULL if `Deleted` is true: %w", core_errors.ErrInvalidArgument)
		}

		if f.CreatedAt.After(*f.DeletedAt) {
			return fmt.Errorf("`DeletedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
		}

		if (*f.DeletedAt).After(time.Now()) {
			return fmt.Errorf("`DeletedAt` can't be after current time: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if f.DeletedAt != nil {
			return fmt.Errorf("`DeletedAt` must be NULL if `Deleted` is false: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
