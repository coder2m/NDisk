.PHONY: build-image
build-image: build-getway-server

.PHONY: build-getway-server
build-getway-server:
	docker build -f build/getway/Dockerfile -t k8s-getway:v1 .

.PHONY: run-getway
run-getway:
	export SERVER_APP_MODE=dev
	export SERVER_APP_ID=idj$sa2@na21
	export SERVER_APP_DEBUG=true
	go run cmd/getway/main.go run -c=config/confin-getway.toml

.PHONY: build-getway
build-getway:
	go build -o  cmd/getway/getway -ldflags "-X github.com/coder2z/g-server/xapp.appName=ndisk_getway -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 -X github.com/coder2z/g-server/xapp.buildHost=`hostname`" cmd/getway/main.go