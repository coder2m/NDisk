.PHONY: build-image
build-image: build-getway-server

.PHONY: build-getway-server
build-getway-server:
	docker build -f build/getway/Dockerfile -t k8s-getway:v1 .

.PHONY: run-getway
run-getway:
	go run cmd/getway/main.go run -c=config/confin-getway.toml

.PHONY: build-getway
build-getway:
	go build -o  cmd/getway/getway cmd/getway/main.go