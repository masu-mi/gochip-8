package main

import (
	"sync"
	"time"

	"github.com/masu-mi/gochip-8/core"
)

type Ignore struct{}

func (i *Ignore) Clear() {}
func (i *Ignore) Draw(x, y uint8, sprite []byte) (collision bool) {
	return false
}

var _ core.Display = &Ignore{}

type Keyboard struct {
	sync.RWMutex
	tty <-chan rune
	time.Duration
	convert map[rune]uint8

	events  chan uint8
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
		convert:  convert,

		events:  make(chan uint8),
		pressed: map[uint8]bool{},
		timers:  map[uint8]*time.Timer{},
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
	defer k.RWMutex.Unlock()

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
