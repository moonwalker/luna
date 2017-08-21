COMMIT = $(shell git rev-parse --short HEAD)

build:
	@rm -rf tmp
	go build -ldflags "-X main.GitCommit=$(COMMIT)" -o tmp/luna-`uname -s`-`uname -m`
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.GitCommit=$(COMMIT)" -o tmp/luna-Linux-x86_64
