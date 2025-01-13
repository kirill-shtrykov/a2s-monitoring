version := $(shell cat VERSION)

test:
	golangci-lint run

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/a2s-monitoring

compress:
	gzip -c build/a2s-monitoring > build/a2s-monitoring-$(version)-linux-amd64.gz

.PHONY: test build compress
