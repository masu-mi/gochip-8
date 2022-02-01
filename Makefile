
dest = ./dest

all: $(dest)/dbg $(dest)/gochip-8
.PHONY: all

$(dest):
	mkdir $(dest)

$(dest)/dbg: ./cmd/dbg/*.go ./core/* $(dest)
	go build -o $@ -tags debug ./$(<D)

$(dest)/gochip-8: ./cmd/gochip-8/*.go ./core/* $(dest)
	go build -o $@ ./$(<D)
