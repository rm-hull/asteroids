package geometry

import "math"

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
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Accumulate(other *Vector) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector) Scale(factor float64) {
	v.X *= factor
	v.Y *= factor
}

func Add(a, b *Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
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
