setup: ## Install all the build and lint dependencies
	go get github.com/goreleaser/goreleaser

build: ## Build a beta version
	go build -o s3eventplay ./cmd/s3eventplay/main.go

.DEFAULT_GOAL := build
