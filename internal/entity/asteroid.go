package entity

import (
	"github.com/rm-hull/asteroids/internal"
	"github.com/rm-hull/asteroids/internal/geometry"
	"github.com/rm-hull/asteroids/internal/sprites"
	"github.com/rm-hull/asteroids/resources/soundfx"

	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const asteroidMaxSpeed = 2

type Asteroid struct {
	sprite       *sprites.Sprite
	size         int
	exploded     bool
	screenBounds *geometry.Dimension
}

func randSize() int {
	n := rand.Intn(10)
	if n < 5 {
		return sprites.Large
	}

	if n < 8 {
		return sprites.Medium
	}

	return sprites.Small
}

func NewAsteroidBelt(n int, seq *internal.Sequence, player *Player, screenBounds *geometry.Dimension) map[int]*Asteroid {
	var asteroids = make(map[int]*Asteroid)
	for i := 0; i < n; i++ {
		idx := seq.GetNext()
		asteroids[idx] = NewAsteroid(randSize(), player.NotNear(), screenBounds)
	}
	return asteroids
}

func NewAsteroid(size int, position *geometry.Vector, screenBounds *geometry.Dimension) *Asteroid {

	sprite := sprites.NewSprite(screenBounds, sprites.Asteroid(size), true)
	sprite.Speed = (rand.Float64() + 0.3) * asteroidMaxSpeed
	sprite.Direction = rand.Float64() * 2 * math.Pi
	sprite.Position.X = position.X
	sprite.Position.Y = position.Y
	sprite.Velocity = geometry.VectorFrom(sprite.Direction, sprite.Speed)
	sprite.Rotation = (rand.Float64() - 0.5) / 20

	return &Asteroid{
		sprite:       sprite,
		size:         size,
		exploded:     false,
		screenBounds: screenBounds,
	}
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	if !a.exploded {
		a.sprite.Draw(screen)
	}
}

func (a *Asteroid) Update() error {
	if !a.exploded {
		if err := a.sprite.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Asteroid) Explode() []*Asteroid {
	a.exploded = true
	sePlayer := audioContext.NewPlayerFromBytes(soundfx.Explosion2)
	sePlayer.SetVolume(0.15)
	sePlayer.Play()

	arr := make([]*Asteroid, 0)
	switch a.size {
	case sprites.Large:
		n := rand.Intn(3) + 1
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Medium, a.sprite.Position, a.screenBounds))
		}
		n = rand.Intn(5 - n)
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Small, a.sprite.Position, a.screenBounds))
		}
	case sprites.Medium:
		n := rand.Intn(2) + 2
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Small, a.sprite.Position, a.screenBounds))
		}
	default:
		break
	}
	return arr
}

func (a *Asteroid) IsExploded() bool {
	return a.exploded
}

func (a *Asteroid) Value() int {
	switch a.size {
	case sprites.Large:
		return 20
	case sprites.Medium:
		return 50
	case sprites.Small:
		return 100
	default:
		return 0
	}
}

func (a *Asteroid) Size() float64 {
	return a.sprite.Centre.Y * 0.70
}

func (a *Asteroid) Position() *geometry.Vector {
	return geometry.Add(a.sprite.Position, a.sprite.Centre).Mod(a.screenBounds)
}
