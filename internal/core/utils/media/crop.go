package core_media_utils

import (
	"image"

	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"golang.org/x/image/draw"
)

var cropAnalyzer = smartcrop.NewAnalyzer(
	nfnt.NewDefaultResizer(),
)

func CropImage(img image.Image, crop image.Rectangle) image.Image {
	cropped := image.NewRGBA(
		image.Rect(0, 0, crop.Dx(), crop.Dy()),
	)

	draw.CatmullRom.Scale(cropped, cropped.Bounds(), img, crop, draw.Over, nil)

	return cropped
}

func FindCrop(img image.Image, width, height int) (image.Rectangle, error) {
	bounds := img.Bounds()

	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return image.Rectangle{}, media_errors.ErrInvalidImage
	}

	result, err := cropAnalyzer.FindBestCrop(
		img,
		width,
		height,
	)

	if err != nil {
		return image.Rectangle{}, err
	}

	return result, nil
}
