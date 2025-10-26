IGNORED_DIRS := internal/.*/router|shared/.*/event|cmd

PKGS := $(shell go list ./... | grep -vE '($(IGNORED_DIRS))')

test:
	@echo $(PKGS)
	@go test -v $(PKGS)

cover:
	@echo $(PKGS)
	@go test -buildvcs=false -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -covermode=atomic -coverprofile=coverage.out $(PKGS)
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

server:
	@go run cmd/main.go
