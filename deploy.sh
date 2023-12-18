#!/usr/bin/env bash

usercenter-rpc() {
  DOCKER_REGISTRY=$DOCKER_REGISTRY cd ./cmd/usercenter/ && make k8s-deploy-usercenter-rpc
}

goproject-nginx() {
  DOCKER_REGISTRY=$DOCKER_REGISTRY cd ./cmd/goproject/ && make k8s-deploy-goproject-nginx
}

goproject-api() {
  DOCKER_REGISTRY=$DOCKER_REGISTRY cd ./cmd/goproject/ && make k8s-deploy-goproject-api
}

echo "deploy $PROJECT"

case $PROJECT in
    usercenter-rpc)
      usercenter-rpc
      ;;
    goproject-nginx)
      goproject-nginx
      ;;
    goproject-api)
      goproject-api
      ;;
    *)
      echo "invalid PROJECT $PROJECT"
esac
