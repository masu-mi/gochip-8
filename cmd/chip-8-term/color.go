package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
)

func NewColorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "color",
		Short: "show color chart",
		RunE:  color,
	}
}

func color(_ *cobra.Command, args []string) error {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)
	pollEvent()
	return nil
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawLine(0, 0, "Press ESC to exit")
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			drawLineColor(j*4, i+3, fmt.Sprintf("%3d", i*16+j), termbox.Attribute(255-(i*16+j)), termbox.Attribute(i*16+j))
		}
	}
	termbox.Flush()
}

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for i, v := range runes {
		termbox.SetCell(x+i, y, v, color, backgroundColor)
	}
}

func drawLineColor(x, y int, str string, fg, bg termbox.Attribute) {
	runes := []rune(str)
	for i, v := range runes {
		termbox.SetCell(x+i, y, v, fg, bg)
	}
}

func key(ev termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			return false
		default:
			draw()
		}
	default:
		draw()
	}
	return true
}

func pollEvent() {
	evc := make(chan termbox.Event)
	go func() {
		for {
			evc <- termbox.PollEvent()
		}
	}()
	draw()
	for {
		select {
		case ev := <-evc:
			if !key(ev) {
				return
			}
		case <-time.After(1 * time.Second):
			draw()
		}
	}
}
