# gochip-8: CHIP-8 emulartor in Go
gochip-8 is emulator of [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8).

## Build
Download source code and build with `make`.

```sh
git clone --recursive github.com/masu-mi/gochip-8
cd ./gochip-8
make
```

### Requirements

Linux/macOS
Go 1.17

## how to use it
rics
```sh
Usage:
  chip-8-term [command]

Available Commands:
  color       show color chart
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       start CHIP-8 emulator

Flags:
  -h, --help   help for chip-8-term

Use "chip-8-term [command] --help" for more information about a command.
```

### Keyboard layout

**[ESC] stop emulator and exit process.**

1 |2 |3 |4(C)
--|--|--|--
Q(4)|W(5)|E(6)|R(D)
A(7)|S(8)|D(9)|F(E)
Z(A)|X(0)|C(B)|V(F)


### example

```sh
## Space Invaders
./dest/chip-8-term start \
  --color 6 --cpu-hz 300 --keyboard-hz 8 \
  --rom './roms/games/Space Invaders [David Winter].ch8'

## Brix
./dest/chip-8-term start \
  --cpu-hz 300 --keyboard-hz 8 \
  --rom './roms/games/Brix [Andreas Gustafsson, 1990].ch8'
```


https://user-images.githubusercontent.com/603602/151666764-e90ff94c-21a7-416a-9d6a-630c34c33b42.mov
ref. [youtube](https://www.youtube.com/watch?v=EjCcRFEkgVc)
