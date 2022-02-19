
dest = ./dest

.PHONY: all
all: $(dest)/dbg $(dest)/gochip-8

$(dest):
	mkdir $(dest)

$(dest)/dbg: ./cmd/dbg/*.go ./core/* $(dest) ./go.mod
	go mod tidy
	go build -o $@ -tags debug ./$(<D)

$(dest)/gochip-8: ./cmd/gochip-8/*.go ./core/* $(dest) ./go.mod
	go mod tidy
	go build -o $@ ./$(<D)
