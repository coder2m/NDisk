/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 19:07
 **/
package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	jwt.StandardClaims
	Name  string
	Email string
	Uid   uint

}

func New() *Info {
	return &Info{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.Time).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    cfg.Issuer,
			NotBefore: time.Now().Unix(),
			Subject:   "JWT",
		},
	}
}

func (user *Info) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	tokens, err := token.SignedString([]byte(cfg.Secret))
	return tokens, err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	}
}

func (user *Info) ParseToken(tokens string) (err error) {
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	user.Email = claim["email"].(string)
	user.Username = claim["name"].(string)
	user.Authority = int(claim["authority"].(float64))
	user.Id = int(claim["id"].(float64))
	return err
}
