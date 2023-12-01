package entity

import (
	"asteroids/internal"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"

	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const asteroidMaxSpeed = 2

type Asteroid struct {
	size      int
	position  geometry.Vector
	velocity  geometry.Vector
	centre    geometry.Vector
	rotation  float64
	direction float64
	bounds    *geometry.Dimension
	sprite    *ebiten.Image
	exploded  bool
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

	direction := rand.Float64() * 2 * math.Pi
	speed := rand.Float64() * asteroidMaxSpeed
	sprite := sprites.Asteroid(size)

	return &Asteroid{
		size:     size,
		position: *position,
		centre:   sprites.Centre(sprite),
		velocity: geometry.VectorFrom(direction, speed),
		rotation: (rand.Float64() - 0.5) / 20,
		bounds:   screenBounds,
		sprite:   sprite,
	}
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-a.centre.X, -a.centre.Y)
	op.GeoM.Rotate(a.direction)
	op.GeoM.Translate(a.centre.X, a.centre.Y)

	op.GeoM.Translate(a.position.X, a.position.Y)

	// vector.DrawFilledCircle(screen, float32(a.position.X+a.centre.X), float32(a.position.Y+a.centre.Y), float32(a.Size()), color.RGBA{0, 128, 255, 128}, false)

	screen.DrawImage(a.sprite, op)
}

func (a *Asteroid) Update() error {
	a.direction += a.rotation
	a.position.Add(&a.velocity)
	a.position.CheckEdges(a.bounds, float64(a.sprite.Bounds().Dx()), float64(a.sprite.Bounds().Dy()))
	return nil
}

func (a *Asteroid) Explode() []*Asteroid {
	a.exploded = true
	arr := make([]*Asteroid, 0)
	switch a.size {
	case sprites.Large:
		n := rand.Intn(2) + 1
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Medium, &a.position, a.bounds))
		}
		n = rand.Intn(4 - n)
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Small, &a.position, a.bounds))
		}
	case sprites.Medium:
		n := rand.Intn(2) + 2
		for i := 0; i < n; i++ {
			arr = append(arr, NewAsteroid(sprites.Small, &a.position, a.bounds))
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
		return 10
	case sprites.Medium:
		return 25
	case sprites.Small:
		return 50
	default:
		return 0
	}
}

func (a *Asteroid) Size() float64 {
	return a.centre.Y * 0.70
}

func (a *Asteroid) Position() *geometry.Vector {
	return geometry.Add(&a.position, &a.centre)
}
