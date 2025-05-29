AWS_REGION := us-west-1
DOCKER_REPO := solargarlicband/mailing-list
IMAGE_NAME := mailing-list
ECR_REGISTRY := 339712990370.dkr.ecr.us-west-1.amazonaws.com

build:
	./scripts/docker_build.sh $(DOCKER_REPO)

run:
	./scripts/docker_run_local.sh $(DOCKER_REPO)

fmt:
	go fmt ./...

push:
	./scripts/docker_push.sh $(DOCKER_REPO) $(ECR_REGISTRY)

login-ecr:
	./scripts/ecr_login.sh $(AWS_REGION) $(ECR_REGISTRY)
