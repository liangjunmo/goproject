#!/usr/bin/env bash

goproject-nginx() {
  DOCKER_REGISTRY=${DOCKER_REGISTRY} VERSION=${VERSION} cd ./cmd/goproject/ && make k8s-deploy-goproject-nginx
}

goproject-api() {
  DOCKER_REGISTRY=${DOCKER_REGISTRY} VERSION=${VERSION} cd ./cmd/goproject/ && make k8s-deploy-goproject-api
}

usercenter-rpc() {
  DOCKER_REGISTRY=${DOCKER_REGISTRY} VERSION=${VERSION} cd ./cmd/usercenter/ && make k8s-deploy-usercenter-rpc
}

echo "deploy ${PROJECT}"

case ${PROJECT} in
    goproject-nginx)
      goproject-nginx
      ;;
    goproject-api)
      goproject-api
      ;;
    usercenter-rpc)
      usercenter-rpc
      ;;
    *)
      echo "invalid PROJECT ${PROJECT}"
esac
