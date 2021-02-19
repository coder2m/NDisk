.PHONY: build-image
build-image: build-authority-server

.PHONY: build-authority-server
build-authority-server:
	docker build -f build/authority/Dockerfile -t ndisk-authority:latest .

.PHONY: run-authority
run-authority:
	go run cmd/authority/main.go

.PHONY: build-authority
build-authority:
	go build -o  cmd/authority/authority cmd/authority/main.go