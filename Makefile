.PHONY: format test simulate

format:
	gofmt -w backend

test:
	go test ./...

simulate:
	go run ./backend/cmd/simulate
