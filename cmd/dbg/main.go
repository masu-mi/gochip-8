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

	chip := &core.Chip8{Cpu: core.NewCpu(nil), Memory: &core.Memory{}}
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
	tty, _ := tty.Open()
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("panic")
		}
		tty.Close()
	}()
	num := 0
	for {
		_, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("tick(%d): Pc: %04x(%d)", num, chip.Cpu.Pc, chip.Cpu.Pc)
		chip.Tick()
		num++
	}
}
