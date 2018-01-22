#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64

go build -i ../http-lambda.go
docker build -t nextmetaphor/ecs-http-lambda:latest .

# docker push nextmetaphor/ecs-http-lambda:latest

docker run -d --name=ecs-http-lambda -p 18080:18080 nextmetaphor/ecs-http-lambda:latest

aws cloudformation create-stack --stack-name ecs-http-lambda --template-body file://ecs.yaml --disable-rollback --parameters file://ecs-parameters.json --capabilities CAPABILITY_IAM