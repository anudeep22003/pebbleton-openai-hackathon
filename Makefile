-include .env
export

.PHONY: format test simulate backend frontend install build start

format:
	gofmt -w backend
	cd frontend && npm run format

test:
	go test ./...

simulate:
	go run ./backend/cmd/simulate

backend:
	go run ./backend/cmd/server

frontend:
	cd frontend && npm run dev

install:
	cd frontend && npm ci

build:
	cd frontend && npm run build

start: build backend
