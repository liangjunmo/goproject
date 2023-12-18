GO := GO111MODULE=on GOPROXY=https://goproxy.cn go
GO_LDFLAGS := "-X 'github.com/liangjunmo/goproject/internal/goproject.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/liangjunmo/goproject/internal/goproject.GoVersion=`go version`' -X 'github.com/liangjunmo/goproject/internal/goproject.GitCommit=`git describe --all --long`'"

TIME := $(shell date +%Y%m%d%H%M%S)

ifndef DOCKER_REPOSITORY
DOCKER_REPOSITORY := localhost:5000/goproject
endif

.DEFAULT_GOAL := build

PHONY = build
build:
	$(GO) build -race -ldflags ${GO_LDFLAGS} -o ./goproject ./

PHONY = docker-publish
docker-publish:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags ${GO_LDFLAGS} -o ./goproject ./
	docker build -f ./Dockerfile -t $(DOCKER_REPOSITORY)/goproject:$(VERSION) .
	docker push $(DOCKER_REPOSITORY)/goproject:$(VERSION)

# deploy k8s >>>

PHONY = k8s-deploy-goproject-nginx
k8s-deploy-goproject-nginx: k8s-docker-build-goproject-nginx k8s-docker-push-goproject-nginx k8s-overwrite-yaml k8s-delete-goproject-nginx k8s-create-goproject-nginx k8s-reset-yaml

PHONY = k8s-deploy-goproject-api
k8s-deploy-goproject-api: k8s-go-build-goproject k8s-docker-build-goproject k8s-docker-push-goproject k8s-overwrite-yaml k8s-delete-goproject-api k8s-create-goproject-api k8s-reset-yaml

PHONY = k8s-deploy-goproject-worker
k8s-deploy-goproject-worker: k8s-go-build-goproject k8s-docker-build-goproject k8s-docker-push-goproject k8s-overwrite-yaml k8s-delete-goproject-worker k8s-create-goproject-worker k8s-reset-yaml

k8s-overwrite-yaml:
	sed -i "s#DOCKER_REPOSITORY#$(DOCKER_REPOSITORY)#g" `find ./ -type f -name "*.yaml"`
	sed -i "s#latest#$(TIME)#g" `find ./ -type f -name "*.yaml"`

k8s-reset-yaml:
	sed -i "s#$(TIME)#latest#g" `find ./ -type f -name "*.yaml"`
	sed -i "s#$(DOCKER_REPOSITORY)#DOCKER_REPOSITORY#g" `find ./ -type f -name "*.yaml"`

k8s-docker-build-goproject-nginx:
	docker build -f ./k8s/nginx/Dockerfile -t $(DOCKER_REPOSITORY)/goproject-nginx:$(TIME) .

k8s-docker-push-goproject-nginx:
	docker push $(DOCKER_REPOSITORY)/goproject-nginx:$(TIME)

k8s-delete-goproject-nginx:
	-kubectl delete -f ./k8s/nginx/deployment.yaml
	-kubectl delete -f ./k8s/nginx/service.yaml

k8s-create-goproject-nginx:
	kubectl create -f ./k8s/nginx/deployment.yaml
	kubectl create -f ./k8s/nginx/service.yaml

k8s-go-build-goproject:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags ${GO_LDFLAGS} -o ./goproject ./

k8s-docker-build-goproject:
	docker build -f ./k8s/Dockerfile -t $(DOCKER_REPOSITORY)/goproject:$(TIME) .

k8s-docker-push-goproject:
	docker push $(DOCKER_REPOSITORY)/goproject:$(TIME)

k8s-delete-goproject-api:
	-kubectl delete -f ./k8s/api/deployment.yaml
	-kubectl delete -f ./k8s/api/service.yaml

k8s-create-goproject-api:
	kubectl create -f ./k8s/api/deployment.yaml
	kubectl create -f ./k8s/api/service.yaml

k8s-delete-goproject-worker:
	-kubectl delete -f ./k8s/worker/deployment.yaml

k8s-create-goproject-worker:
	kubectl create -f ./k8s/worker/deployment.yaml

# <<< deploy k8s
