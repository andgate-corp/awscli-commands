build:
	go get github.com/aws/aws-lambda-go/lambda
	env GOOS=linux go build -ldflags="-s -w" -o bin/igniter handler/igniter/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/interaction handler/interaction/main.go