package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Bullet struct {
	position     geometry.Vector
	velocity     geometry.Vector
	direction    float64
	bounds       *geometry.Dimension
	sprite       *ebiten.Image
	timer        *internal.Timer
	directHit    bool
}

func NewBullet(screenBounds *geometry.Dimension, position *geometry.Vector, direction float64, size int) *Bullet {
	bulletSpeed := float64(480 / ebiten.TPS())
	bounds := sprites.Bullet1.Bounds()
	halfW := float64(bounds.Dx() / 2)
	halfH := float64(bounds.Dy() / 2)

	return &Bullet{
		direction:    direction,
		position:     geometry.Vector{X: position.X - halfW, Y: position.Y - halfH},
		velocity:     geometry.VectorFrom(direction, bulletSpeed),
		sprite:       sprites.Bullet(size),
		bounds:       screenBounds,
		timer:        internal.NewTimer(3 * time.Second),
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
	// bounds2 := b.Bounds()
	// ebitenutil.DrawRect(screen, float64(bounds2.Min.X), float64(bounds2.Min.Y), float64(bounds2.Dx()), float64(bounds2.Dy()), color.RGBA{128, 255, 0, 88})

	colorm.DrawImage(screen, b.sprite, cm, op)
}

func (b *Bullet) Update() error {
	b.timer.Update()
	if !b.IsExpired() {

		bounds := b.sprite.Bounds()
		halfW := float64(bounds.Dx() / 2)
		halfH := float64(bounds.Dy() / 2)

		b.position.Add(&b.velocity)
		b.position.CheckEdges(b.bounds, halfW, halfH)
	}
	return nil
}

func (b *Bullet) IsExpired() bool {
	return b.directHit || b.timer.IsReady()
}

func (b *Bullet) Bounds() *image.Rectangle {
	point := image.Point{X: int(b.position.X), Y: int(b.position.Y)}
	return &image.Rectangle{
		Min: point,
		Max: b.sprite.Bounds().Max.Add(point),
	}
}

func (b *Bullet) AsteroidCollisionDetection(asteroids map[int]*Asteroid) []*Asteroid {
	if b.timer.PercentComplete() < 90 {
		bounds := b.Bounds()
		for _, asteroid := range asteroids {
			if bounds.In(*asteroid.Bounds()) {
				b.directHit = true
				return asteroid.Explode()
			}
		}
	}

	return make([]*Asteroid, 0)
}

func (b *Bullet) AlienCollisionDetection(alien *Alien) bool {
	if b.timer.PercentComplete() < 90 && alien.IsAlive() {
		if hit := b.Bounds().In(*alien.Bounds()); hit {
			b.directHit = true
			alien.Kill()
			return true
		}
	}
	return false
}
