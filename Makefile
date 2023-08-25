GO := GO111MODULE=on GOPROXY=https://goproxy.cn go
GO_LDFLAGS := "-X 'github.com/liangjunmo/goproject/internal/version.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/liangjunmo/goproject/internal/version.GoVersion=`go version`' -X 'github.com/liangjunmo/goproject/internal/version.GitCommit=`git describe --all --long`'"

TIME := $(shell date +%Y%m%d%H%M%S)

DOCKER_REPOSITORY_GOPROJECT := localhost:5000/goproject
ifdef DOCKER_REGISTRY
DOCKER_REPOSITORY_GOPROJECT = $(DOCKER_REGISTRY)/goproject
endif

.DEFAULT_GOAL := default

PHONY = default
default:
	@echo "please specify a target to run"

PHONY = build
build:
	@mkdir -p tmp
	$(GO) build -race -ldflags ${GO_LDFLAGS} -o ./tmp/goproject-server ./cmd/server/

# deploy docker >>>

PHONY = run-dev-middleware
run-dev-middleware:
	docker compose -f ./deploy/docker/middleware/docker-compose.yaml -p "goproject-dev-middleware" up -d

PHONY = start-dev-middleware
start-dev-middleware:
	docker compose -f ./deploy/docker/middleware/docker-compose.yaml -p "goproject-dev-middleware" start

PHONY = stop-dev-middleware
stop-dev-middleware:
	docker compose -f ./deploy/docker/middleware/docker-compose.yaml -p "goproject-dev-middleware" stop

PHONY = run-dev-server
run-dev-server: run-dev-middleware
	@mkdir -p tmp
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags ${GO_LDFLAGS} -o tmp/goproject-server ./cmd/server/
	chmod u+x ./tmp/goproject-server
	docker build -f ./deploy/docker/server/Dockerfile -t goproject-server:latest .
	docker compose -f ./deploy/docker/server/docker-compose.yaml -p "goproject-dev-server" up -d
	-docker rmi $(shell docker images -q goproject-server)

PHONY = start-dev-server
start-dev-server: start-dev-middleware
	docker compose -f ./deploy/docker/server/docker-compose.yaml -p "goproject-dev-server" start api worker1

PHONY = stop-dev-server
stop-dev-server:
	docker compose -f ./deploy/docker/server/docker-compose.yaml -p "goproject-dev-server" stop api worker1

# <<< deploy docker

# deploy k8s >>>

rewrite-yaml:
	sed -i "s#DOCKER_REPOSITORY_GOPROJECT#$(DOCKER_REPOSITORY_GOPROJECT)#g" `find ./ -type f -name "*.yaml"`
	sed -i "s#latest#$(TIME)#g" `find ./ -type f -name "*.yaml"`

recover-yaml:
	sed -i "s#$(TIME)#latest#g" `find ./ -type f -name "*.yaml"`
	sed -i "s#$(DOCKER_REPOSITORY_GOPROJECT)#DOCKER_REPOSITORY_GOPROJECT#g" `find ./ -type f -name "*.yaml"`

define panic-and-recover-yaml
	(sed -i "s#$(TIME)#latest#g" `find ./ -type f -name "*.yaml"` && sed -i "s#$(DOCKER_REPOSITORY_GOPROJECT)#DOCKER_REPOSITORY_GOPROJECT#g" `find ./ -type f -name "*.yaml"` && exit 1)
endef

PHONY = redeploy-nginx
redeploy-nginx: docker-build-nginx docker-push-nginx rewrite-yaml delete-nginx create-nginx recover-yaml

docker-build-nginx:
	docker build -f ./deploy/k8s/goproject/nginx/Dockerfile -t $(DOCKER_REPOSITORY_GOPROJECT)/goproject-nginx:$(TIME) .

docker-push-nginx:
	docker push $(DOCKER_REPOSITORY_GOPROJECT)/goproject-nginx:$(TIME)

create-nginx:
	kubectl create -f ./deploy/k8s/goproject/nginx/nginx-deployment.yaml || $(call panic-and-recover-yaml)
	kubectl create -f ./deploy/k8s/goproject/nginx/nginx-service.yaml || $(call panic-and-recover-yaml)

delete-nginx:
	-kubectl delete -f ./deploy/k8s/goproject/nginx/nginx-deployment.yaml
	-kubectl delete -f ./deploy/k8s/goproject/nginx/nginx-service.yaml

PHONY = redeploy-server-api
redeploy-server-api: go-build-server docker-build-server docker-push-server delete-server-api delete-server-worker1 rewrite-yaml create-server-api create-server-worker1 recover-yaml

go-build-server:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags ${GO_LDFLAGS} -o ./tmp/goproject-server ./cmd/server/

docker-build-server:
	docker build -f ./deploy/k8s/goproject/server/Dockerfile -t $(DOCKER_REPOSITORY_GOPROJECT)/goproject-server:$(TIME) .

docker-push-server:
	docker push $(DOCKER_REPOSITORY_GOPROJECT)/goproject-server:$(TIME)

create-server-api:
	kubectl create -f ./deploy/k8s/goproject/server/api-deployment.yaml || $(call panic-and-recover-yaml)
	kubectl create -f ./deploy/k8s/goproject/server/api-service.yaml || $(call panic-and-recover-yaml)

delete-server-api:
	-kubectl delete -f ./deploy/k8s/goproject/server/api-deployment.yaml
	-kubectl delete -f ./deploy/k8s/goproject/server/api-service.yaml

create-server-worker1:
	kubectl create -f ./deploy/k8s/goproject/server/worker1-deployment.yaml || $(call panic-and-recover-yaml)

delete-server-worker1:
	-kubectl delete -f ./deploy/k8s/goproject/server/worker1-deployment.yaml

# <<< deploy k8s
