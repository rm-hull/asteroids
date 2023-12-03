package entity

import (
	"time"

	"github.com/rm-hull/asteroids/internal"
	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/internal/sprites"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	sprite       *sprites.Sprite
	timer        *internal.Timer
	screenBounds *geometry.Dimension
	directHit    bool
}

func NewBullet(screenBounds *geometry.Dimension, position *geometry.Vector, direction float64, size int) *Bullet {
	bulletSpeed := float64(480 / ebiten.TPS())
	sprite := sprites.NewSprite(screenBounds, sprites.Bullet(size), false)
	sprite.Direction = direction
	sprite.Position.X = position.X - sprite.Centre.X
	sprite.Position.Y = position.Y - sprite.Centre.Y
	sprite.Velocity = geometry.VectorFrom(direction, bulletSpeed)

	return &Bullet{
		sprite:       sprite,
		timer:        internal.NewTimer(2 * time.Second),
		screenBounds: screenBounds,
		directHit:    false,
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.IsExpired() {
		return
	}

	if pctComplete := b.timer.PercentComplete(); pctComplete > 0.75 {
		fade := ((1.0 - pctComplete) / 0.25)
		b.sprite.ColorModel.Scale(1, 1, 1, fade)
	}

	b.sprite.Draw(screen)
}

func (b *Bullet) Update() error {
	b.timer.Update()
	if !b.IsExpired() {
		if err := b.sprite.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bullet) IsExpired() bool {
	return b.directHit || b.timer.IsReady()
}

func (b *Bullet) Position() *geometry.Vector {
	return geometry.Add(b.sprite.Position, b.sprite.Centre)
}

func (b *Bullet) Size() float64 {
	return b.sprite.Centre.X / 2
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
