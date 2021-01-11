module github.com/myxy99/ndisk

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.4.4
	github.com/go-resty/resty/v2 v2.3.0
	github.com/golang/protobuf v1.4.3
	github.com/myxy99/component v0.4.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/grpc v1.27.0
	gorm.io/gorm v1.20.9
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
