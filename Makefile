BINARY_APP=bin/auth-api

CMD_AUTH_API = services/auth/cmd/main.go

.PHONY: all
all: build

.PHONY: build
build-auth:
	go build -o $(BINARY_APP) $(CMD_AUTH_API)

.PHONY: run-auth
run-auth:
	go run $(CMD_AUTH_API)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_APP)
	rm -f coverage.out
	rm -f coverage.html
	rm -rf bin/
	rm -rf vendor/
	rm -rf go.sum

.PHONY: lint
lint:
	golangci-lint run ./...
	golangci-lint run ./... --fix --timeout 5m