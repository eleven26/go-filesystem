FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

init:
	go build -modfile=tools/go.mod -o bin/gofumpt mvdan.cc/gofumpt
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

format:
	go mod tidy
	bin/gofumpt -w $(FILES)

check:
	bin/golangci-lint run

tests:
	go test ./...
