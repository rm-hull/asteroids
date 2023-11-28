package main

import (
	"asteroids/internal"
	"asteroids/internal/entity"
	"asteroids/internal/geometry"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Player     *entity.Player
	Alien      *entity.Alien
	Asteroids  map[int]*entity.Asteroid
	Sequence   *internal.Sequence
	fullscreen bool
	paused     bool
	level      int
}

var screenSize = geometry.Dimension{W: 1024, H: 768}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return errors.New("dejar de ser un desertor")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset(10)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
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

	g.HandleCollisionDetection()

	err = g.Alien.Update()
	if err != nil {
		return err
	}

	if len(g.Asteroids) == 0 {
		g.NextLevel()
	}

	return nil
}

func (g *Game) HandleCollisionDetection() {
	g.Player.Bullets(func(bullet *entity.Bullet) {
		for _, asteroid := range g.Asteroids {
			if !asteroid.IsExploded() && bullet.CollisionDetected(asteroid) {
				for _, fragment := range asteroid.Explode() {
					g.Asteroids[g.Sequence.GetNext()] = fragment
				}
			}
		}

		if g.Alien.IsAlive() && bullet.CollisionDetected(g.Alien) {
			g.Alien.Kill()
			g.Player.UpdateScore(g.Alien.Value())
		}
	})

	g.Alien.Bullets(func(bullet *entity.Bullet) {
		if g.Player.IsAlive() && bullet.CollisionDetected(g.Player) {
			g.Player.Kill()
		}
	})

	for _, asteroid := range g.Asteroids {
		if g.Player.IsAlive() && asteroid.CollisionDetected(g.Player) {
			g.Player.Kill()
		}
	}

	if g.Player.IsAlive() && g.Alien.IsAlive() && g.Alien.Bounds().Overlaps(*g.Player.Bounds()) {
		g.Player.Kill()
	}

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

func (g *Game) Reset(n int) {
	g.level = 1
	g.Player = entity.NewPlayer(&screenSize)
	g.Alien = entity.NewAlien(g.level, g.Player.NotNear(), g.Player.CurrentPosition(), &screenSize)
	g.Asteroids = entity.NewAsteroidBelt(n, g.Sequence, g.Player, &screenSize)
}

func (g *Game) NextLevel() {
	g.level++
	g.Alien = entity.NewAlien(g.level, g.Player.NotNear(), g.Player.CurrentPosition(), &screenSize)
	g.Asteroids = entity.NewAsteroidBelt(5+g.level, g.Sequence, g.Player, &screenSize)
}

func main() {
	player := entity.NewPlayer(&screenSize)
	seq := internal.NewSequence()
	g := &Game{
		Sequence:   seq,
		Player:     player,
		Alien:      entity.NewAlien(1, player.NotNear(), player.CurrentPosition(), &screenSize),
		Asteroids:  entity.NewAsteroidBelt(5, seq, player, &screenSize),
		fullscreen: false,
		level:      1,
	}

	// ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(screenSize.W), int(screenSize.H))
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
