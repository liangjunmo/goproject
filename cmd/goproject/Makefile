GO := GO111MODULE=on GOPROXY=https://goproxy.cn go
GO_LDFLAGS := "-X 'github.com/liangjunmo/goproject/internal/goproject.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/liangjunmo/goproject/internal/goproject.GoVersion=`go version`' -X 'github.com/liangjunmo/goproject/internal/goproject.GitCommit=`git describe --all --long`'"

ifndef DOCKER_REPOSITORY
DOCKER_REPOSITORY := localhost:5000/goproject
endif

.DEFAULT_GOAL := build

PHONY = build
build:
	$(GO) build -race -ldflags ${GO_LDFLAGS} -o ./goproject ./

PHONY = docker-publish-goproject-nginx
docker-publish-goproject-nginx:
	docker build -f ./nginx/Dockerfile -t $(DOCKER_REPOSITORY)/goproject-nginx:$(VERSION) .
	docker push $(DOCKER_REPOSITORY)/goproject-nginx:$(VERSION)

PHONY = docker-publish-goproject
docker-publish-goproject:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags ${GO_LDFLAGS} -o ./goproject ./
	docker build -f ./Dockerfile -t $(DOCKER_REPOSITORY)/goproject:$(VERSION) .
	docker push $(DOCKER_REPOSITORY)/goproject:$(VERSION)

# deploy k8s >>>

PHONY = k8s-deploy-goproject-nginx
k8s-deploy-goproject-nginx: k8s-delete-goproject-nginx k8s-create-goproject-nginx

k8s-delete-goproject-nginx:
	-kubectl delete -f ./k8s/nginx/deployment.yaml
	-kubectl delete -f ./k8s/nginx/service.yaml

k8s-create-goproject-nginx:
	DOCKER_REPOSITORY=$(DOCKER_REPOSITORY) VERSION=$(VERSION) kubectl create -f ./k8s/nginx/deployment.yaml
	kubectl create -f ./k8s/nginx/service.yaml

PHONY = k8s-deploy-goproject-api
k8s-deploy-goproject-api: k8s-delete-goproject-api k8s-create-goproject-api

k8s-delete-goproject-api:
	-kubectl delete -f ./k8s/api/deployment.yaml
	-kubectl delete -f ./k8s/api/service.yaml

k8s-create-goproject-api:
	DOCKER_REPOSITORY=$(DOCKER_REPOSITORY) VERSION=$(VERSION) kubectl create -f ./k8s/api/deployment.yaml
	kubectl create -f ./k8s/api/service.yaml

PHONY = k8s-deploy-goproject-worker
k8s-deploy-goproject-worker: k8s-delete-goproject-worker k8s-create-goproject-worker

k8s-delete-goproject-worker:
	-kubectl delete -f ./k8s/worker/deployment.yaml

k8s-create-goproject-worker:
	DOCKER_REPOSITORY=$(DOCKER_REPOSITORY) VERSION=$(VERSION) kubectl create -f ./k8s/worker/deployment.yaml

# <<< deploy k8s
