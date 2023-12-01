package geometry

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
}

var Zero = Vector{X: 0, Y: 0}

func (v *Vector) Normalize() Vector {
	magnitude := v.Magnitude()
	return Vector{v.X / magnitude, v.Y / magnitude}
}

func (v *Vector) Magnitude() float64 {
	return v.DistanceFrom(&Zero)
}

func (v Vector) String() string {
	return fmt.Sprintf("%0.1f,%0.1f", v.X, v.Y)
}

func (v *Vector) AngleTo(other *Vector) float64 {
	dx := other.X - v.X
	dy := other.Y - v.Y
	signedAngle := math.Atan2(dy, dx)
	if signedAngle < 0 {
		signedAngle += 2 * math.Pi
	}
	return signedAngle
}

func (v *Vector) SquareDistanceFrom(other *Vector) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return (dx*dx + dy*dy)
}

func (v *Vector) DistanceFrom(other *Vector) float64 {
	return math.Sqrt(v.SquareDistanceFrom(other))
}

func Add(a, b *Vector) *Vector {
	return &Vector{X: a.X + b.X, Y: a.Y + b.Y}
}

func (v *Vector) Add(other *Vector) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector) Scale(factor float64) {
	v.X *= factor
	v.Y *= factor
}

func VectorFrom(direction float64, speed float64) Vector {
	return Vector{
		X: speed * math.Cos(direction),
		Y: speed * math.Sin(direction),
	}
}

func (v *Vector) CheckEdges(screenSize *Dimension, w, h float64) {
	if v.X+w < 0 {
		v.X = screenSize.W
	} else if v.X-w > screenSize.W {
		v.X = -w
	}
	if v.Y+h < 0 {
		v.Y = screenSize.H
	} else if v.Y-h > screenSize.H {
		v.Y = -h
	}
}
