#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64

go build -i ../http-lambda.go
docker build -t nextmetaphor/ecs-http-lambda:latest .

docker run -d --name=ecs-http-lambda -p 18080:18080 nextmetaphor/ecs-http-lambda:latest