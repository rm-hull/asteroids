package entity

import (
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const asteroidMaxSpeed = 2

type Asteroid struct {
	position  geometry.Vector
	velocity  geometry.Vector
	rotation  float64
	direction float64
	bounds    *geometry.Dimension
	sprite    *ebiten.Image
}

func NewAsteroid(size int, playerPosition *geometry.Vector, screenBounds *geometry.Dimension) *Asteroid {

	direction := rand.Float64() * 2 * math.Pi
	speed := rand.Float64() * asteroidMaxSpeed

	return &Asteroid{
		position: notNear(playerPosition, screenBounds),
		velocity: geometry.VectorFrom(direction, speed),
		rotation: (rand.Float64() - 0.5) / 20,
		bounds:   screenBounds,
		sprite:   sprites.Asteroid(size),
	}
}

func notNear(v *geometry.Vector, bounds *geometry.Dimension) geometry.Vector {
	halfH := bounds.H / 2
	for {
		position := geometry.Vector{
			X: rand.Float64() * bounds.W,
			Y: rand.Float64() * bounds.H,
		}

		if v.DistanceFrom(&position) > halfH {
			return position
		}
	}
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	halfW := float64(a.sprite.Bounds().Dx()) / 2
	halfH := float64(a.sprite.Bounds().Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(a.direction)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(a.position.X, a.position.Y)

	screen.DrawImage(a.sprite, op)
}

func (a *Asteroid) Update() error {
	a.direction += a.rotation
	a.position.Accumulate(&a.velocity)
	a.position.CheckEdges(a.bounds, float64(a.sprite.Bounds().Dx()), float64(a.sprite.Bounds().Dy()))
	return nil
}
