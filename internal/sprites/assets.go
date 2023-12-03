package sprites

import (
	"bytes"
	"image"
	_ "image/png"
	"math/rand"

	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/resources/images"

	"github.com/hajimehoshi/ebiten/v2"
)

var spriteSheet = mustLoadImage(images.Asteroids2X_png)

func mustLoadImage(b []byte) *ebiten.Image {
	r := bytes.NewReader(b)
	img, _, err := image.Decode(r)
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
var LargeAsteroids = []*ebiten.Image{
	sprite(spriteSheet, 0, 0, 160, 160),
	sprite(spriteSheet, 160, 0, 160, 160),
	sprite(spriteSheet, 320, 0, 160, 160),
}

var MediumAsteroids = []*ebiten.Image{
	sprite(spriteSheet, 0, 160, 80, 96),
	sprite(spriteSheet, 96, 160, 80, 96),
	sprite(spriteSheet, 192, 160, 80, 96),
}

var SmallAsteroids = []*ebiten.Image{
	sprite(spriteSheet, 0, 254, 64, 64),
	sprite(spriteSheet, 64, 254, 64, 64),
	sprite(spriteSheet, 128, 254, 64, 64),
}

var AlienSpaceShip = sprite(spriteSheet, 416, 160, 96, 80)

var SpaceShip1 = sprite(spriteSheet, 192, 254, 96, 64)
var SpaceShip2 = sprite(spriteSheet, 288, 254, 96, 64)

var Bullet1 = sprite(spriteSheet, 448, 286, 32, 32)
var Bullet2 = sprite(spriteSheet, 480, 286, 32, 32)

const (
	Large = iota
	Medium
	Small
)

func Asteroid(size int) *ebiten.Image {
	idx := rand.Intn(3)
	switch size {
	case Large:
		return LargeAsteroids[idx]
	case Medium:
		return MediumAsteroids[idx]
	default:
		return SmallAsteroids[idx]
	}
}

func Bullet(size int) *ebiten.Image {
	switch size {
	case Large:
		return Bullet2
	default:
		return Bullet1
	}
}

func Centre(sprite *ebiten.Image) geometry.Vector {
	bounds := sprite.Bounds()
	return geometry.Vector{
		X: float64(bounds.Dx()) / 2,
		Y: float64(bounds.Dy()) / 2,
	}
}

func Size(sprite *ebiten.Image) *geometry.Dimension {
	bounds := sprite.Bounds()
	return &geometry.Dimension{
		W: float64(bounds.Dx()),
		H: float64(bounds.Dy()),
	}
}
