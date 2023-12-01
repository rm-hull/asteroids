package entity

import "asteroids/internal/geometry"

type Collider interface {
	Position() *geometry.Vector
	Size() float64
}

func CollisionDetected(a Collider, b Collider) bool {
	actaulSquareDist := a.Position().SquareDistanceFrom(b.Position())
	minDist := a.Size() + b.Size()
	hit := actaulSquareDist < minDist*minDist
	return hit
}
