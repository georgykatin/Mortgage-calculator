# container config
CONTAINER_NAME=mortgage_service
IMAGE_NAME=mortgage_service:latest

# run tests 
.PHONY: test
test:
	go test -cover ./...

# run lint checkout
.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

# build docker image
.PHONY: build
build:
	@docker build -t $(IMAGE_NAME) .

# run container
.PHONY: run
run:
	@docker run -d --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)

# stop and delete container
.PHONY: stop
stop:
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true
