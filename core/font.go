package core

func fontAddr(i uint8) uint16 {
	return uint16(i * 5)
}

var Font = [(0x10)][5]byte{
	{
		0b11110000,
		0b10010000,
		0b10010000,
		0b10010000,
		0b11110000,
	},
	{
		0b00100000,
		0b01100000,
		0b00100000,
		0b00100000,
		0b01110000,
	},
	{
		0b11110000,
		0b00010000,
		0b11110000,
		0b10000000,
		0b11110000,
	},
	{
		0b11110000,
		0b00010000,
		0b11110000,
		0b00010000,
		0b11110000,
	},
	{
		0b10010000,
		0b10010000,
		0b11110000,
		0b00010000,
		0b00010000,
	},
	{
		0b11110000,
		0b10000000,
		0b11110000,
		0b00010000,
		0b11110000,
	},
	{
		0b11110000,
		0b10000000,
		0b11110000,
		0b10010000,
		0b11110000,
	},
	{
		0b11110000,
		0b00010000,
		0b00100000,
		0b10000000,
		0b10000000,
	},
	{
		0b11110000,
		0b10010000,
		0b11110000,
		0b10010000,
		0b11110000,
	},
	{
		0b11110000,
		0b10010000,
		0b11110000,
		0b00010000,
		0b11110000,
	},
	{
		0b11110000,
		0b10010000,
		0b11110000,
		0b10010000,
		0b10010000,
	},
	{
		0b11100000,
		0b10010000,
		0b11100000,
		0b10010000,
		0b11100000,
	},
	{
		0b11110000,
		0b10000000,
		0b10000000,
		0b10000000,
		0b11110000,
	},
	{
		0b11100000,
		0b10010000,
		0b10010000,
		0b10010000,
		0b11100000,
	},
	{
		0b11110000,
		0b10000000,
		0b11110000,
		0b10000000,
		0b11110000,
	},
	{
		0b11110000,
		0b10000000,
		0b11110000,
		0b10000000,
		0b10000000,
	},
}