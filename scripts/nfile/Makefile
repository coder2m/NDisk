.PHONY: build-image
build-image: build-nfile-server

.PHONY: build-nfile-server
build-nfile-server:
	docker build -f build/nfile/Dockerfile -t k8s-nfile:v1 .

.PHONY: run-nfile
run-nfile:
	go run cmd/nfile/main.go run -c=config/confin-nfile.toml

.PHONY: build-nfile
build-nfile:
	go build -o  cmd/nfile/nfile cmd/nfile/main.go