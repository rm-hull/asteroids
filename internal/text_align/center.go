package text_align

import (

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rm-hull/asteroids/internal/geometry"
)

func Center(screenBounds *geometry.Dimension, message string, face text.Face) (int, int) {
	width, height := text.Measure(message, face, 0)
	return int(screenBounds.W-width) / 2, int(screenBounds.H-height) / 2
}
