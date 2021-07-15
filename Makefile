init:
	bash build/install_server.sh

docker_up:
	docker-compose -f build/docker-compose.yml up -d

docker_down:
	docker-compose -f build/docker-compose.yml down

unit_test:
	go test ./...

integration_test:
	go test -tags=integration ./...

run_directly:
	go run *.go -state local

build_image:
	docker build -t tiki-go .

run_container:
	docker run tiki-go