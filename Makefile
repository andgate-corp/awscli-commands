aws_profile = default
region = us-east-1

build-docker-golang:
	docker build -f Dockerfile . -t awscli-command:latest
dep-ensure:
	dep ensure
build-golang:
	go get github.com/aws/aws-lambda-go/lambda
	env GOOS=linux go build -ldflags="-s -w" -o bin/igniter handler/igniter/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/interaction handler/interaction/main.go

run-dep-ensure:
	docker run --rm -v "$$PWD":/go/src/app -w /go/src/app awscli-command:latest make dep-ensure

run-build:
ifeq ($$(OS),Windows_NT)
# for Windows
	docker run --rm -v "$$PWD":/go/src/app -w /go/src/app awscli-command:latest make build-golang
else
	docker run --rm -v "$$PWD":/go/src/app -w /go/src/app awscli-command:latest make build-golang
endif

sls-deploy: 
	docker run --rm -v $${HOME}/.aws/credentials:/root/.aws/credentials -v $${PWD}:/go/src/app \
	-e AWS_PROFILE=${aws_profile} -e AWS_DEFAULT_REGION=${region} \
	-w /go/src/app \
	softinstigate/serverless sls deploy
