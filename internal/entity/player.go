package entity

import (
	"asteroids/internal"
	"asteroids/internal/fonts"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Player struct {
	position         geometry.Vector
	velocity         geometry.Vector
	direction        float64
	speed            float64
	sprite           *ebiten.Image
	deadTimer        *internal.Timer
	cannotDieTimer   *internal.Timer
	shootCooldown    *internal.Timer
	bounds           *geometry.Dimension
	livesLeft        int
	score            int
	bullets          map[int]*Bullet
	sequence         *internal.Sequence
	godMode          bool
	maxSalvo         int
	shootingAccuracy float64
}

const numLives = 3
const maxSpeed = 5.0
const blastRadius = 40.0

var spaceshipWidth = sprites.SpaceShip1.Bounds().Dx()
var spaceshipHeight = sprites.SpaceShip1.Bounds().Dy()

var spaceshipHalfW = float64(spaceshipWidth / 2)
var spaceshipHalfH = float64(spaceshipHeight / 2)

const deathDuration = 2 * time.Second
const cannotDieDuration = 3 * time.Second
const cooldownTime = 100 * time.Millisecond

func NewPlayer(screenBounds *geometry.Dimension) *Player {
	return &Player{
		direction: 0,
		speed:     0,
		position: geometry.Vector{
			X: screenBounds.W/2 - spaceshipHalfW,
			Y: screenBounds.H/2 - spaceshipHalfH,
		},
		cannotDieTimer:   internal.NewTimer(cannotDieDuration),
		shootCooldown:    internal.NewTimer(cooldownTime),
		sprite:           sprites.SpaceShip1,
		bounds:           screenBounds,
		livesLeft:        numLives,
		score:            0,
		bullets:          make(map[int]*Bullet),
		sequence:         internal.NewSequence(),
		maxSalvo:         3,
		shootingAccuracy: 1.0,
	}
}

func (p *Player) CurrentPosition() *geometry.Vector {
	return &p.position
}

func (p *Player) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("LIVES: %d", p.livesLeft), fonts.AsteroidsDisplayFont16, 0, 30, color.White)
	text.Draw(screen, fmt.Sprintf("SCORE: %d", p.score), fonts.AsteroidsDisplayFont16, 350, 30, color.White)

	if p.livesLeft == 0 {
		message := "GAME OVER"
		rect := text.BoundString(fonts.AsteroidsDisplayFont32, message)

		x := int(p.bounds.W-float64(rect.Dx())) / 2
		y := int(p.bounds.H-float64(rect.Dy())) / 2

		text.Draw(screen, message, fonts.AsteroidsDisplayFont32, x, y, color.White)
		return
	}

	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}

	cm := colorm.ColorM{}
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(-spaceshipHalfW, -spaceshipHalfH)
	op.GeoM.Rotate(p.direction)
	op.GeoM.Translate(spaceshipHalfW, spaceshipHalfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	if p.IsDying() {
		fade := 1.0 - p.deadTimer.PercentComplete()
		cm.Scale(1.0, 1.0, 1.0, fade)
	} else if p.CannotDie() {
		fade := p.cannotDieTimer.PercentComplete()
		cm.Scale(1.0, 1.0, 1.0, fade)
	}

	colorm.DrawImage(screen, p.sprite, cm, op)
}

func (p *Player) NoseTip() *geometry.Vector {
	return &geometry.Vector{
		X: p.position.X + spaceshipHalfW + (math.Cos(p.direction) * blastRadius),
		Y: p.position.Y + spaceshipHalfH + (math.Sin(p.direction) * blastRadius),
	}
}

func (p *Player) Update() error {
	if p.livesLeft == 0 {
		// TODO: Update game state
		return nil
	}

	for idx, bullet := range p.bullets {
		err := bullet.Update()
		if err != nil {
			return err
		}

		if bullet.IsExpired() {
			delete(p.bullets, idx)
		}
	}

	if p.deadTimer != nil {
		p.SpinOutOfControl()
	} else {
		p.HandleMovement()
	}

	p.HandleShooting()

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		p.ToggleGodMode()
	}

	p.cannotDieTimer.Update()
	p.position.Add(&p.velocity)
	p.position.CheckEdges(p.bounds, float64(spaceshipWidth), float64(spaceshipHeight))

	return nil
}

func (p *Player) HandleMovement() {
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

func (p *Player) HandleShooting() {
	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && len(p.bullets) < p.maxSalvo && (ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeySpace)) {
		p.shootCooldown.Reset()
		p.bullets[p.sequence.GetNext()] = NewBullet(p.bounds, p.NoseTip(), p.direction+p.ShootingJitter(), sprites.Small)
	}
}

func (p *Player) ToggleGodMode() {
	if p.godMode {
		p.godMode = false
		p.maxSalvo = 3
		p.shootCooldown.ResetTarget(cooldownTime)
		p.shootingAccuracy = 1.0
	} else {
		p.godMode = true
		p.maxSalvo = 200
		p.shootCooldown.ResetTarget(50 * time.Millisecond)
		p.shootingAccuracy = 0.6
	}
}

func (p *Player) ShootingJitter() float64 {
	return (rand.Float64() - 0.5) * (1 - p.shootingAccuracy)
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
	p.cannotDieTimer.Reset()
	p.livesLeft--
}

func (p *Player) Kill() {
	if p.CannotDie() {
		return
	}
	p.deadTimer = internal.NewTimer(deathDuration)
}

func (p *Player) NotNear() *geometry.Vector {
	halfH := p.bounds.H / 3
	for {
		position := geometry.Vector{
			X: rand.Float64() * p.bounds.W,
			Y: rand.Float64() * p.bounds.H,
		}

		if p.position.DistanceFrom(&position) > halfH {
			return &position
		}
	}
}

func (p *Player) UpdateScore(value int) {
	p.score += value
}

func (p *Player) IsDying() bool {
	return p.deadTimer != nil
}

func (p *Player) IsAlive() bool {
	return !p.IsDying() && !p.CannotDie()
}

func (p *Player) CannotDie() bool {
	return p.godMode || !p.cannotDieTimer.IsReady()
}

func (p *Player) Bullets(callback func(bullet *Bullet)) {
	for _, bullet := range p.bullets {
		callback(bullet)
	}
}

func (p *Player) Bounds() *image.Rectangle {
	point := image.Point{X: int(p.position.X), Y: int(p.position.Y)}
	return &image.Rectangle{
		Min: point,
		Max: p.sprite.Bounds().Max.Add(point),
	}
}
