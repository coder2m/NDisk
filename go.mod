module ndisk

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/myxy99/component latest
	github.com/spf13/pflag v1.0.3 // indirect
	gorm.io/gorm v1.20.8
)

replace (
	github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
