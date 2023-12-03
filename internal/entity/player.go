package entity

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/rm-hull/asteroids/internal"
	"github.com/rm-hull/asteroids/internal/fonts"
	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/internal/sprites"
	"github.com/rm-hull/asteroids/resources/soundfx"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Player struct {
	sprite           *sprites.Sprite
	deadTimer        *internal.Timer
	cannotDieTimer   *internal.Timer
	shootCooldown    *internal.Timer
	screenBounds     *geometry.Dimension
	livesLeft        int
	score            int
	bullets          map[int]*Bullet
	sequence         *internal.Sequence
	godMode          bool
	maxSalvo         int
	shootingAccuracy float64
}

const (
	numLives           = 3
	maxSpeed           = 5.0
	blastRadius        = 40.0
	deathDuration      = 2 * time.Second
	cannotDieDuration  = 3 * time.Second
	cooldownTime       = 100 * time.Millisecond
	sampleRate         = 44100
	extraLifeThreshold = 10000
)

var audioContext = audio.NewContext(sampleRate)

func NewPlayer(screenBounds *geometry.Dimension) *Player {
	sprite := sprites.NewSprite(screenBounds, sprites.SpaceShip1, true)
	sprite.Position.X = screenBounds.W/2 - sprite.Centre.X
	sprite.Position.Y = screenBounds.H/2 - sprite.Centre.Y

	return &Player{
		sprite:           sprite,
		cannotDieTimer:   internal.NewTimer(cannotDieDuration),
		shootCooldown:    internal.NewTimer(cooldownTime),
		screenBounds:     screenBounds,
		livesLeft:        numLives,
		score:            0,
		bullets:          make(map[int]*Bullet),
		sequence:         internal.NewSequence(),
		maxSalvo:         3,
		shootingAccuracy: 1.0,
		godMode:          false,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("LIVES: %d", p.livesLeft), fonts.AsteroidsDisplayFont16, 0, 30, color.White)
	text.Draw(screen, fmt.Sprintf("SCORE: %d", p.score), fonts.AsteroidsDisplayFont16, 350, 30, color.White)
	text.Draw(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), fonts.AsteroidsDisplayFont16, 700, 30, color.White)

	if p.livesLeft == 0 {
		message := "GAME OVER"
		bounds, _ := font.BoundString(fonts.AsteroidsDisplayFont32, message)
		dx := float64(bounds.Max.X.Round() - bounds.Min.X.Round())
		dy := float64(bounds.Max.Y.Round() - bounds.Min.Y.Round())

		x := int(p.screenBounds.W-dx) / 2
		y := int(p.screenBounds.H-dy) / 2

		text.Draw(screen, message, fonts.AsteroidsDisplayFont32, x, y, color.White)
		return
	}

	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}

	if p.IsDying() {
		fade := 1.0 - p.deadTimer.PercentComplete()
		p.sprite.ColorModel.Scale(1.0, 1.0, 1.0, fade)
	} else if p.CannotDie() {
		fade := p.cannotDieTimer.PercentComplete()
		p.sprite.ColorModel.Scale(1.0, 1.0, 1.0, fade)
	}

	p.sprite.Draw(screen)
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

	if p.IsDying() {
		p.SpinOutOfControl()
	} else {
		p.HandleMovement()
		p.HandleShooting()

		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			p.ToggleGodMode()
		}
	}

	p.cannotDieTimer.Update()
	if err := p.sprite.Update(); err != nil {
		return err
	}

	return nil
}

func (p *Player) HandleMovement() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.sprite.Direction -= math.Pi / float64(ebiten.TPS())
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.sprite.Direction += math.Pi / float64(ebiten.TPS())
	}
	p.sprite.Orientation = p.sprite.Direction

	// Thrusting?
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.sprite.MoveForward(0.2, maxSpeed)
		p.sprite.Image = sprites.SpaceShip2

	} else {
		// Back to normal
		p.sprite.Image = sprites.SpaceShip1
	}
}

func (p *Player) HandleShooting() {
	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && len(p.bullets) < p.maxSalvo && (ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeySpace)) {
		p.shootCooldown.Reset()

		direction := p.sprite.Direction + p.ShootingJitter()
		spawnPosn := geometry.Add(p.Position(), geometry.VectorFrom(p.sprite.Direction, blastRadius))
		p.bullets[p.sequence.GetNext()] = NewBullet(p.screenBounds, spawnPosn, direction, sprites.Small)

		sfxPlayer := audioContext.NewPlayerFromBytes(soundfx.LazerGunShot1)
		sfxPlayer.Play()
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
	p.sprite.Orientation += 3 * math.Pi / float64(ebiten.TPS())
	p.deadTimer.Update()

	if p.deadTimer.IsReady() {
		p.Prepare()
		p.livesLeft--
	}
}

func (p *Player) Prepare() {
	p.deadTimer = nil
	p.sprite.Reset()
	p.sprite.Position.X = p.screenBounds.W/2 - p.sprite.Centre.X
	p.sprite.Position.Y = p.screenBounds.H/2 - p.sprite.Centre.Y
	p.sprite.Image = sprites.SpaceShip1
	p.cannotDieTimer.Reset()
	for idx := range p.bullets {
		delete(p.bullets, idx)
	}
}

func (p *Player) Kill() {
	if p.CannotDie() {
		return
	}
	p.deadTimer = internal.NewTimer(deathDuration)

	sePlayer := audioContext.NewPlayerFromBytes(soundfx.Explosion1)
	sePlayer.Play()
}

func (p *Player) NotNear() *geometry.Vector {
	sqHalfH := math.Pow(p.screenBounds.H/2, 2)
	for {
		position := geometry.Vector{
			X: rand.Float64() * p.screenBounds.W,
			Y: rand.Float64() * p.screenBounds.H,
		}

		if p.sprite.Position.SquareDistanceFrom(&position) > sqHalfH {
			return &position
		}
	}
}

func (p *Player) UpdateScore(value int) {
	if math.Mod(float64(p.score), extraLifeThreshold) > math.Mod(float64(p.score+value), extraLifeThreshold) {
		p.livesLeft++

		sfxPlayer := audioContext.NewPlayerFromBytes(soundfx.ExtraLife)
		sfxPlayer.Play()
	}

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

func (p *Player) Size() float64 {
	return p.sprite.Centre.Y * 0.65
}

func (p *Player) Position() *geometry.Vector {
	return geometry.Add(p.sprite.Position, p.sprite.Centre).Mod(p.screenBounds)
}
