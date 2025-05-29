#!/bin/bash

IMAGE_NAME=$1
IMAGE_ID=$(docker image ls | grep $IMAGE_NAME | tr -s ' ' | cut -f2,3 -d ' ' | sort -r | head -n1 | cut -f2 -d ' ')
CONTAINER_NAME=${2:-solar-garlic-mailing-list}
PORT=${3:-8080}

docker run -p $PORT:8080 --env-file .env -d --name $CONTAINER_NAME --network solar-garlic $IMAGE_ID
