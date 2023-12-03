package text_align

import (
	"github.com/rm-hull/asteroids/internal/geometry"
	"golang.org/x/image/font"
)

func Center(screenBounds *geometry.Dimension, text string, face font.Face) (int, int) {
	bounds, _ := font.BoundString(face, text)
	dx := float64(bounds.Max.X.Round() - bounds.Min.X.Round())
	dy := float64(bounds.Max.Y.Round() - bounds.Min.Y.Round())

	return int(screenBounds.W-dx) / 2, int(screenBounds.H-dy) / 2
}
