.PHONY: run-auth build-auth tidy test

# Auth module
run-auth:
	go run modules/auth/main.go

build-auth:
	go build -o bin/auth modules/auth/main.go

# General
tidy:
	go mod tidy

test:
	go test ./... -v
