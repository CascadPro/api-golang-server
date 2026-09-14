package core_avatar_utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func processGIF(data []byte) (*ProcessedAvatar, error) {
	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode gif: %v: %w", err, media_errors.ErrInvalidImage)
	}

	if len(g.Image) == 0 {
		return nil, media_errors.ErrInvalidImage
	}

	frames := renderGIFFrames(g)

	if len(frames) == 0 {
		return nil, media_errors.ErrInvalidImage
	}

	crop, err := core_media_utils.FindCrop(frames[0], domain.FileAvatarS3Size, domain.FileAvatarS3Size)
	if err != nil {
		return nil, fmt.Errorf("find gif crop: %w", err)
	}

	processed := make([]*image.Paletted, len(frames))

	for i, frame := range frames {
		processed[i] = toPaletted(processGifFrame(frame, crop))
	}

	outGIF := &gif.GIF{
		Image:           processed,
		Delay:           append([]int(nil), g.Delay...),
		Disposal:        append([]byte(nil), g.Disposal...),
		LoopCount:       g.LoopCount,
		BackgroundIndex: g.BackgroundIndex,
		Config: image.Config{
			Width:      domain.FileAvatarS3Size,
			Height:     domain.FileAvatarS3Size,
			ColorModel: g.Config.ColorModel,
		},
	}

	var out bytes.Buffer

	if err := gif.EncodeAll(&out, outGIF); err != nil {
		return nil, fmt.Errorf("encode gif: %w", err)
	}

	return &ProcessedAvatar{
		Data:      out.Bytes(),
		Format:    core_media_utils.FormatGIF,
		MimeType:  domain.FileMimeTypeGif,
		Extension: ".gif",
	}, nil
}
