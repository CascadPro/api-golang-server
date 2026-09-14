package core_avatar_utils

import (
	"image"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
)

func processImage(img image.Image, format core_media_utils.Format, crop image.Rectangle) image.Image {
	img = core_media_utils.CropImage(img, crop)

	resized, _ := core_media_utils.ResizeImage(img, format, domain.FileAvatarS3Size, domain.FileAvatarS3Size)

	return resized
}

func processGifFrame(frame image.Image, crop image.Rectangle) image.Image {
	frame = core_media_utils.CropImage(frame, crop)

	resized, _ := core_media_utils.ResizeImage(frame, core_media_utils.FormatPNG, domain.FileAvatarS3Size, domain.FileAvatarS3Size)

	return resized
}
