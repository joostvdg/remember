
PROJECT_ID := remember-272515

build: fmt
	CGO_ENABLED=0 GOOS=linux go build -v -o remember

test:
	CGO_ENABLED=0 GOOS=linux go test -v -coverprofile cover.out ./...

fmt:
	go fmt ./...

fly-deploy:
	flyctl deploy

fly-logs:
	flyctl --app remember logs

dbuild: fmt
	DOCKER_BUILDKIT=1 docker build --tag caladreas/remember:latest .

dpush: dbuild
	docker push caladreas/remember:latest

gbuild: fmt
	gcloud builds submit --tag gcr.io/$(PROJECT_ID)/remember

gpush: dbuild
	docker tag caladreas/remember:latest gcr.io/$(PROJECT_ID)/remember:latest
	docker push gcr.io/$(PROJECT_ID)/remember:latest

gdeploy: dpush
	gcloud run deploy remember --image=gcr.io/$(PROJECT_ID)/remember:latest --platform managed