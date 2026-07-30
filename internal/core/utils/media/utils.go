package core_media_utils

import (
	"image"
	"image/color"
	"io"

	webp "github.com/urtie/gowebp"
)

const (
	placeholderResizeMax = 64
	placeholderBlurIter  = 256
	placeholderQuality   = 20
)

// Mode    – задаёт тип сжатия (lossy / lossless). Мы используем lossy.
// Quality – качество 0‑100 (тип LossyQuality).
func encodeWebP(w io.Writer, img image.Image, quality int) error {
	opts := &webp.EncodeOptions{
		Mode:    webp.EncodeLossy,
		Quality: webp.LossyQuality(quality),
	}
	return webp.Encode(w, img, opts)
}

// simpleBoxBlur делает один проход «box‑blur» (среднее значение 3×3).
// Это полностью реализовано на чистом Go, без сторонних библиотек.
func simpleBoxBlur(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	// Пройдем по каждому пикселю, посчитаем среднее значение соседей.
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b2, a, cnt := sumNeighbors(src, x, y, b)

			// Усредняем
			r /= cnt
			g /= cnt
			b2 /= cnt
			a /= cnt

			dst.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b2 >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	return dst
}

// sumNeighbors sums RGBA values in a 3x3 neighborhood around (x,y) within bounds b.
func sumNeighbors(src image.Image, x, y int, b image.Rectangle) (r, g, b2, a, cnt uint32) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px := x + dx
			py := y + dy
			if px < b.Min.X || px >= b.Max.X || py < b.Min.Y || py >= b.Max.Y {
				continue
			}
			rr, gg, bb, aa := src.At(px, py).RGBA()
			r += rr
			g += gg
			b2 += bb
			a += aa
			cnt++
		}
	}

	if cnt == 0 {
		cnt = 1
	}

	return r, g, b2, a, cnt
}
