PROJECT=ascii_counter
NAME=alextonkonogov/${PROJECT}
APP=main
PORT=9090

clear:
	rm -f ${APP} || true

build: clear
	CGO_ENABLED=0
	GOARCH=amd64
	GOOS=linux
	go build -tags netgo -a -v cmd/${APP}.go

run: build
	./${APP}

docker_build:
	rm -f ${APP} || true
	go build main.go
	docker build -t ${NAME} .

docker_run_local: build
	docker build -t test/${PROJECT} .
	docker run -ip ${PORT}:${PORT} test/${PROJECT}

docker_push:
	docker push ${NAME}