.PHONY: generate build test serve clean

generate:
	go tool templ generate

build: generate
	go run ./cmd/build

test: generate
	go test ./...

serve:
	go run ./cmd/preview

clean:
	rm -rf dist
