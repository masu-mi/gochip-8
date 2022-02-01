package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/masu-mi/gochip-8/core"
	"github.com/mattn/go-tty"
	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
)

var (
	cpuHz      int
	fps        uint8
	path       string
	blockColor int64
)

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start CHIP-8 emulator",
		RunE:  start,
	}
	cmd.PersistentFlags().IntVar(&cpuHz, "cpu-hz", 1000000, "reciprocal of duration of key pressed (default: 1MHz)")
	cmd.PersistentFlags().Uint8Var(&fps, "keyboard-hz", 10, "reciprocal of duration of key pressed (default: 10Hz)")
	cmd.PersistentFlags().StringVar(&path, "rom", "", "rom image file path")
	cmd.PersistentFlags().Int64Var(&blockColor, "color", 16, "display active cell's color(defalt: 16)")
	return cmd
}

func start(_ *cobra.Command, args []string) error {
	f, e := os.Open(path)
	if e != nil {
		log.Fatalf("can't open `%s`\n", path)
	}
	defer f.Close()

	tty, _ := tty.Open()
	defer func() {
		if v := recover(); v != nil {
			fmt.Printf("panic: %v\n", v)
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
	ctx, dsp, kb, e := StarTermbox(context.Background(), termbox.Attribute(blockColor))
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	chip := &core.Chip8{
		Cpu:      core.NewCpu(time.NewTicker(time.Second/time.Duration(cpuHz)), nil),
		Memory:   &core.Memory{},
		Display:  dsp,
		Keyboard: kb,
	}
	_, e = chip.Init(f)
	if e != nil {
		log.Fatalln(e)
	}
	chip.Run(ctx)
	return nil
}
