# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Installation

install:
	go install github.com/a-h/templ/cmd/templ@latest
	brew install entr
	brew install ollama
	brew install tailscale

docker:
	docker pull dyrnq/open-webui:main

# ==============================================================================
# Tailscale

tailscale-up:
	tailscale serve --bg --tcp 4001 4000

tailscale-down:
	tailscale serve --tcp 4001 off

# ==============================================================================
# Ollama

ollama-up:
	export OLLAMA_MODELS="zarf/ollama/models" && \
	export OLLAMA_DEBUG=1 && \
	ollama serve

ollama-pull:
	ollama pull llama3.2

ollama-logs:
	tail -f -n 100 ~/.ollama/logs/server.log

# ==============================================================================
# OpenWebUI

compose-up:
	docker compose -f zarf/docker/compose.yaml up

compose-down:
	docker compose -f zarf/docker/compose.yaml down

compose-logs:
	docker compose logs -n 100

openwebui:
	open -a "Google Chrome" http://localhost:3000/

# ==============================================================================
# Chat

hack:
	go run api/tooling/hack/main.go

run-cap:
	go run api/services/cap/main.go | go run api/tooling/logfmt/main.go

run-tui:
	go run api/clients/tui/main.go

run-tui-ai:
	go run api/clients/tui/main.go --aimode=true

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
