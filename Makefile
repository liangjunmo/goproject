GO := GO111MODULE=on GOPROXY=https://goproxy.cn go

.DEFAULT_GOAL := test

PHONY = test
test:
	$(GO) test -count=1 -race ./...
