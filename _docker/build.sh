#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64

go build -i ../http-lambda.go
docker build -t nextmetaphor/ecs-http-lambda:latest .

# ...to run locally...
# docker run -d --name=ecs-http-lambda -p 18080:18080 -l info nextmetaphor/ecs-http-lambda:latest

# ...to push to DockerHub...
# docker push nextmetaphor/ecs-http-lambda:latest