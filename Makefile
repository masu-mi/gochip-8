

./dbg: ./core/*.go ./cmd/dbg/main.go
	go build -tags debug ./cmd/dbg/
