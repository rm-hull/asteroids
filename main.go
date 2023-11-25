package main

import (
	"asteroids/internal/entity"
	"asteroids/internal/geometry"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player     entity.Player
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

	err := g.Player.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(screenSize.W), int(screenSize.H)
}

func main() {
	g := &Game{
		Player:     *entity.NewPlayer(screenSize),
		fullscreen: false,
	}

	// ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(screenSize.W), int(screenSize.H))
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
