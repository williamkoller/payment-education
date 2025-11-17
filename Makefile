IGNORED_DIRS := internal/.*/router|shared/.*/event|cmd|internal/.*/domain/port/|pkg/logger|pkg/ulid|shared/middleware|internal/.*/model|config|internal/.*/presentation

PKGS := $(shell go list ./... | grep -vE '($(IGNORED_DIRS))')

test:
	@go test -v $(PKGS)

cover:
	@go test -buildvcs=false -covermode=atomic \
		-coverpkg=$(shell go list ./... | grep -vE '($(IGNORED_DIRS))' | tr '\n' ',') \
		-coverprofile=coverage.out ./... 2>/dev/null
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html


docker-build:
	docker build -t williamkoller/system-education:latest .

k8s-load-image:
	minikube image load williamkoller/system-education:latest

k8s-restart:
	kubectl rollout restart deployment system-education

minikube-launch:
	minikube service system-education

k8s-apply:
	kubectl apply -k k8s/

run-all:
	make docker-build && make k8s-load-image && make k8s-apply && make k8s-restart && make minikube-launch

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Erro: você precisa passar o nome da migração. Ex: make migrate-create name=create_users_table"; \
		exit 1; \
	else \
		migrate create -ext sql -dir db/migrations -seq $(name); \
	fi
