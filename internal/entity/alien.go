package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Alien struct {
	position         geometry.Vector
	velocity         geometry.Vector
	centre           geometry.Vector
	direction        float64
	sprite           *ebiten.Image
	bounds           *geometry.Dimension
	respawnTimer     *internal.Timer
	shootCooldown    *internal.Timer
	bullets          map[int]*Bullet
	sequence         *internal.Sequence
	playerPosition   *geometry.Vector
	shootingAccuracy float64
	maxSalvo         int
}

const respawnDuration = 30 * time.Second

func NewAlien(level int, position *geometry.Vector, playerPosition *geometry.Vector, screenBounds *geometry.Dimension) *Alien {
	return &Alien{
		direction:        0,
		position:         *position,
		sprite:           sprites.AlienSpaceShip,
		centre:           sprites.Centre(sprites.AlienSpaceShip),
		bounds:           screenBounds,
		respawnTimer:     internal.NewTimer(respawnDuration),
		shootCooldown:    internal.NewTimer(5 * time.Second),
		sequence:         internal.NewSequence(),
		bullets:          make(map[int]*Bullet),
		playerPosition:   playerPosition,
		shootingAccuracy: 0.8,
		maxSalvo:         3 + level,
	}
}

func (a *Alien) Draw(screen *ebiten.Image) {
	for _, bullet := range a.bullets {
		bullet.Draw(screen)
	}

	if a.respawnTimer.IsReady() {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(a.position.X, a.position.Y)

		// vector.DrawFilledCircle(screen, float32(a.position.X+a.centre.X), float32(a.position.Y+a.centre.Y), float32(a.Size()), color.RGBA{128, 128, 0, 255}, false)

		screen.DrawImage(a.sprite, op)
	}
}

func (a *Alien) Update() error {
	for idx, bullet := range a.bullets {
		err := bullet.Update()
		if err != nil {
			return err
		}

		if bullet.IsExpired() {
			delete(a.bullets, idx)
		}
	}

	a.respawnTimer.Update()
	if a.respawnTimer.IsReady() {
		a.HandleMovement()
		a.HandleShooting()

		a.position.Add(&a.velocity)
		a.position.CheckEdges(a.bounds, a.centre.X, a.centre.Y)
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

func randomDuration(min, max time.Duration) time.Duration {
	if min > max {
		min, max = max, min
	}

	return time.Duration(rand.Int63n(int64(max-min))) + min
}

func (a *Alien) HandleShooting() {
	a.shootCooldown.Update()
	if a.shootCooldown.IsReady() && len(a.bullets) < a.maxSalvo {
		duration := randomDuration(1*time.Second, 8*time.Second)
		a.shootCooldown.ResetTarget(duration)

		direction := a.position.AngleTo(a.playerPosition) + a.ShootingJitter()
		spawnPosn := &geometry.Vector{
			X: a.position.X + a.centre.X + (math.Cos(direction) * 60),
			Y: a.position.Y + a.centre.Y + (math.Sin(direction) * 60),
		}
		a.bullets[a.sequence.GetNext()] = NewBullet(a.bounds, spawnPosn, direction, sprites.Large)
	}
}

func (a *Alien) ShootingJitter() float64 {
	return (rand.Float64() - 0.5) * (1 - a.shootingAccuracy)
}

func (a *Alien) Value() int {
	return 1000
}

func (a *Alien) Position() *geometry.Vector {
	return geometry.Add(&a.position, &a.centre)
}

func (a *Alien) Size() float64 {
	return a.centre.Y * 0.75
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
