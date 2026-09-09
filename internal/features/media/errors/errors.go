package media_errors

import "errors"

var (
	ErrFileNotFound            = errors.New("file not found")
	ErrFilePlaceholderNotFound = errors.New("file placeholder empty")
	ErrFileAlreadyExists       = errors.New("file already exists")
	ErrFileAccessDenied        = errors.New("file access denied")
	ErrFileDeleted             = errors.New("file was deleted")
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrFileTooLarge            = errors.New("file too large")
)
