package fonts

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rm-hull/asteroids/resources/fonts"
)

var (
	AsteroidsFace32 *text.GoTextFace
	AsteroidsFace64 *text.GoTextFace
)

func init() {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.AsteroidsDisplay_ttf))
	if err != nil {
		log.Fatal(err)
	}

	AsteroidsFace32 = &text.GoTextFace{
		Source: source,
		Size:   32,
	}

	AsteroidsFace64 = &text.GoTextFace{
		Source: source,
		Size:   64,
	}
}
