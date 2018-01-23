#!/usr/bin/env bash

aws cloudformation create-stack --stack-name ecs-http-lambda --template-body file://ecs.yaml --disable-rollback --parameters file://ecs-parameters.json --capabilities CAPABILITY_IAM