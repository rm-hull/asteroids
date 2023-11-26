package main

import (
	"asteroids/internal"
	"asteroids/internal/entity"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player        *entity.Player
	Asteroids     []*entity.Asteroid
	Bullets       map[time.Time]*entity.Bullet
	fullscreen    bool
	shootCooldown *internal.Timer
}

var screenSize = geometry.Dimension{W: 1024, H: 768}
var cooldownTime = 100 * time.Millisecond
var maxSalvo = 3

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("dejar de ser un desertor")
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	g.shootCooldown.Update()
	if g.shootCooldown.IsReady() && len(g.Bullets) < maxSalvo && ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.shootCooldown.Reset()
		g.Bullets[time.Now()] = g.Player.FireBullet()
	}

	for _, asteroid := range g.Asteroids {
		err := asteroid.Update()
		if err != nil {
			return err
		}
	}

	for idx, bullet := range g.Bullets {
		err := bullet.Update()
		if err != nil {
			return err
		}

		if bullet.IsExpired() {
			delete(g.Bullets, idx)
		}
	}

	err := g.Player.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, asteroid := range g.Asteroids {
		asteroid.Draw(screen)
	}

	for _, bullet := range g.Bullets {
		bullet.Draw(screen)
	}

	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(screenSize.W), int(screenSize.H)
}

func main() {
	player := entity.NewPlayer(&screenSize)
	posn := player.CurrentPosition()
	g := &Game{
		shootCooldown: internal.NewTimer(cooldownTime),
		Player:        player,
		Bullets:       make(map[time.Time]*entity.Bullet),
		Asteroids: []*entity.Asteroid{
			entity.NewAsteroid(sprites.Large, posn, &screenSize),
			entity.NewAsteroid(sprites.Large, posn, &screenSize),
			entity.NewAsteroid(sprites.Large, posn, &screenSize),
			entity.NewAsteroid(sprites.Medium, posn, &screenSize),
			entity.NewAsteroid(sprites.Medium, posn, &screenSize),
			entity.NewAsteroid(sprites.Small, posn, &screenSize),
			entity.NewAsteroid(sprites.Small, posn, &screenSize),
			entity.NewAsteroid(sprites.Small, posn, &screenSize),
		},
		fullscreen: false,
	}

	// ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(screenSize.W), int(screenSize.H))
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
