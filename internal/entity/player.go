package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	position  geometry.Vector
	velocity  geometry.Vector
	direction float64
	speed     float64
	sprite    *ebiten.Image
	deadTimer *internal.Timer
	bounds    geometry.Dimension
	livesLeft int
}

const maxSpeed = 5

var spaceshipWidth = sprites.SpaceShip1.Bounds().Dx()
var spaceshipHeight = sprites.SpaceShip1.Bounds().Dy()

var spaceshipHalfW = float64(spaceshipWidth / 2)
var spaceshipHalfH = float64(spaceshipHeight / 2)

func NewPlayer(screenBounds geometry.Dimension) *Player {
	return &Player{
		direction: 0,
		speed:     0,
		position: geometry.Vector{
			X: screenBounds.W/2 - spaceshipHalfW,
			Y: screenBounds.H/2 - spaceshipHalfH,
		},
		sprite:    sprites.SpaceShip1,
		bounds:    screenBounds,
		livesLeft: 3,
	}
}

func (p *Player) Position() *geometry.Vector {
	return &p.position
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.livesLeft == 0 {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", 0, 0)
		return
	}

	cm := colorm.ColorM{}
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(-spaceshipHalfW, -spaceshipHalfH)
	op.GeoM.Rotate(p.direction)
	op.GeoM.Translate(spaceshipHalfW, spaceshipHalfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	if p.deadTimer != nil {
		fade := 1.0 - p.deadTimer.PercentComplete()
		cm.Scale(1.0, 1.0, 1.0, fade)
	}

	colorm.DrawImage(screen, p.sprite, cm, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: (%d,%d)", int(p.position.X), int(p.position.Y)), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %0.2f", p.speed), 150, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lives: %d", p.livesLeft), 250, 0)
}

func (p *Player) Update() error {
	if p.livesLeft == 0 {
		// TODO: Update game state
		return nil
	}

	// Temporary -
	if ebiten.IsKeyPressed(ebiten.KeyK) && p.deadTimer == nil {
		p.Kill()
	}

	if p.deadTimer != nil {
		p.SpinOutOfControl()
	} else {
		p.HandleMovement()
	}

	p.position.Add(&p.velocity)
	p.position.CheckEdges(&p.bounds, float64(spaceshipWidth), float64(spaceshipHeight))

	return nil
}

func (p *Player) HandleMovement() {
	if p.livesLeft == 0 {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.direction -= math.Pi / float64(ebiten.TPS())
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.direction += math.Pi / float64(ebiten.TPS())
	}

	// Thrusting?
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		newVector := geometry.VectorFrom(p.direction, 0.2)
		newVector.Add(&p.velocity)
		p.speed = newVector.Magnitude()

		if p.speed < maxSpeed {
			p.velocity = newVector
		} else {
			newVector.Scale(maxSpeed / p.speed)
			p.velocity = newVector
		}
		p.sprite = sprites.SpaceShip2

	} else {
		// Back to normal
		p.sprite = sprites.SpaceShip1
	}
}

func (p *Player) SpinOutOfControl() {
	p.direction += 3 * math.Pi / float64(ebiten.TPS())
	p.deadTimer.Update()

	if p.deadTimer.IsReady() {
		p.Reset()
	}
}

func (p *Player) Reset() {
	p.deadTimer = nil
	p.direction = 0
	p.speed = 0
	p.velocity.X = 0
	p.velocity.Y = 0
	p.position.X = p.bounds.W/2 - spaceshipHalfW
	p.position.Y = p.bounds.H/2 - spaceshipHalfW
	p.sprite = sprites.SpaceShip1
	p.livesLeft--
}

func (p *Player) Kill() {
	p.deadTimer = internal.NewTimer(2 * time.Second)
}
