.PHONY: build-image
build-image: build-nuser-server

.PHONY: build-nuser-server
build-nuser-server:
	docker build -f build/nuser/Dockerfile -t k8s-nuser:v1 .

.PHONY: run-nuser
run-nuser:
	export SERVER_APP_MODE=dev
	export SERVER_APP_ID=id#j$s*@na21
	export SERVER_APP_DEBUG=true
	go run cmd/nuser/main.go

.PHONY: build-nuser
build-nuser:
	go build -o  cmd/nuser/nuser -ldflags "-X github.com/coder2z/g-server/xapp.appName=ndisk_nuser -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 -X github.com/coder2z/g-server/xapp.buildHost=`hostname`" cmd/nuser/main.go