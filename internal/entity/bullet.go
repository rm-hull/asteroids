package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Bullet struct {
	position  geometry.Vector
	velocity  geometry.Vector
	centre    geometry.Vector
	direction float64
	bounds    *geometry.Dimension
	sprite    *ebiten.Image
	timer     *internal.Timer
	directHit bool
}

func NewBullet(screenBounds *geometry.Dimension, position *geometry.Vector, direction float64, size int) *Bullet {
	bulletSpeed := float64(480 / ebiten.TPS())
	sprite := sprites.Bullet(size)
	centre := sprites.Centre(sprite)

	return &Bullet{
		direction: direction,
		position:  geometry.Vector{X: position.X - centre.X, Y: position.Y - centre.Y},
		velocity:  geometry.VectorFrom(direction, bulletSpeed),
		centre:    centre,
		sprite:    sprite,
		bounds:    screenBounds,
		timer:     internal.NewTimer(2 * time.Second),
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.IsExpired() {
		return
	}

	cm := colorm.ColorM{}
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(b.position.X, b.position.Y)

	if pctComplete := b.timer.PercentComplete(); pctComplete > 0.75 {
		fade := ((1.0 - pctComplete) / 0.25)
		cm.Scale(1, 1, 1, fade)
	}

	colorm.DrawImage(screen, b.sprite, cm, op)
}

func (b *Bullet) Update() error {
	b.timer.Update()
	if !b.IsExpired() {
		b.position.Add(&b.velocity)
	}
	return nil
}

func (b *Bullet) IsExpired() bool {
	return b.directHit || b.timer.IsReady()
}

func (b *Bullet) Position() *geometry.Vector {
	return geometry.Add(&b.position, &b.centre)
}

func (b *Bullet) Size() float64 {
	return b.centre.X / 2
}

func (b *Bullet) CollisionDetected(collider Collider) bool {
	if b.timer.PercentComplete() < 90 && !b.directHit {
		if hit := CollisionDetected(b, collider); hit {
			b.directHit = true
			return true
		}
	}
	return false
}
