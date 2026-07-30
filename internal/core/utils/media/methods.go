package core_media_utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"

	"github.com/disintegration/gift"
	"golang.org/x/image/draw"
)

// src – полные байты GIF‑анимации.
// targetW / targetH – желаемые размеры.
// Возвращает GIF той же анимации, но уже с нужными сторонами.
func ResizeAnimatedGIF(src []byte, targetW, targetH int) ([]byte, error) {
	gifImg, err := gif.DecodeAll(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("gif decode: %w", err)
	}

	out := &gif.GIF{
		BackgroundIndex: gifImg.BackgroundIndex,
		LoopCount:       gifImg.LoopCount,
	}

	filter := gift.New(gift.Resize(targetW, targetH, gift.LanczosResampling))

	for i, srcPaletted := range gifImg.Image {
		bounds := srcPaletted.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, srcPaletted, bounds.Min, draw.Src)

		dstRGBA := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
		filter.Draw(dstRGBA, rgba)

		dstPaletted := image.NewPaletted(dstRGBA.Bounds(), srcPaletted.Palette)

		draw.FloydSteinberg.Draw(dstPaletted, dstRGBA.Bounds(), dstRGBA, image.Point{})

		out.Image = append(out.Image, dstPaletted)
		out.Delay = append(out.Delay, gifImg.Delay[i])
		if len(gifImg.Disposal) > i {
			out.Disposal = append(out.Disposal, gifImg.Disposal[i])
		} else {
			out.Disposal = append(out.Disposal, gif.DisposalNone)
		}
	}
	out.Config.Width = targetW
	out.Config.Height = targetH

	var buf bytes.Buffer
	if err = gif.EncodeAll(&buf, out); err != nil {
		return nil, fmt.Errorf("gif encode: %w", err)
	}
	return buf.Bytes(), nil
}

// Src – любые растровые данные (jpeg, png, webp, static gif, bmp, …).
// TargetW / TargetH – желаемый размер.
// Возвращает изображение того же MIME‑типа, что и входное.
func ResizeImage(src []byte, targetW, targetH, quality int) ([]byte, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("detect format: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	dstRect := image.Rect(0, 0, targetW, targetH)
	dst := image.NewRGBA(dstRect)
	draw.CatmullRom.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	var out bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&out, dst, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&out, dst)
	case "webp":
		err = encodeWebP(&out, dst, quality)
	case "gif":
		var paletted *image.Paletted
		if p, ok := img.(*image.Paletted); ok {
			paletted = image.NewPaletted(dstRect, p.Palette)
		} else {
			paletted = image.NewPaletted(dstRect, nil)
		}

		draw.FloydSteinberg.Draw(paletted, dstRect, dst, image.Point{})

		g := &gif.GIF{
			Image: []*image.Paletted{paletted},
			Delay: []int{0},
		}
		g.Config.Width = targetW
		g.Config.Height = targetH

		err = gif.EncodeAll(&out, g)
	default:
		err = fmt.Errorf("unsupported image format '%s': %w", format, core_errors.ErrInvalidArgument)
	}
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", format, err)
	}
	return out.Bytes(), nil
}

func ResizeAny(src []byte, targetW, targetH int, quality *int) ([]byte, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("detect format: %w", err)
	}

	if format == "gif" {
		if g, err := gif.DecodeAll(bytes.NewReader(src)); err == nil && len(g.Image) > 1 {
			return ResizeAnimatedGIF(src, targetW, targetH)
		}
	}

	if quality == nil {
		tmp := 90
		quality = &tmp
	}

	return ResizeImage(src, targetW, targetH, *quality)
}

// GeneratePlaceholder – берёт оригинальный байтовый массив изображения,
// уменьшает его до 64×{Y} px, размывает и кодирует в JPEG с низким качеством.
func GeneratePlaceholder(original []byte) ([]byte, error) {
	srcImg, _, err := image.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	origBounds := srcImg.Bounds()
	origW, origH := origBounds.Dx(), origBounds.Dy()
	if origW == 0 || origH == 0 {
		return nil, fmt.Errorf("zero dimension image: %w", core_errors.ErrInvalidArgument)
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
