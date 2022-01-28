package main

import (
	"fmt"
	"log"
	"os"

	"github.com/masu-mi/gochip-8/core"
	"github.com/mattn/go-tty"
	"github.com/spf13/cobra"
)

var (
	fps        uint8
	path       string
	blockColor uint64
)

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start CHIP-8 emulator",
		RunE:  start,
	}
	cmd.PersistentFlags().Uint8Var(&fps, "keyboard-hz", 60, "reciprocal of duration of key pressed (default: 60Hz)")
	cmd.PersistentFlags().StringVar(&path, "rom", "", "rom image file path")
	cmd.PersistentFlags().Uint64Var(&blockColor, "color", 16, "display active cell's color(defalt: 16)")
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
	n, e := chip.Init(f)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Printf("load: %d[byte]\n", n)

	chip.Memory.WriteTo(os.Stdout)
	num := 0
	for range forRepl {
		fmt.Printf("tick(%d): Pc: %04x(%d)\n", num, chip.Cpu.Pc, chip.Cpu.Pc)
		chip.Tick()
		num++
	}
	return nil
}
