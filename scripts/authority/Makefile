.PHONY: build-image
build-image: build-authority-server

.PHONY: build-authority-server
build-authority-server:
	docker build -f build/authority/Dockerfile -t ndisk-authority:latest .

.PHONY: run-authority
run-authority:
	export SERVER_APP_MODE=dev
	export SERVER_APP_ID=idjaosa2@na21
	export SERVER_APP_DEBUG=true
	go run cmd/authority/main.go

.PHONY: build-authority
build-authority:
	go build -o  cmd/authority/authority -ldflags "-X github.com/coder2z/g-server/xapp.appName=ndisk_authority -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 -X github.com/coder2z/g-server/xapp.buildHost=`hostname`" cmd/authority/main.go
