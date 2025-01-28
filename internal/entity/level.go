package entity

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/rm-hull/asteroids/internal"
	"github.com/rm-hull/asteroids/internal/fonts"
	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/internal/text_align"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Level struct {
	position geometry.Vector
	velocity *geometry.Vector
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
	op := &text.DrawOptions{}
	op.GeoM.Translate(l.position.X, l.position.Y)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, l.message, fonts.AsteroidsFace64, op)
}

func (l *Level) Update() error {
	l.timer.Update()
	if !l.IsExpired() {
		l.position.Add(l.velocity)
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
	x, y := text_align.Center(l.bounds, l.message, fonts.AsteroidsFace64)

	l.position.X = float64(x)
	l.position.Y = float64(y)

	l.timer.Reset()
}
