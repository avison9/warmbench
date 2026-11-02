.PHONY: build test vet fmt-check

build:
	go build -o bin/ ./cmd/...

test:
	go test ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
