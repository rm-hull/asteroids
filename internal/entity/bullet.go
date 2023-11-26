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
	direction float64
	bounds    *geometry.Dimension
	sprite    *ebiten.Image
	timer     *internal.Timer
}

func NewBullet(screenBounds *geometry.Dimension, position *geometry.Vector, direction float64) *Bullet {
	bulletSpeed := float64(480 / ebiten.TPS())
	return &Bullet{
		direction: direction,
		position:  *position,
		velocity:  geometry.VectorFrom(direction, bulletSpeed),
		sprite:    sprites.Bullet1,
		bounds:    screenBounds,
		timer:     internal.NewTimer(3 * time.Second),
	}
}

var bulletWidth = sprites.Bullet1.Bounds().Dx()
var bulletHeight = sprites.Bullet1.Bounds().Dy()

var bulletHalfW = float64(bulletWidth / 2)
var bulletHalfH = float64(bulletHeight / 2)

func (b *Bullet) Draw(screen *ebiten.Image) {

	if b.IsExpired() {
		return
	}

	cm := colorm.ColorM{}
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(b.position.X-bulletHalfW, b.position.Y-bulletHalfH)

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
		b.position.CheckEdges(b.bounds, bulletHalfW, bulletHalfH)
	}
	return nil
}

func (b *Bullet) IsExpired() bool {
	return b.timer.IsReady()
}
