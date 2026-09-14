package core_media_utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"

	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"

	"golang.org/x/image/draw"
)

// GeneratePlaceholder – берёт оригинальный байтовый массив изображения,
// уменьшает его до 64×{Y} px, размывает и кодирует в JPEG с низким качеством.
func GeneratePlaceholder(original []byte) ([]byte, error) {
	srcImg, _, err := image.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, fmt.Errorf("decode image: %v: %w", err, media_errors.ErrInvalidImage)
	}

	origBounds := srcImg.Bounds()
	origW, origH := origBounds.Dx(), origBounds.Dy()
	if origW == 0 || origH == 0 {
		return nil, fmt.Errorf("zero dimension image: %w", media_errors.ErrInvalidImage)
	}

	scaleW := float64(placeholderResizeMax) / float64(origW)
	scaleH := float64(placeholderResizeMax) / float64(origH)
	scale := math.Min(scaleW, scaleH)
	if scale > 1 {
		scale = 1
	}

	targetW := int(math.Round(float64(origW) * scale))
	targetH := int(math.Round(float64(origH) * scale))
	if targetW < 1 {
		targetW = 1
	}
	if targetH < 1 {
		targetH = 1
	}

	thumb := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	draw.ApproxBiLinear.Scale(thumb, thumb.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	blurred := thumb
	for range placeholderBlurIter {
		blurred = simpleBoxBlur(blurred)
	}

	var out bytes.Buffer
	if err := jpeg.Encode(&out, blurred, &jpeg.Options{Quality: placeholderQuality}); err != nil {
		return nil, fmt.Errorf("jpeg encode: %w", err)
	}

	return out.Bytes(), nil
}
