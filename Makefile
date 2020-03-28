build: fmt
	CGO_ENABLED=0 GOOS=linux go build -v -o remember

fmt:
	go fmt ./...

fly-deploy:
	flyctl deploy

fly-logs:
	flyctl --app remember logs

