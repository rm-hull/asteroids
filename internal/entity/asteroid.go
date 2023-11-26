package entity

import (
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"image"

	// "image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const asteroidMaxSpeed = 2

type Asteroid struct {
	size      int
	position  geometry.Vector
	velocity  geometry.Vector
	rotation  float64
	direction float64
	bounds    *geometry.Dimension
	sprite    *ebiten.Image
	exploded  bool
}

func NewAsteroid(size int, position *geometry.Vector, screenBounds *geometry.Dimension) *Asteroid {

	direction := rand.Float64() * 2 * math.Pi
	speed := rand.Float64() * asteroidMaxSpeed

	return &Asteroid{
		size:     size,
		position: *position,
		velocity: geometry.VectorFrom(direction, speed),
		rotation: (rand.Float64() - 0.5) / 20,
		bounds:   screenBounds,
		sprite:   sprites.Asteroid(size),
	}
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	bounds := a.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(a.direction)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(a.position.X, a.position.Y)

	// bounds := a.Bounds()
	// ebitenutil.DrawRect(screen, float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Dx()), float64(bounds.Dy()), color.RGBA{67, 0, 255, 255})
	// ebitenutil.DrawCircle(screen, a.position.X+halfW, a.position.Y+halfH, halfH*0.75, color.RGBA{0, 128, 255, 128})

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

func (a *Asteroid) Bounds() *image.Rectangle {
	point := image.Point{X: int(a.position.X), Y: int(a.position.Y)}
	return &image.Rectangle{
		Min: point,
		Max: a.sprite.Bounds().Max.Add(point),
	}
}
