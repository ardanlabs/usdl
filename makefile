# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

install:
	go install github.com/a-h/templ/cmd/templ@latest
	brew install entr

# ==============================================================================
# Chat

hack:
	go run chat/api/tooling/hack/main.go

run-cap:
	go run chat/api/services/cap/main.go | go run chat/api/tooling/logfmt/main.go

run-client:
	templ generate chat/api/frontends/client/ui/web/
	go run chat/api/frontends/client/main.go

run-client-reload:
	find . -name "*.go" -o -name "*.html" -o -name "*.css" | entr -r go run chat/api/frontends/client/main.go & \
	find . -name "*.templ" | entr -r templ generate chat/api/frontends/client/ui/web/ & \
	wait;

chat-test:
	curl -i -X GET http://localhost:3000/test

chat-docker:
	docker pull nats:2.10

chat-nats:
	docker run -p 4222:4222 nats:2.10 -js

chat-nats-down:
	docker stop nats:2.10
	docker rm nats:2.10 -v

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check
