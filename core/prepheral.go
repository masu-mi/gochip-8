package core

import "context"

// Buzzer is called by SoundTimer
type Buzzer interface {
	Start()
	Stop()
}

// Display
//
// > The original implementation of the Chip-8 language used a 64x32-pixel monochrome display with this format:
// > http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2.4
type Display interface {
	Clear()
	Draw(x, y uint8, sprite []byte) (collision bool)
}

const (
	WIDTH  = 64
	HEIGHT = 32
)

type Keyboard interface {
	IsPressed(key uint8) bool
	Wait(ctx context.Context, key uint8)
}
