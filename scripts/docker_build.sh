#!/bin/bash

REPO=$1
LAST_COMMIT_SHA=$(git log --oneline | head -n1 | cut -f1 -d ' ')
TIMESTAMP=$(date +%s)

TAG=$TIMESTAMP-$LAST_COMMIT_SHA

docker build -t $REPO:$TAG .
