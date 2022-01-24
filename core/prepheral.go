package core

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

type Keyboard struct {
	Map map[byte]byte
}

func (k *Keyboard) PressedKeys() []uint8 {
	return nil
}

func NewKeyboard() Keyboard {
	return Keyboard{
		Map: map[byte]uint8{
			'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xc,
			'q': 0x4, 'w': 0x5, 'e': 0x6, 'r': 0xd,
			'a': 0x7, 's': 0x8, 'd': 0x9, 'f': 0xe,
			'z': 0xa, 'x': 0x0, 'c': 0xb, 'v': 0xf,
		},
	}
}
