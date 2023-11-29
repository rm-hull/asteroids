package entity

import (
	"asteroids/internal"
	"asteroids/internal/fonts"
	"asteroids/internal/geometry"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	level := &Level{
		velocity: geometry.VectorFrom(-math.Pi/2, 0.4),
		timer:    internal.NewTimer(3 * time.Second),
		bounds:   screenBounds,
	}

	level.Reset(1)
	return level
}

func (l *Level) Draw(screen *ebiten.Image) {
	if l.IsExpired() {
		return
	}
	text.Draw(screen, l.message, fonts.AsteroidsDisplayFont32, int(l.position.X), int(l.position.Y), color.White)
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

	l.message = fmt.Sprintf("LEVEL %d", l.current)
	rect := text.BoundString(fonts.AsteroidsDisplayFont32, l.message)

	l.position.X = (l.bounds.W - float64(rect.Dx())) / 2
	l.position.Y = (l.bounds.H - float64(rect.Dy())) / 2

	l.timer.Reset()
}
