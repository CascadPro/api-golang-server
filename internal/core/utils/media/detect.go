package core_media_utils

import (
	"bytes"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	_ "golang.org/x/image/draw"
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
