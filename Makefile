
dest = ./dest

all: $(dest)/dbg $(dest)/chip-8-term
.PHONY: all

$(dest):
	mkdir $(dest)

$(dest)/dbg: ./cmd/dbg/*.go ./core/* $(dest)
	go build -o $@ -tags debug ./$(<D)

$(dest)/chip-8-term: ./cmd/chip-8-term/*.go ./core/* $(dest)
	go build -o $@ ./$(<D)
