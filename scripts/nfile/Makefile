.PHONY: build-image
build-image: build-nfile-server

.PHONY: build-nfile-server
build-nfile-server:
	docker build -f build/nfile/Dockerfile -t k8s-nfile:v1 .

.PHONY: run-nfile
run-nfile:
	export SERVER_APP_MODE=dev
	export SERVER_APP_ID=idj$s*@na21
	export SERVER_APP_DEBUG=true
	go run cmd/nfile/main.go run -c=config/confin-nfile.toml

.PHONY: build-nfile
build-nfile:
	go build -o  cmd/nfile/nfile -ldflags "-X github.com/coder2z/g-server/xapp.appName=ndisk_nfile -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 -X github.com/coder2z/g-server/xapp.buildHost=`hostname`" cmd/nfile/main.go