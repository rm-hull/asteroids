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
	return v.DistanceFrom(&Zero)
}

func (v *Vector) DistanceFrom(other *Vector) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
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
