package token

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/xcfg"
	"strconv"
	"time"
)

const (
	DefaultAccessKey = "ecol123og1ysK#xo"
	AccessTokenSalt  = iota + 1
	RefreshTokenSalt
)

type (
	AccessTokenTicket struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	accessTokenConfig struct {
		AccessTokenKey   string        `mapStructure:"accessTokenKey"`
		AccessTokenTime  time.Duration `mapStructure:"accessTokenAt"`  //token持续时间
		RefreshTokenTime time.Duration `mapStructure:"refreshTokenAt"` //刷新token再token过期后多久有效
	}

	AccessToken interface {
		CreateAccessToken(ctx context.Context, uid uint64) (resp *AccessTokenTicket, err error)
		CheckAccessToken(ctx context.Context, token string) bool
		RefreshAccessToken(ctx context.Context, token string) (resp *AccessTokenTicket, err error)
		ClearAccessToken(ctx context.Context, uid uint64) (err error)
	}
)

func DefaultAccessTokenConfig() *accessTokenConfig {
	return &accessTokenConfig{
		AccessTokenKey:   DefaultAccessKey,
		AccessTokenTime:  time.Hour * 72,
		RefreshTokenTime: time.Hour * 24,
	}
}

var (
	AccessTokenCfg = xcfg.UnmarshalWithExpect("access.token", DefaultAccessTokenConfig()).(*accessTokenConfig)
)

func (a *AccessTokenTicket) Encode(uid uint64) (err error) {
	now := time.Now()
	accessTokenClaims := &jwt.StandardClaims{
		ExpiresAt: now.Add(AccessTokenCfg.AccessTokenTime).Unix(),
		Id:        strconv.FormatUint(uid, 10),
		IssuedAt:  now.Unix(),
		Issuer:    `NDisk_User`,
		NotBefore: now.Unix(),
		Subject:   `JWT`,
	}
	accessTokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	a.AccessToken, err = accessTokenWithClaims.SignedString([]byte(AccessTokenCfg.AccessTokenKey))

	refreshTokenClaims := &jwt.StandardClaims{
		ExpiresAt: now.Add(AccessTokenCfg.AccessTokenTime).Add(AccessTokenCfg.RefreshTokenTime).Unix(),
		Id:        strconv.FormatUint(uid, 10),
		IssuedAt:  now.Unix(),
		Issuer:    `NDisk_User`,
		NotBefore: now.Unix(),
		Subject:   `JWT`,
	}
	refreshTokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	a.RefreshToken, err = refreshTokenWithClaims.SignedString([]byte(AccessTokenCfg.AccessTokenKey))
	return err
}

func (a *AccessTokenTicket) Decode() (uid uint64, err error) {
	if a.RefreshToken == "" && a.AccessToken == "" {
		return 0, errors.New("nil")
	}
	var (
		token  *jwt.Token
		secret = func() jwt.Keyfunc {
			return func(token *jwt.Token) (interface{}, error) {
				return []byte(AccessTokenCfg.AccessTokenKey), nil
			}
		}
	)
	if a.AccessToken != "" {
		token, err = jwt.Parse(a.AccessToken, secret())

	}
	if a.RefreshToken != "" {
		token, err = jwt.Parse(a.RefreshToken, secret())
	}
	if err != nil {
		return
	}
	claimMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to StandardClaims")
		return
	}
	return xcast.ToUint64(claimMap["jti"]), err
}
