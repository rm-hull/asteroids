package soundfx

import (
	_ "embed"
)

//go:embed laser-gun-shot1.wav
var LazerGunShot1 []byte

//go:embed laser-gun-shot2.wav
var LazerGunShot2 []byte

//go:embed explosion1.wav
var Explosion1 []byte

//go:embed explosion1.wav
var Explosion2 []byte
