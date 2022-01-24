package core

import (
	"fmt"
	"io"
)

// Memory is Chip-8's RAM.
//
// > The Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM, from location 0x000 (0) to 0xFFF (4095).
// > ref. http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2.1
type Memory struct {
	Buf [(0x1000)]uint8
}

// The first 512 bytes, from 0x000 to 0x1FF, are where the original interpreter was located, and should not be used by programs.
func (m *Memory) Load(start uint16, rom io.Reader) (int, error) {
	n, e := rom.Read(m.Buf[start:])
	// _, e := io.Copy(bytes.NewBuffer(m.Buf[start:start]), rom)
	if e != nil {
		// }&& e != io.EOF {
		return n, e
	}
	return n, nil
}
func (m *Memory) WriteTo(dst io.Writer) (int64, error) {
	var num int64
	for i := StartOfProgram; i < len(m.Buf); i += 2 {
		s, e := dst.Write([]byte(fmt.Sprintf("%v\n", NewInstruction(m.Buf[i:i+2]))))
		num += int64(s)
		if e != nil {
			return num, e
		}
	}
	return num, nil
}
