package core_media_utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"

	"github.com/disintegration/gift"
	"golang.org/x/image/draw"
)

func ResizeAnimatedGIF(gifImg *gif.GIF, targetW, targetH int) (*gif.GIF, error) {
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

	return out, nil
}

func ResizeImage(img image.Image, format Format, targetW, targetH int) (image.Image, error) {
	dstRect := image.Rect(0, 0, targetW, targetH)

	dst := image.NewRGBA(dstRect)

	draw.CatmullRom.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	return dst, nil
}

func ResizeAny(data []byte, targetW, targetH int, quality *int) ([]byte, error) {
	format, err := DetectFormat(data)
	if err != nil {
		return nil, fmt.Errorf("detect format: %w", err)
	}

	if format == FormatGIF {
		g, err := gif.DecodeAll(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("decode gif: %v: %w", err, media_errors.ErrInvalidImage)
		}

		result, err := ResizeAnimatedGIF(g, targetW, targetH)
		if err != nil {
			return nil, fmt.Errorf("resize animated gif: %w", err)
		}

		var out bytes.Buffer
		if err := gif.EncodeAll(&out, result); err != nil {
			return nil, fmt.Errorf("encode gif: %w", err)
		}

		return out.Bytes(), nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %v: %w", err, media_errors.ErrInvalidImage)
	}

	dst, err := ResizeImage(img, format, targetW, targetH)
	if err != nil {
		return nil, fmt.Errorf("resize image: %w", err)
	}

	if quality == nil {
		q := 75
		quality = &q
	}

	var out bytes.Buffer
	switch format {
	case FormatJPEG:
		err = jpeg.Encode(&out, dst, &jpeg.Options{Quality: *quality})
	case FormatPNG:
		err = png.Encode(&out, dst)
	case FormatWebP:
		err = encodeWebP(&out, dst, *quality)
	default:
		err = fmt.Errorf("unsupported image format '%s': %w", format, media_errors.ErrUnsupportedFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", format, err)
	}

	return out.Bytes(), nil
}
