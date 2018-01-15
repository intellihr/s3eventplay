setup: ## Install all the build and lint dependencies
	go get github.com/goreleaser/goreleaser
	dep ensure -vendor-only

build: ## Build a beta version
	go build -o s3eventplay ./cmd/s3eventplay/main.go

build-alpine:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o s3eventplay ./cmd/s3eventplay/main.go

.DEFAULT_GOAL := build
