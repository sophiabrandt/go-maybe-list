SHELL := /bin/bash

export PROJECT = go-maybe-list

# ==============================================================================
# Development

run: up dev

up:
	docker-compose up -d

dev:
	go run ./cmd/web

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Administration

migrate:
	go run ./cmd/admin -action="migrate"

seed: migrate
	go run ./cmd/admin -action="seed"

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck ./...
