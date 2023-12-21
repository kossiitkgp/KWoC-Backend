#!/usr/bin/env bash

set -xe

REPO_PATH=$(git rev-parse --show-toplevel)

cd $REPO_PATH/tests
source .env

docker compose -f docker-compose.test.yaml up --build -d

go test ./... -p 1

docker compose -f docker-compose.test.yaml down
