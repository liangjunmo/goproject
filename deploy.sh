#!/usr/bin/env bash

nginx() {
  DOCKER_REGISTRY=$DOCKER_REGISTRY make redeploy-nginx
}

server-api() {
  DOCKER_REGISTRY=$DOCKER_REGISTRY make redeploy-server-api
}

echo "deploy $PROJECT"

case $PROJECT in
    nginx)
      nginx
      ;;
    server-api)
      server-api
      ;;
    *)
      echo "invalid PROJECT $PROJECT"
esac
