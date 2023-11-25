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

var width = sprites.SpaceShip1.Bounds().Dx()
var height = sprites.SpaceShip1.Bounds().Dy()

var halfW = float64(width / 2)
var halfH = float64(height / 2)

func NewPlayer(screenBounds geometry.Dimension) *Player {
	return &Player{
		direction: 0,
		speed:     0,
		position: geometry.Vector{
			X: screenBounds.W/2 - halfW,
			Y: screenBounds.H/2 - halfH,
		},
		thrusting: false,
		bounds:    screenBounds,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.direction)
	op.GeoM.Translate(halfW, halfH)

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
		p.thrusting = true
	} else {
		keyVector = geometry.Zero
		p.thrusting = false
	}

	newVector := geometry.Add(&p.velocity, &keyVector)
	p.speed = newVector.Magnitude()

	if p.speed < maxSpeed {
		p.velocity = newVector
	} else {
		newVector.Scale(maxSpeed / p.speed)
		p.velocity = newVector
	}

	p.position.Accumulate(&p.velocity)

	p.position.CheckEdges(&p.bounds, float64(width), float64(height))

	return nil
}
