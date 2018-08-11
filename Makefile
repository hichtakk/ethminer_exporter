.PHONY: build clean test

DOCKERHUB_ID := hichtakk
IMAGE_NAME := ethminer_exporter
IMAGE_TAG := latest
IMAGE_REPOSITORY := ${DOCKERHUB_ID}/${IMAGE_NAME}:${IMAGE_TAG}

bin/ethminer_exporter:
	@GOOS=linux GOARCH=amd64 go build -o bin/ethminer_exporter ./main.go

build: bin/ethminer_exporter
	@docker build -t ${IMAGE_REPOSITORY} .

release: build
	@docker login
	@docker push ${IMAGE_REPOSITORY}

test:
	docker run -d --name test -p 8555:8555 ${IMAGE_REPOSITORY}

clean:
	@rm -rf ./bin
	#@for i in `docker images | grep ${IMAGE_NAME} | awk '{print $$1 ":" $$2}'`; do docker rmi $$1 ; done
