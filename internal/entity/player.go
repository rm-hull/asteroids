package entity

import (
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	position  geometry.Vector
	velocity  geometry.Vector
	direction float64
	speed     float64
	thrusting bool
	bounds    geometry.Dimension
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
		thrusting: false,
		bounds:    screenBounds,
	}
}

func (p *Player) Position() *geometry.Vector {
	return &p.position
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-spaceshipHalfW, -spaceshipHalfH)
	op.GeoM.Rotate(p.direction)
	op.GeoM.Translate(spaceshipHalfW, spaceshipHalfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	if p.thrusting {
		screen.DrawImage(sprites.SpaceShip2, op)
	} else {
		screen.DrawImage(sprites.SpaceShip1, op)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: (%d,%d)", int(p.position.X), int(p.position.Y)), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %0.2f", p.speed), 150, 0)
}

func (p *Player) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.direction -= math.Pi / float64(ebiten.TPS())
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.direction += math.Pi / float64(ebiten.TPS())
	}

	var keyVector geometry.Vector

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		keyVector = geometry.VectorFrom(p.direction, 0.2)
		newVector := geometry.Add(&p.velocity, &keyVector)
		p.speed = newVector.Magnitude()

		if p.speed < maxSpeed {
			p.velocity = newVector
		} else {
			newVector.Scale(maxSpeed / p.speed)
			p.velocity = newVector
		}
		p.thrusting = true

	} else {
		p.thrusting = false
	}

	p.position.Accumulate(&p.velocity)
	p.position.CheckEdges(&p.bounds, float64(spaceshipWidth), float64(spaceshipHeight))

	return nil
}
