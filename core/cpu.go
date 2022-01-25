package core

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/masu-mi/gochip-8/env"
)

// CHIP-8 emulator
// https://en.wikipedia.org/wiki/CHIP-8
type Chip8 struct {
	*Cpu
	*Memory
	Display
	*Keyboard
	Buzzer
}

const StartOfProgram = 0x200

func (chip *Chip8) Init(rom io.Reader) (int, error) {
	for i, f := range Font {
		copy(chip.Memory.Buf[fontAddr(uint8(i)):fontAddr(uint8(i))+5], f[0:])
	}
	return chip.Memory.Load(StartOfProgram, rom)
}

func (chip *Chip8) Run() {
	chip.Cpu.Run(chip.Memory, chip.Display, chip.Keyboard, chip.Buzzer)
}
func (chip *Chip8) Tick() {
	chip.Cpu.Tick(chip.Memory, chip.Display, chip.Keyboard, chip.Buzzer)
}

type Cpu struct {
	*rand.Rand

	V [16]uint8
	I uint16

	Dt *DelayedTimer
	St *DelayedTimer

	Pc    uint16
	Sp    uint8
	Stack [16]uint16
}

func NewCpu(buz Buzzer) *Cpu {
	c := &Cpu{
		Pc: StartOfProgram,
		Dt: NewDelayedTimer(60, nil),
		St: NewDelayedTimer(60, buz),
	}
	return c
}

func (cpu *Cpu) Run(ram *Memory, disp Display, keys *Keyboard, buz Buzzer) {
	for {
		if cpu.Pc >= uint16(len(ram.Buf)) {
			break
		}
		cpu.Tick(ram, disp, keys, buz)
	}
}

// Tick
func (cpu *Cpu) Tick(ram *Memory, disp Display, keys *Keyboard, buz Buzzer) {
	op := ram.Buf[cpu.Pc : cpu.Pc+2]
	inst := NewInstruction(op)
	switch inst.o1 {
	case 0x0:
		switch {
		case inst == instruction{0, 0, 0xe, 0}:
			trace("00E0 - CLS")
			disp.Clear()
		case inst == instruction{0, 0, 0xe, 0xe}:
			trace("00EE - RET")
			cpu.Pc = cpu.Stack[cpu.Sp-1]
			cpu.Sp--
		case inst.o1 == 0:
			next := addr(inst.o2, inst.o3, inst.o4)
			trace("0nnn - SYS 0x%X", next)
			cpu.Pc = next
			return
		}
	case 0x1:
		next := addr(inst.o2, inst.o3, inst.o4)
		trace("1nnn - JP 0x%x", next)
		cpu.Pc = next
		return
	case 0x2:
		next := addr(inst.o2, inst.o3, inst.o4)
		trace("2nnn - CALL addr 0x%x", next)
		cpu.Sp++
		cpu.Stack[cpu.Sp-1] = cpu.Pc
		cpu.Pc = next
		return
	case 0x3:
		kk := bite(inst.o3, inst.o4)
		cv := cpu.V[inst.o2]
		trace("3xkk - SE V%d(0x%x), 0x%x", inst.o2, cv, kk)
		if cv == kk {
			cpu.Pc += 2
		}
	case 0x4:
		kk := bite(inst.o3, inst.o4)
		cv := cpu.V[inst.o2]
		trace("4xkk - SNE V%d(0x%x), 0x%x", inst.o2, cv, kk)
		if cv != kk {
			cpu.Pc += 2
		}
	case 0x5:
		if inst.o4 != 0x0 {
			panic(fmt.Sprintf("N/A: `%v`", inst))
		}
		cx := cpu.V[inst.o2]
		cy := cpu.V[inst.o3]
		trace("5xy0 - SE V%d(0x%x), V%d(0x%x)", inst.o2, cx, inst.o3, cy)
		if cx == cy {
			cpu.Pc += 2
		}
	case 0x6:
		v := bite(inst.o3, inst.o4)
		trace("6xkk - LD V%d, 0x%x", inst.o2, v)
		cpu.V[inst.o2] = v
	case 0x7:
		v := bite(inst.o3, inst.o4)
		trace("7xkk - ADD V%d, 0x%x", inst.o2, v)
		cpu.V[inst.o2] += v
	case 0x8:
		switch inst.o4 {
		case 0x0:
			trace("8xy0 - LD V%d, V%d", inst.o2, inst.o3)
			cpu.V[inst.o2] = cpu.V[inst.o3]
		case 0x1:
			trace("8xy1 - OR V%d, V%d", inst.o2, inst.o3)
			cpu.V[inst.o2] |= cpu.V[inst.o3]
		case 0x2:
			trace("8xy2 - AND V%d, V%d", inst.o2, inst.o3)
			cpu.V[inst.o2] &= cpu.V[inst.o3]
		case 0x3:
			trace("8xy3 - XOR V%d, V%d", inst.o2, inst.o3)
			and := cpu.V[inst.o2] & cpu.V[inst.o3]
			cpu.V[inst.o2] = (cpu.V[inst.o2] | cpu.V[inst.o3]) & ^and
		case 0x4:
			trace("8xy4 - ADD V%d, V%d", inst.o2, inst.o3)
			add := uint16(cpu.V[inst.o2]) + uint16(cpu.V[inst.o3])
			if add > 0xff {
				add = add & 0xff
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[inst.o2] = uint8(add)
		case 0x5:
			trace("8xy5 - SUB V%d, V%d", inst.o2, inst.o3)
			vx := cpu.V[inst.o2]
			vy := cpu.V[inst.o3]
			if vx > vy {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[inst.o2] = vx - vy
		case 0x6:
			trace("8xy6 - SHR V%d {, V%d}", inst.o2, inst.o3)
			cpu.V[0xF] = cpu.V[inst.o2] & 0x1
			cpu.V[inst.o2] >>= 1
		case 0x7:
			trace("8xy7 - SUBN V%d, V%d", inst.o2, inst.o3)
			vx := cpu.V[inst.o2]
			vy := cpu.V[inst.o3]
			if vx < vy {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[inst.o2] = vy - vx
		case 0xE:
			trace("8xyE - SHL V%d {, V%d}", inst.o2, inst.o3)
			cpu.V[0xF] = cpu.V[inst.o2] >> 7 & 0x1
			cpu.V[inst.o2] <<= 1
		}
	case 0x9:
		if inst.o4 != 0x0 {
			panic(fmt.Sprintf("N/A: `%v`", inst))
		}
		vx, vy := cpu.V[inst.o2], cpu.V[inst.o3]
		trace("9xy0 - SNE V%d(0x%x), V%d(0x%x)", inst.o2, vx, inst.o3, vy)
		if vx != vy {
			cpu.Pc += 2
		}
	case 0xA:
		p := addr(inst.o2, inst.o3, inst.o4)
		trace("Annn - LD I, *(0x%x)", p)
		cpu.I = p
	case 0xB:
		p := addr(inst.o2, inst.o3, inst.o4)
		trace("Bnnn - JP V0, *(0x%x)", p)
		cpu.Pc = p + uint16(cpu.V[0x0])
		return
	case 0xC:
		v := bite(inst.o3, inst.o4)
		trace("Cxkk - RND V%d, 0x%x", inst.o2, v)
		r := uint8(rand.Intn(256))
		if cpu.Rand != nil {
			r = uint8(cpu.Rand.Intn(256))
		}
		cpu.V[inst.o2] = r & v
	case 0xD:
		trace("Dxyn - DRW V%d, V%d, %d[byte]", inst.o2, inst.o3, inst.o4)
		sprite := ram.Buf[cpu.I : cpu.I+uint16(inst.o4)]
		if disp.Draw(cpu.V[inst.o2], cpu.V[inst.o3], sprite) {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
	case 0xE:
		// TODO
		if inst.o3 == 0x9 && inst.o4 == 0xE {
			trace("Ex9E - SKP V%d", inst.o2)
			pressed := keys.PressedKeys()
			target := cpu.V[inst.o2]
			for _, v := range pressed {
				if v == target {
					cpu.Pc += 2
					break
				}
			}
		} else if inst.o3 == 0xA && inst.o4 == 0x1 {
			trace("ExA1 - SKNP V%d", inst.o2)
			cpu.Pc += 2 // skip is default
			pressed := keys.PressedKeys()
			target := cpu.V[inst.o2]
			for _, v := range pressed {
				if v == target {
					cpu.Pc -= 2
					break
				}
			}
		} else {
			panic(fmt.Sprintf("N/A: `%v`", inst))
		}
	case 0xF:
		// TODO
		switch {
		case inst.o3 == 0x0 && inst.o4 == 0x7:
			trace("Fx07 - LD V%d, DT", inst.o2)
			cpu.V[inst.o2] = cpu.Dt.GetV()
		case inst.o3 == 0x0 && inst.o4 == 0xA:
			trace("Fx0A - LD V%d, K", inst.o2)
			// TODO more less CPU consuption
		BLOCK:
			for {
				pressed := keys.PressedKeys()
				for _, key := range pressed {
					if key == cpu.V[inst.o2] {
						break BLOCK
					}
				}
			}
			// Fx0A - LD Vx, K; Wait for a key press, store the value of the key in Vx.
		case inst.o3 == 0x1 && inst.o4 == 0x5:
			trace("Fx15 - LD DT, V%d", inst.o2)
			cpu.Dt.SetV(cpu.V[inst.o2])
		case inst.o3 == 0x1 && inst.o4 == 0x8:
			trace("Fx18 - LD ST, V%d", inst.o2)
			cpu.St.SetV(cpu.V[inst.o2])
		case inst.o3 == 0x1 && inst.o4 == 0xE:
			trace("Fx1E - ADD I, V%d", inst.o2)
			cpu.I += uint16(cpu.V[inst.o2])
		case inst.o3 == 0x2 && inst.o4 == 0x9:
			trace("Fx29 - LD F, V%d", inst.o2)
			cpu.I = fontAddr(inst.o2)
		case inst.o3 == 0x3 && inst.o4 == 0x3:
			trace("Fx33 - LD B, V%d", inst.o2)
			ram.Buf[cpu.I], ram.Buf[cpu.I+1], ram.Buf[cpu.I+2] = bcd(cpu.V[inst.o2])
		case inst.o3 == 0x5 && inst.o4 == 0x5:
			trace("Fx55 - LD [I], V%d", inst.o2)
			copy(ram.Buf[cpu.I:(cpu.I+uint16(inst.o2)+1)], cpu.V[0:inst.o2+1])
		case inst.o3 == 0x6 && inst.o4 == 0x5:
			trace("Fx65 - LD V%d, [I]", inst.o2)
			copy(cpu.V[0:inst.o2+1], ram.Buf[cpu.I:(cpu.I+uint16(inst.o2)+1)])
		}
	}
	// All instructions are 2 bytes long and are stored most-significant-byte first.
	// In memory, the first byte of each instruction should be located at an even addresses.
	// If a program includes sprite data, it should be padded so any instructions following it will be properly situated in RAM.
	cpu.Pc += 2
}

func bcd(v uint8) (h, t, o uint8) {
	o = v % 10
	v /= 10
	t = v % 10
	h = v / 10
	return
}
func addr(n1, n2, n3 uint8) uint16 {
	return uint16(n1)<<8 + uint16(n2)<<4 + uint16(n3)
}
func bite(n1, n2 uint8) uint8 {
	return n1<<4 + n2
}

func trace(msg string, d ...interface{}) {
	if env.DEBUG {
		fmt.Printf(fmt.Sprintf("%s\n", msg), d...)
	}
}

type instruction struct {
	o1, o2, o3, o4 uint8
}

func NewInstruction(seg []byte) instruction {
	return instruction{
		seg[0] >> 4,
		seg[0] & 0b00001111,
		seg[1] >> 4,
		seg[1] & 0b00001111,
	}
}
