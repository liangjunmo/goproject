ifndef DOCKER_REPOSITORY
	DOCKER_REPOSITORY = localhost:5000/goproject
endif

.DEFAULT_GOAL := run

PHONY = run
run:
	DOCKER_REPOSITORY=${DOCKER_REPOSITORY} docker compose -f ./docker-compose.yaml -p "goproject" up -d
