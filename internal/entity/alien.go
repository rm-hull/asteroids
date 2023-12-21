package entity

import (
	"math"
	"math/rand"
	"time"

	"github.com/rm-hull/asteroids/internal"
	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/internal/sprites"
	"github.com/rm-hull/asteroids/resources/soundfx"

	"github.com/hajimehoshi/ebiten/v2"
)

type Alien struct {
	sprite           *sprites.Sprite
	screenBounds     *geometry.Dimension
	deadTimer        *internal.Timer
	respawnTimer     *internal.Timer
	shootCooldown    *internal.Timer
	bullets          map[int]*Bullet
	sequence         *internal.Sequence
	playerPosition   func() *geometry.Vector
	shootingAccuracy float64
	maxSalvo         int
}

const respawnDuration = 30 * time.Second

func NewAlien(level int, position *geometry.Vector, playerPosition func() *geometry.Vector, screenBounds *geometry.Dimension) *Alien {
	sprite := sprites.NewSprite(screenBounds, sprites.AlienSpaceShip, true)
	sprite.Position = position

	return &Alien{
		sprite:           sprite,
		screenBounds:     screenBounds,
		respawnTimer:     internal.NewTimer(respawnDuration),
		shootCooldown:    internal.NewTimer(5 * time.Second),
		sequence:         internal.NewSequence(),
		bullets:          make(map[int]*Bullet),
		playerPosition:   playerPosition,
		shootingAccuracy: 0.8,
		maxSalvo:         3 + level,
	}
}

func (a *Alien) Draw(screen *ebiten.Image) {
	for _, bullet := range a.bullets {
		bullet.Draw(screen)
	}

	if a.respawnTimer.IsReady() {

		if a.IsDying() {
			fade := 1.0 - a.deadTimer.PercentComplete()
			a.sprite.ColorModel.Scale(1.0, 1.0, 1.0, fade)
		}

		a.sprite.Draw(screen)
	}
}

func (a *Alien) Update() error {
	for idx, bullet := range a.bullets {
		err := bullet.Update()
		if err != nil {
			return err
		}

		if bullet.IsExpired() {
			delete(a.bullets, idx)
		}
	}

	a.respawnTimer.Update()
	if a.respawnTimer.IsReady() {

		if a.IsDying() {
			a.SpinOutOfControl()
		} else {
			a.HandleMovement()
			a.HandleShooting()
		}
		if err := a.sprite.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Alien) HandleMovement() {
	delta := (rand.Float64() - 0.5) * 0.6
	a.sprite.Direction += delta

	thrusting := rand.Float64() > 0.3
	if thrusting {
		a.sprite.MoveForward(0.3, maxSpeed)
	}
}

func randomDuration(min, max time.Duration) time.Duration {
	if min > max {
		min, max = max, min
	}

	return time.Duration(rand.Int63n(int64(max-min))) + min
}

func (a *Alien) HandleShooting() {
	a.shootCooldown.Update()
	if a.shootCooldown.IsReady() && len(a.bullets) < a.maxSalvo {
		duration := randomDuration(1*time.Second, 8*time.Second)
		a.shootCooldown.ResetTarget(duration)

		direction := a.sprite.Position.AngleTo(a.playerPosition()) + a.ShootingJitter()
		spawnPosn := geometry.Add(a.Position(), geometry.VectorFrom(direction, 60))
		a.bullets[a.sequence.GetNext()] = NewBullet(a.screenBounds, spawnPosn, direction, sprites.Large)

		sfxPlayer := audioContext.NewPlayerFromBytes(soundfx.LazerGunShot2)
		sfxPlayer.SetVolume(0.5)
		sfxPlayer.Play()
	}
}

func (a *Alien) ShootingJitter() float64 {
	return (rand.Float64() - 0.5) * (1 - a.shootingAccuracy)
}

func (a *Alien) Value() int {
	return 1000
}

func (a *Alien) Position() *geometry.Vector {
	return geometry.Add(a.sprite.Position, a.sprite.Centre).Mod(a.screenBounds)
}

func (a *Alien) Size() float64 {
	return a.sprite.Centre.Y * 0.75
}

func (a *Alien) Kill() {
	a.deadTimer = internal.NewTimer(deathDuration)

	sePlayer := audioContext.NewPlayerFromBytes(soundfx.Explosion2)
	sePlayer.SetVolume(0.15)
	sePlayer.Play()
}

func (a *Alien) IsAlive() bool {
	return a.respawnTimer.IsReady() && !a.IsDying()
}

func (a *Alien) Bullets(callback func(bullet *Bullet)) {
	for _, bullet := range a.bullets {
		callback(bullet)
	}
}

func (a *Alien) IsDying() bool {
	return a.deadTimer != nil
}

func (a *Alien) SpinOutOfControl() {
	a.sprite.Orientation += 3 * math.Pi / float64(ebiten.TPS())
	a.deadTimer.Update()

	if a.deadTimer.IsReady() {
		a.respawnTimer.Reset()
		a.sprite.Reset()
		a.deadTimer = nil
	}
}
