package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Alien struct {
	position      geometry.Vector
	velocity      geometry.Vector
	direction     float64
	sprite        *ebiten.Image
	bounds        *geometry.Dimension
	respawnTimer  *internal.Timer
	shootCooldown *internal.Timer
	bullets       map[int]*Bullet
	sequence      *internal.Sequence
}

const respawnDuration = 30 * time.Second

func NewAlien(screenBounds *geometry.Dimension, position *geometry.Vector) *Alien {
	return &Alien{
		direction:     0,
		position:      *position,
		sprite:        sprites.AlienSpaceShip,
		bounds:        screenBounds,
		respawnTimer:  internal.NewTimer(respawnDuration),
		shootCooldown: internal.NewTimer(5 * time.Second),
		sequence:      internal.NewSequence(),
		bullets:       make(map[int]*Bullet),
	}
}

func (a *Alien) Draw(screen *ebiten.Image) {
	if a.respawnTimer.IsReady() {

		for _, bullet := range a.bullets {
			bullet.Draw(screen)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(a.position.X, a.position.Y)

		screen.DrawImage(a.sprite, op)
	}
}

func (a *Alien) Update() error {
	a.respawnTimer.Update()
	if a.respawnTimer.IsReady() {

		for idx, bullet := range a.bullets {
			err := bullet.Update()
			if err != nil {
				return err
			}

			if bullet.IsExpired() {
				delete(a.bullets, idx)
			}
		}

		a.HandleMovement()
		a.HandleShooting()

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

func (a *Alien) HandleShooting() {
	a.shootCooldown.Update()
	if a.shootCooldown.IsReady() && len(a.bullets) < maxSalvo {
		a.shootCooldown.Reset()

		bounds := a.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		direction := rand.Float64() * 2 * math.Pi
		spawnPosn := &geometry.Vector{
			X: a.position.X + halfW + (math.Cos(direction) * 60),
			Y: a.position.Y + halfH + (math.Sin(direction) * 60),
		}
		a.bullets[a.sequence.GetNext()] = NewBullet(a.bounds, spawnPosn, direction, sprites.Large)
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

func (a *Alien) Bullets(callback func(bullet *Bullet)) {
	for _, bullet := range a.bullets {
		callback(bullet)
	}
}