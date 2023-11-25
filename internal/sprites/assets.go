package sprites

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/asteroids-2x.png
var assets embed.FS

var spriteSheet = mustLoadImage("assets/asteroids-2x.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func sprite(spriteSheet *ebiten.Image, x, y, w, h int) *ebiten.Image {
	img := spriteSheet.SubImage(image.Rectangle{
		image.Point{x, y},
		image.Point{x + w, y + h},
	})

	return ebiten.NewImageFromImage(img)
}

// 2x
var BigAsteroid1 = sprite(spriteSheet, 0, 0, 160, 160)
var BigAsteroid2 = sprite(spriteSheet, 160, 0, 160, 160)
var BigAsteroid3 = sprite(spriteSheet, 320, 0, 160, 160)

var MediumAsteroid1 = sprite(spriteSheet, 0, 160, 96, 96)
var MediumAsteroid2 = sprite(spriteSheet, 96, 160, 96, 96)
var MediumAsteroid3 = sprite(spriteSheet, 192, 160, 96, 96)

var SmallAsteroid1 = sprite(spriteSheet, 0, 254, 64, 64)
var SmallAsteroid2 = sprite(spriteSheet, 64, 254, 64, 64)
var SmallAsteroid3 = sprite(spriteSheet, 128, 254, 64, 64)

var AlienSpaceShip = sprite(spriteSheet, 416, 160, 96, 80)

var SpaceShip1 = sprite(spriteSheet, 192, 254, 96, 64)
var SpaceShip2 = sprite(spriteSheet, 288, 254, 96, 64)

var Bullet1 = sprite(spriteSheet, 448, 286, 32, 32)
var Bullet2 = sprite(spriteSheet, 480, 286, 32, 32)
