GO := GO111MODULE=on GOPROXY=https://goproxy.cn go
GO_LDFLAGS := "-X 'github.com/liangjunmo/goproject/internal/usercenter.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/liangjunmo/goproject/internal/usercenter.GoVersion=`go version`' -X 'github.com/liangjunmo/goproject/internal/usercenter.GitCommit=`git describe --all --long`'"

ifndef DOCKER_REPOSITORY
	DOCKER_REPOSITORY = localhost:5000/goproject
endif

.DEFAULT_GOAL := build

PHONY = build
build:
	$(GO) build -race -ldflags ${GO_LDFLAGS} -o ./usercenter ./

PHONY = docker-publish-usercenter
docker-publish-usercenter:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags ${GO_LDFLAGS} -o ./usercenter ./
	docker build -f ./Dockerfile -t $(DOCKER_REPOSITORY)/goproject-usercenter:$(VERSION) .
	docker push $(DOCKER_REPOSITORY)/goproject-usercenter:$(VERSION)

# deploy k8s >>>

PHONY = k8s-deploy-usercenter-rpc
k8s-deploy-usercenter-rpc: k8s-delete-usercenter-rpc k8s-create-usercenter-rpc

k8s-delete-usercenter-rpc:
	-kubectl delete -f ./k8s/rpc/deployment.yaml
	-kubectl delete -f ./k8s/rpc/service.yaml

k8s-create-usercenter-rpc:
	DOCKER_REPOSITORY=$(DOCKER_REPOSITORY) VERSION=$(VERSION) kubectl create -f ./k8s/rpc/deployment.yaml
	kubectl create -f ./k8s/rpc/service.yaml

# <<< deploy k8s
