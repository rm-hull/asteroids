package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Alien struct {
	position     geometry.Vector
	velocity     geometry.Vector
	direction    float64
	sprite       *ebiten.Image
	bounds       *geometry.Dimension
	respawnTimer *internal.Timer
}

const respawnDuration = 30 * time.Second

func NewAlien(screenBounds *geometry.Dimension, position *geometry.Vector) *Alien {
	return &Alien{
		direction:    0,
		position:     *position,
		sprite:       sprites.AlienSpaceShip,
		bounds:       screenBounds,
		respawnTimer: internal.NewTimer(respawnDuration),
	}
}

func (a *Alien) Draw(screen *ebiten.Image) {
	if a.respawnTimer.IsReady() {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(a.position.X, a.position.Y)

		screen.DrawImage(a.sprite, op)
	}
}

func (a *Alien) Update() error {
	a.respawnTimer.Update()
	if a.respawnTimer.IsReady() {
		a.HandleMovement()

		a.position.Add(&a.velocity)
		a.position.CheckEdges(a.bounds, float64(spaceshipWidth), float64(spaceshipHeight))
	}
	return nil
}

func (a *Alien) HandleMovement() {
	delta := rand.Float64() - 0.3
	a.direction += delta

	thrusting := rand.Float64() > 0.3
	if thrusting {
		newVector := geometry.VectorFrom(a.direction, 0.3)
		newVector.Add(&a.velocity)
		speed := newVector.Magnitude()

		if speed < maxSpeed {
			a.velocity = newVector
		} else {
			newVector.Scale(maxSpeed / speed)
			a.velocity = newVector
		}
	}
}

func (a *Alien) Value() int {
	return 1000
}

func (a *Alien) Bounds() *image.Rectangle {
	point := image.Point{X: int(a.position.X), Y: int(a.position.Y)}
	return &image.Rectangle{
		Min: point,
		Max: a.sprite.Bounds().Max.Add(point),
	}
}

func (a *Alien) Kill() {
	a.respawnTimer.Reset()
}

func (a *Alien) IsAlive() bool {
	return a.respawnTimer.IsReady()
}
