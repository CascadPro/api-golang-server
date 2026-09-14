package core_avatar_utils

import (
	"bytes"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	"github.com/chai2010/webp"
)

func processStatic(data []byte) (*ProcessedAvatar, error) {
	img, format, err := core_media_utils.DecodeImage(data)
	if err != nil {
		return nil, err
	}

	switch format {
	case core_media_utils.FormatJPEG, core_media_utils.FormatPNG, core_media_utils.FormatWebP:
	default:
		return nil, fmt.Errorf("%w: %s", media_errors.ErrUnsupportedFormat, format)
	}

	crop, err := core_media_utils.FindCrop(img, domain.FileAvatarS3Size, domain.FileAvatarS3Size)
	if err != nil {
		return nil, fmt.Errorf("find crop: %w", err)
	}

	img = processImage(img, format, crop)

	opts := &webp.Options{
		Quality: float32(domain.FileAvatarS3Quality),
	}

	var out bytes.Buffer

	if err := webp.Encode(&out, img, opts); err != nil {
		return nil, fmt.Errorf("encode webp: %w", err)
	}

	return &ProcessedAvatar{
		Data:      out.Bytes(),
		Format:    core_media_utils.FormatWebP,
		MimeType:  domain.FileMimeTypeWebp,
		Extension: ".webp",
	}, nil
}
