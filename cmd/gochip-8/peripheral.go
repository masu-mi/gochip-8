package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/masu-mi/gochip-8/core"
	"github.com/nsf/termbox-go"
)

type Display struct {
	color termbox.Attribute
}

func StarTermbox(ctx context.Context, color termbox.Attribute) (context.Context, *Display, *Keyboard, error) {
	c, cancel := context.WithCancel(ctx)
	e := termbox.Init()
	if e != nil {
		termbox.Close()
		cancel()
		return c, nil, nil, e
	}
	w, h := termbox.Size()
	if !IsDisplaySizeSufficient(w, h) {
		termbox.Close()
		cancel()
		return c, nil, nil, errors.New("terminal is too small")
	}
	ch := make(chan rune)
	kb := NewKeyboard(ch, DefaultConvert)
	go func() {
	MAINLOOP:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					termbox.Close()
					cancel()
					break MAINLOOP
				default:
					ch <- ev.Ch
				}
			}
		}
	}()
	return c, &Display{color: color}, kb, nil
}

func IsDisplaySizeSufficient(w, h int) bool {
	return w >= 64 && h >= 32
}

func (t *Display) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}
func (t *Display) Draw(x, y uint8, sprite []byte) (collision bool) {
	for dh, b := range sprite {
		for rdw := 0; rdw < 8; rdw++ {
			input := (b >> rdw) & 1
			cx, cy := (int(x)+7-rdw)%core.WIDTH, (int(y)+dh)%core.HEIGHT
			cell := termbox.GetCell(cx, cy)
			col := cell.Bg == t.color && input == 1
			if (cell.Bg == t.color || input == 1) && !col {
				termbox.SetBg(cx, cy, t.color)
			} else {
				termbox.SetBg(cx, cy, termbox.ColorDefault)
			}
			collision = collision || col
		}
	}
	termbox.Flush()
	return collision
}

var _ core.Display = &Display{}

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
		Duration: time.Second / time.Duration(fps),
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
	k.Lock()
	defer k.Unlock()

	k.pressed[key] = true
	t, ok := k.timers[key]
	if !ok {
		t = time.AfterFunc(k.Duration, func() { k.up(key) })
	} else {
		t.Reset(k.Duration)
	}
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
func (k *Keyboard) Wait(ctx context.Context, key uint8) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case pressed := <-k.events:
			if pressed == key {
				break LOOP
			}
		}
	}
	return
}
