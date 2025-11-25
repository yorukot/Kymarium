BINARY_NAME=knocker

build:
	go build -o tmp/$(BINARY_NAME) ./cmd
	chmod +x tmp/$(BINARY_NAME)

run: build
	./tmp/$(BINARY_NAME)

web:
	cd website && pnpm run dev

dev:
	air --build.cmd "go build -o tmp/$(BINARY_NAME) ./cmd" --build.bin "./tmp/$(BINARY_NAME)"

api:
	air --build.cmd "go build -o tmp/$(BINARY_NAME) ./cmd" --build.bin "./tmp/$(BINARY_NAME) api"

worker:
	air --build.cmd "go build -o tmp/$(BINARY_NAME) ./cmd" --build.bin "./tmp/$(BINARY_NAME) worker"

schedular:
	air --build.cmd "go build -o tmp/$(BINARY_NAME) ./cmd" --build.bin "./tmp/$(BINARY_NAME) schedular"

test:
	go test ./...

lint:
	go fmt ./...
	go vet ./...
	golint ./...

generate-docs:
	swag init -g cmd/main.go -o ./docs

clean:
	rm -rf tmp/

seed:
	go run ./cmd/seed

.PHONY: build run web dev api worker schedular test lint generate-docs clean seed
