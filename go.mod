module github.com/myxy99/ndisk

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.3
	github.com/mitchellh/mapstructure v1.4.0
	github.com/myxy99/component v0.3.4
	google.golang.org/grpc v1.27.0
	gorm.io/gorm v1.20.9
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
