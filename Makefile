GOPATH:=$(shell go env GOPATH)
APP?=apigateway

.PHONY: install-tools
## install-tools: installs dependencies for tools
install-tools:
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: build
## build: build the application(api)
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o bin/${APP}.linux.amd64 cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags="-w -s" -o bin/${APP}.linux.arm64 cmd/main.go
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o bin/${APP} cmd/main.go

.PHONY: run
## run: run the application(api)
run:
	go run -v -race cmd/main.go

.PHONY: format-go
## format-go: formats go files
format-go: install-tools
	find . -not -path './gen/*' -not -path './tools.go' -not -name '*_mock.go' -name '*.go' -print0 | xargs -0 -I {} goimports-reviser -rm-unused -format -company-prefixes github.com/taehoio {}
	go mod tidy

.PHONY: format
## format: formats files
format: format-go

.PHONY: test
## test: runs tests
test: install-tools
	gotest -p 1 -race -cover -v ./...

.PHONY: coverage
## coverage: runs tests with coverage
coverage: install-tools
	gotest -p 1 -race -coverprofile=coverage.out -covermode=atomic -v ./...

.PHONY: generate-mock
## generate-mock: generates mock files
generate-mock: install-tools
	go generate ./...

.PHONY: generate
## generate: generates files
generate: generate-mock

.PHONY: lint-go
## lint-go: lints go files
lint-go: install-tools
	golangci-lint run ./...
	find . -not -path './gen/*' -not -path './tools.go' -not -name '*_mock.go' -name '*.go' -print0 | xargs -0 -I {} goimports-reviser -rm-unused -format -company-prefixes github.com/taehoio -list-diff -set-exit-status {}
	go mod verify

.PHONY: lint
## lint: lints files
lint: lint-go

.PHONY: clean
## clean: cleans generated files
clean:
	rm -rf gen

.PHONY: diff
## diff: shows diff
diff:
	git diff --exit-code
	if [ -n "$(git status --porcelain)" ]; then git status; exit 1; else exit 0; fi

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'
