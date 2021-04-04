module github.com/coder2z/ndisk

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/casbin/casbin/v2 v2.12.0
	github.com/casbin/gorm-adapter/v3 v3.0.4
	github.com/coder2z/g-saber v1.0.5
	github.com/coder2z/g-server v1.0.7
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.7.1
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/grpc v1.36.0
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.21.2
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
