PORT = 8000
CONFIG = ./config.yaml

run_tests:
	go test ./... -v

build:
	docker build -t http-lb .

run:
	make build
	docker run -v $(CONFIG):/app/config.yaml:ro -p $(PORT):$(PORT) -d --name http-lb http-lb

run_dev:
	make build
	docker run -v $(CONFIG):/app/config.yaml:ro -p $(PORT):$(PORT) -d -e development=1 --name http-lb http-lb

kill:
	docker container kill http-lb

watch:
	docker logs -f http-lb
