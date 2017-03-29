build-linux:
	GOOS=linux GOARCH=amd64 go build

build-docker: build-linux
	docker build -t fresh8/consulship .

push-docker: build-docker
	docker push fresh8/consulship
