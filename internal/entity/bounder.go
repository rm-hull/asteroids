package entity

import "image"

type Bounder interface {
	Bounds() *image.Rectangle
}
