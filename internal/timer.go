package internal

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Timer struct {
	currentTicks int
	targetTicks  int
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(d.Milliseconds()) * ebiten.TPS() / 1000,
	}
}

func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

func (t *Timer) IsReady() bool {
	return t.currentTicks >= t.targetTicks
}

func (t *Timer) Reset() {
	t.currentTicks = 0
}

func (t *Timer) ResetTarget(d time.Duration) {
	t.currentTicks = 0
	t.targetTicks = int(d.Milliseconds()) * ebiten.TPS() / 1000
}

func (t *Timer) CurrentTicks() int {
	return t.currentTicks
}

func (t *Timer) PercentComplete() float64 {
	return float64(t.currentTicks) / float64(t.targetTicks)
}
