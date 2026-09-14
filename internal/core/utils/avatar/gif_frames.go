package core_avatar_utils

import (
	"image"
	"image/draw"
	"image/gif"
)

func renderGIFFrames(g *gif.GIF) []image.Image {
	if len(g.Image) == 0 {
		return nil
	}

	canvas := image.NewRGBA(
		image.Rect(0, 0, g.Config.Width, g.Config.Height),
	)

	frames := make([]image.Image, 0, len(g.Image))

	for i, frame := range g.Image {
		bounds := frame.Bounds()

		draw.Draw(
			canvas,
			bounds,
			frame,
			bounds.Min,
			draw.Over,
		)

		copyFrame := image.NewRGBA(
			canvas.Bounds(),
		)

		draw.Draw(
			copyFrame,
			copyFrame.Bounds(),
			canvas,
			image.Point{},
			draw.Src,
		)

		frames = append(frames, copyFrame)

		if i >= len(g.Disposal) {
			continue
		}

		switch g.Disposal[i] {
		case gif.DisposalBackground:
			draw.Draw(
				canvas,
				bounds,
				image.Transparent,
				image.Point{},
				draw.Src,
			)
		}
	}

	return frames
}
