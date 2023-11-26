package main

import (
	"asteroids/internal/entity"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player     *entity.Player
	Asteroids  []*entity.Asteroid
	fullscreen bool
}

var screenSize = geometry.Dimension{W: 1024, H: 768}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("dejar de ser un desertor")
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	for _, asteroid := range g.Asteroids {
		err := asteroid.Update()
		if err != nil {
			return err
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
	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(screenSize.W), int(screenSize.H)
}

func main() {
	player := entity.NewPlayer(screenSize)
	posn := player.Position()
	g := &Game{
		Player: player,
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
