package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	position geometry.Vector
	velocity geometry.Vector
	timer    *internal.Timer
	bounds   *geometry.Dimension
	message  string
	current  int
}

func NewLevel(screenBounds *geometry.Dimension) *Level {
	return &Level{
		position: geometry.Vector{
			X: screenBounds.W / 2,
			Y: screenBounds.H / 2,
		},
		velocity: geometry.VectorFrom(-math.Pi/2, 0.4),
		timer:    internal.NewTimer(3 * time.Second),
		bounds:   screenBounds,
		message:  "Level: 1",
		current:    1,
	}
}

func (l *Level) Draw(screen *ebiten.Image) {
	if l.IsExpired() {
		return
	}

	// cm := colorm.ColorM{}
	// op := &colorm.DrawImageOptions{}
	// op.GeoM.Translate(l.position.X, l.position.Y)

	// if pctComplete := l.timer.PercentComplete(); pctComplete > 0.75 {
	// 	fade := ((1.0 - pctComplete) / 0.25)
	// 	cm.Scale(1, 1, 1, fade)
	// }

	// colorm.DrawImage(screen, b.sprite, cm, op)
	ebitenutil.DebugPrintAt(screen, l.message, int(l.position.X), int(l.position.Y))
}

func (l *Level) Update() error {
	l.timer.Update()
	if !l.IsExpired() {
		l.position.Add(&l.velocity)
	}
	return nil
}

func (l *Level) IsExpired() bool {
	return l.timer.IsReady()
}

func (l *Level) Current() int {
	return l.current
}

func (l *Level) Next() {
	l.Reset(l.current + 1)
}

func (l *Level) Reset(level int) {
	l.position.X = l.bounds.W / 2
	l.position.Y = l.bounds.H / 2
	l.current = level
	l.message = fmt.Sprintf("Level: %d", l.current)
	l.timer.Reset()
}
