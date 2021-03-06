//go:build debug

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
		Cpu:      core.NewCpu(nil, nil),
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
	num := 0
	for range forRepl {
		fmt.Printf("tick(%d): Pc: %04x(%d)\n", num, chip.Cpu.Pc, chip.Cpu.Pc)
		chip.Cycle()
		num++
	}
}
