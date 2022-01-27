//go:build debug

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/masu-mi/gochip-8/core"
	"github.com/mattn/go-tty"
)

func main() {
	path := flag.String("rom", "", "rom file path")
	flag.Parse()

	tty, _ := tty.Open()
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("panic")
		}
		tty.Close()
	}()
	forKeys := make(chan rune)
	forRepl := make(chan rune)
	go func() {
		defer close(forKeys)
		defer close(forRepl)
		for {
			r, e := tty.ReadRune()
			if e != nil {
				break
			}
			forKeys <- r
			forRepl <- r
		}
	}()

	chip := &core.Chip8{
		Cpu:      core.NewCpu(nil),
		Memory:   &core.Memory{},
		Display:  &Ignore{},
		Keyboard: NewKeyboard(forKeys, DefaultConvert),
	}
	f, e := os.Open(*path)
	if e != nil {
		log.Fatalf("can't open `%s`\n", *path)
	}
	defer f.Close()
	n, e := chip.Init(f)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Printf("load: %d[byte]\n", n)

	chip.Memory.WriteTo(os.Stdout)
	// fmt.Printf("%v\n", chip.Memory.Buf)
	num := 0
	for range forRepl {
		fmt.Printf("tick(%d): Pc: %04x(%d)", num, chip.Cpu.Pc, chip.Cpu.Pc)
		chip.Tick()
		num++
	}
}

type Ignore struct{}

func (i *Ignore) Clear() {}
func (i *Ignore) Draw(x, y uint8, sprite []byte) (collision bool) {
	return false
}

var _ core.Display = &Ignore{}

type Keyboard struct {
	sync.RWMutex
	time.Duration

	tty    <-chan rune
	events chan uint8

	convert map[rune]uint8
	pressed map[uint8]bool
	timers  map[uint8]*time.Timer
}

var _ core.Keyboard = &Keyboard{}

var DefaultConvert = map[rune]uint8{
	'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xc,
	'q': 0x4, 'w': 0x5, 'e': 0x6, 'r': 0xd,
	'a': 0x7, 's': 0x8, 'd': 0x9, 'f': 0xe,
	'z': 0xa, 'x': 0x0, 'c': 0xb, 'v': 0xf,
}

func NewKeyboard(tty <-chan rune, convert map[rune]uint8) *Keyboard {
	dev := &Keyboard{
		tty:      tty,
		Duration: time.Second / time.Duration(60),
		pressed:  map[uint8]bool{},
		convert:  map[rune]uint8{},
		events:   make(chan uint8),
	}

	go func() {
		for {
			r := <-dev.tty
			k, ok := dev.convert[r]
			if !ok {
				continue
			}
			dev.press(k)
		}
	}()
	return dev
}

func (k *Keyboard) up(key uint8) {
	k.Lock()
	defer k.Unlock()
	k.pressed[key] = false
	delete(k.timers, key)
}

func (k *Keyboard) press(key uint8) {
	k.RWMutex.Lock()
	defer k.RWMutex.RUnlock()

	k.pressed[key] = true
	t, ok := k.timers[key]
	if !ok {
		t = time.AfterFunc(k.Duration, func() { k.up(key) })
	}
	t.Reset(k.Duration)
	k.timers[key] = t
	select {
	case k.events <- key:
	default:
	}
}

func (k *Keyboard) IsPressed(key uint8) bool {
	k.RLock()
	defer k.RUnlock()
	return k.pressed[key]
}
func (k *Keyboard) Wait(key uint8) {
	ch := make(chan struct{})
	go func() {
		for {
			pressed := <-k.events
			if pressed == key {
				break
			}
		}
		defer close(ch)
	}()
	<-ch
	return
}
