FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

init:                       ## Install linters.
	go build -modfile=tools/go.mod -o bin/gofumports mvdan.cc/gofumpt/gofumports
	go build -modfile=tools/go.mod -o bin/gofumpt mvdan.cc/gofumpt
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
	go build -modfile=tools/go.mod -o bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog

format:                     ## Format source code.
	go mod tidy
	bin/gofumpt -w -s $(FILES)
	bin/gofumports -local gitlab.gzjztw.com/alert-api -l -w $(FILES)

check:                      ## Run checks/linters
	bin/golangci-lint run

tests:
	go test ./...
