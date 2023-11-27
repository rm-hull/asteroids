package main

import (
	"asteroids/internal"
	"asteroids/internal/entity"
	"asteroids/internal/geometry"
	"asteroids/internal/sprites"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player     *entity.Player
	Alien      *entity.Alien
	Asteroids  map[int]*entity.Asteroid
	Sequence   *internal.Sequence
	fullscreen bool
	paused     bool
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

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

	if g.paused {
		return nil
	}

	for idx, asteroid := range g.Asteroids {
		err := asteroid.Update()
		if err != nil {
			return err
		}

		if asteroid.IsExploded() {
			g.Player.UpdateScore(asteroid.Value())
			delete(g.Asteroids, idx)
		}
	}
	
	err := g.Player.Update()
	if err != nil {
		return err
	}
	g.Player.Bullets(func(bullet *entity.Bullet) {
		for _, fragment := range bullet.AsteroidCollisionDetection(g.Asteroids) {
			g.Asteroids[g.Sequence.GetNext()] = fragment
		}

		if bullet.AlienCollisionDetection(g.Alien) {
			g.Player.UpdateScore(g.Alien.Value())
		}
	})

	err = g.Alien.Update()
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
	g.Alien.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(screenSize.W), int(screenSize.H)
}

func main() {
	player := entity.NewPlayer(&screenSize)
	seq := internal.NewSequence()
	g := &Game{
		Sequence: seq,
		Player:   player,
		Alien:    entity.NewAlien(&screenSize, player.NotNear()),
		Asteroids: map[int]*entity.Asteroid{
			seq.GetNext(): entity.NewAsteroid(sprites.Large, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Large, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Large, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Large, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Large, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Medium, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Medium, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Small, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Small, player.NotNear(), &screenSize),
			seq.GetNext(): entity.NewAsteroid(sprites.Small, player.NotNear(), &screenSize),
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
