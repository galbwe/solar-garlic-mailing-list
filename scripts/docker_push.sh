#!/bin/bash

REPO=$1
REGISTRY=$2
IMAGE_ID=$(docker image ls | grep ^$REPO | tr -s ' ' | cut -f2,3 -d ' ' | sort -r | head -n1 | cut -f2 -d ' ');
LAST_COMMIT_SHA=$(git log --oneline | head -n1 | cut -f1 -d ' ');

echo IMAGE_ID=$IMAGE_ID
echo LAST_COMMIT_SHA=$LAST_COMMIT_SHA

echo tagging ...
echo tag $IMAGE_ID $REGISTRY/$REPO:$LAST_COMMIT_SHA; 
docker tag $IMAGE_ID $REGISTRY/$REPO:$LAST_COMMIT_SHA; 

echo pushing ...
echo push $REGISTRY/$REPO:$LAST_COMMIT_SHA; 
docker push $REGISTRY/$REPO:$LAST_COMMIT_SHA; 
