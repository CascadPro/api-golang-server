package core_avatar_utils

import (
	"fmt"
	"io"

	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func ProcessAvatar(r io.Reader) (*ProcessedAvatar, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read avatar: %w", err)
	}

	if len(data) == 0 {
		return nil, media_errors.ErrInvalidImage
	}

	format, err := core_media_utils.DetectFormat(data)
	if err != nil {
		return nil, err
	}

	switch format {
	case core_media_utils.FormatGIF:
		return processGIF(data)

	case core_media_utils.FormatJPEG, core_media_utils.FormatPNG, core_media_utils.FormatWebP:
		return processStatic(data)

	default:
		return nil, fmt.Errorf("%w: %s", media_errors.ErrUnsupportedFormat, format)
	}
}
