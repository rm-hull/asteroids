package fonts

import (
	"log"

	"github.com/rm-hull/asteroids/resources/fonts"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	AsteroidsDisplayFont16 font.Face
	AsteroidsDisplayFont32 font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.AsteroidsDisplay_ttf)
	if err != nil {
		log.Fatal(err)
	}

	AsteroidsDisplayFont16, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     144,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	AsteroidsDisplayFont32, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     144,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
