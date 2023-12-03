package sprites

import (
	"github.com/rm-hull/asteroids/internal/geometry"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Sprite struct {
	Position     *geometry.Vector
	Velocity     *geometry.Vector
	Orientation  float64
	Direction    float64
	Speed        float64
	Rotation     float64
	Centre       *geometry.Vector
	Size         *geometry.Dimension
	Image        *ebiten.Image
	screenBounds *geometry.Dimension
	wraparound   bool
	ColorModel   colorm.ColorM
	DrawOptions  *colorm.DrawImageOptions
}

func NewSprite(screenBounds *geometry.Dimension, image *ebiten.Image, wraparound bool) *Sprite {
	bounds := image.Bounds()
	centre := geometry.Vector{
		X: float64(bounds.Dx()) / 2,
		Y: float64(bounds.Dy()) / 2,
	}

	return &Sprite{
		Position:     geometry.Zero(),
		Velocity:     geometry.Zero(),
		Orientation:  0,
		Direction:    0,
		Speed:        0,
		Rotation:     0,
		Centre:       &centre,
		Size:         Size(image),
		Image:        image,
		screenBounds: screenBounds,
		wraparound:   wraparound,
		ColorModel:   colorm.ColorM{},
		DrawOptions:  &colorm.DrawImageOptions{},
	}
}

func (s *Sprite) Reset() {
	s.Orientation = 0
	s.Rotation = 0
	s.Direction = 0
	s.Speed = 0
	s.Velocity.X = 0
	s.Velocity.Y = 0
}

func (s *Sprite) Update() error {
	s.Orientation += s.Rotation
	s.Position.Add(s.Velocity)
	if s.wraparound {
		s.checkEdges()
	}
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) {

	s.DrawOptions.GeoM.Translate(-s.Centre.X, -s.Centre.Y)
	s.DrawOptions.GeoM.Rotate(s.Orientation)
	s.DrawOptions.GeoM.Translate(s.Centre.X, s.Centre.Y)

	// vector.DrawFilledRect(screen,
	// 	float32(s.Position.X), float32(s.Position.Y),
	// 	float32(s.Size.W), float32(s.Size.H),
	// 	color.RGBA{255, 128, 0, 255}, false)

	s.DrawOptions.GeoM.Translate(s.Position.X, s.Position.Y)
	colorm.DrawImage(screen, s.Image, s.ColorModel, s.DrawOptions)

	if s.wraparound {
		s.DrawOptions.GeoM.Translate(s.screenBounds.W, 0)
		colorm.DrawImage(screen, s.Image, s.ColorModel, s.DrawOptions)

		s.DrawOptions.GeoM.Translate(-s.screenBounds.W, +s.screenBounds.H)
		colorm.DrawImage(screen, s.Image, s.ColorModel, s.DrawOptions)

		s.DrawOptions.GeoM.Translate(-s.screenBounds.W, -s.screenBounds.H)
		colorm.DrawImage(screen, s.Image, s.ColorModel, s.DrawOptions)

		s.DrawOptions.GeoM.Translate(+s.screenBounds.W, -s.screenBounds.H)
		colorm.DrawImage(screen, s.Image, s.ColorModel, s.DrawOptions)
	}

	s.ColorModel.Reset()
	s.DrawOptions.GeoM.Reset()
}

func (s *Sprite) checkEdges() {
	if s.Position.X > s.screenBounds.W {
		s.Position.X = 0
	} else if s.Position.X < -s.Size.W {
		s.Position.X = s.screenBounds.W - s.Size.W
	}

	if s.Position.Y > s.screenBounds.H {
		s.Position.Y = 0
	} else if s.Position.Y < -s.Size.H {
		s.Position.Y = s.screenBounds.H - s.Size.H
	}
}

func (s *Sprite) MoveForward(acceleration, maxSpeed float64) {
	newVector := geometry.VectorFrom(s.Direction, acceleration)
	newVector.Add(s.Velocity)
	s.Speed = newVector.Magnitude()

	if s.Speed >= maxSpeed {
		newVector.Scale(maxSpeed / s.Speed)
	}
	s.Velocity.X = newVector.X
	s.Velocity.Y = newVector.Y
}
