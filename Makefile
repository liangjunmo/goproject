GO := GO111MODULE=on GOPROXY=https://goproxy.cn go
GO_LDFLAGS := "-X 'github.com/liangjunmo/goproject/internal/version.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/liangjunmo/goproject/internal/version.GoVersion=`go version`' -X 'github.com/liangjunmo/goproject/internal/version.GitCommit=`git describe --all --long`'"

.DEFAULT_GOAL := default

PHONY = default
default:
	@echo "please specify a target to run"

PHONY = build
build:
	@mkdir -p tmp
	$(GO) build -race -ldflags ${GO_LDFLAGS} -o ./tmp/goproject-server ./cmd/server/
