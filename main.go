package main

import (
	"asteroids/internal/sprites"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(sprites.BigAsteroid1, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(150, 200)
	screen.DrawImage(sprites.BigAsteroid2, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(190, 20)
	screen.DrawImage(sprites.BigAsteroid3, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(20, 150)
	screen.DrawImage(sprites.MediumAsteroid1, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(300, 180)
	screen.DrawImage(sprites.MediumAsteroid2, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 140)
	screen.DrawImage(sprites.MediumAsteroid3, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(400, 400)
	screen.DrawImage(sprites.SmallAsteroid1, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(300, 250)
	screen.DrawImage(sprites.SmallAsteroid2, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(100, 400)
	screen.DrawImage(sprites.SmallAsteroid3, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(500, 350)
	screen.DrawImage(sprites.SpaceShip1, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(450, 250)
	screen.DrawImage(sprites.SpaceShip2, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(400, 150)
	screen.DrawImage(sprites.AlienSpaceShip, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(350, 15)
	screen.DrawImage(sprites.Bullet1, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(15, 450)
	screen.DrawImage(sprites.Bullet2, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
