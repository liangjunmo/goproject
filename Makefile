GO := GO111MODULE=on GOPROXY=https://goproxy.cn go

.DEFAULT_GOAL := test

PHONY = test
test:
	$(GO) test -count=1 -race ./...

PHONE = swag-goproject-api
swag-goproject-api:
	swag fmt -d ./internal/goproject/api/ -g router.go && swag init -d ./internal/goproject/api/ -g router.go -o ./internal/goproject/api/swagger
