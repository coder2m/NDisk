.PHONY: build-image
build-image: build-nuser-server

.PHONY: build-nuser-server
build-nuser-server:
	docker build -f build/nuser/Dockerfile -t k8s-nuser:v1 .

.PHONY: run-nuser
run-nuser:
	go run cmd/nuser/main.go

.PHONY: build-nuser
build-nuser:
	go build -o  cmd/nuser/nuser cmd/nuser/main.go