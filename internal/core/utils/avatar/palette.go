package core_avatar_utils

import (
	"image"
	"image/color"
	"sort"

	"golang.org/x/image/draw"
)

type colorBucket struct {
	R uint8
	G uint8
	B uint8
	N uint64
}

func toPaletted(img image.Image) *image.Paletted {
	palette := quantizePalette(img, 256)

	dst := image.NewPaletted(
		img.Bounds(),
		palette,
	)

	draw.FloydSteinberg.Draw(
		dst,
		dst.Bounds(),
		img,
		image.Point{},
	)

	return dst
}

func quantizePalette(
	img image.Image,
	size int,
) color.Palette {
	const buckets = 16 * 16 * 16

	histogram := make([]uint64, buckets)

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if a < 0x8000 {
				continue
			}

			ri := uint8(r >> 12)
			gi := uint8(g >> 12)
			bi := uint8(b >> 12)

			index :=
				int(ri)*256 +
					int(gi)*16 +
					int(bi)

			histogram[index]++
		}
	}

	result := make(
		[]colorBucket,
		0,
		buckets,
	)

	for i, count := range histogram {
		if count == 0 {
			continue
		}

		result = append(result, colorBucket{
			R: uint8((i / 256) * 17),
			G: uint8(((i / 16) % 16) * 17),
			B: uint8((i % 16) * 17),
			N: count,
		})
	}

	sort.Slice(
		result,
		func(i, j int) bool {
			return result[i].N > result[j].N
		},
	)

	if len(result) > size {
		result = result[:size]
	}

	palette := make(
		color.Palette,
		0,
		len(result),
	)

	for _, c := range result {
		palette = append(
			palette,
			color.RGBA{
				R: c.R,
				G: c.G,
				B: c.B,
				A: 255,
			},
		)
	}

	if len(palette) == 0 {
		return color.Palette{
			color.Black,
			color.White,
		}
	}

	return palette
}
