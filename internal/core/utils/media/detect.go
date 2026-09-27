package core_media_utils

import (
	"bytes"
	"fmt"
	"image"

	// Import for GIF decoder registration
	_ "image/gif"

	// Import for JPEG decoder registration
	_ "image/jpeg"

	// Import for PNG decoder registration
	_ "image/png"

	// Import for Image draw registration
	_ "golang.org/x/image/draw"

	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func DetectFormat(data []byte) (Format, error) {
	_, f, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("decode config: %v: %w", err, media_errors.ErrInvalidImage)
	}

	format := Format(f)

	switch format {
	case FormatGIF, FormatJPEG, FormatPNG, FormatWebP:
		return format, nil

	default:
		return "", media_errors.ErrUnsupportedFormat
	}
}
