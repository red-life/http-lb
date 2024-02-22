CONFIG = ./config.yaml
CONTAINER_PORT = $(shell cat $(CONFIG) | grep "listen:" | cut -d ":" -f 3)
CONTAINER_NAME = http-lb
VOLUMES = -v $(CONFIG):/app/config.yaml:ro # -v ./key.key:/app/key.key:ro -v ./cert.crt:/app/cert.crt:ro
HOST_PORT = 5000
PORTS = -p $(HOST_PORT):$(CONTAINER_PORT)
RUN_ARGS = $(VOLUMES) $(PORTS) -d --name $(CONTAINER_NAME)

run_tests:
	go test ./... -v

build:
	docker build -t http-lb .

run:
	make build
	docker run $(RUN_ARGS) http-lb

run_dev:
	make build
	docker run $(RUN_ARGS) -e development=1 http-lb

kill:
	docker container kill $(CONTAINER_NAME)

watch:
	docker logs -f $(CONTAINER_NAME)
