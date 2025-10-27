IGNORED_DIRS := internal/.*/router|shared/.*/event|cmd|internal/.*/domain/port/

PKGS := $(shell go list ./... | grep -vE '($(IGNORED_DIRS))')

test:
	@go test -v $(PKGS)

cover:
	@go test -buildvcs=false -covermode=atomic \
		-coverpkg=$(shell go list ./... | grep -vE '($(IGNORED_DIRS))' | tr '\n' ',') \
		-coverprofile=coverage.out ./... 2>/dev/null
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html


server:
	@go run cmd/main.go
