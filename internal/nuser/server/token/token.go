package token

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/myxy99/component/pkg/xjson"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/ndisk/pkg/aes"
	"time"
)

const (
	DefaultAccessKey = "ecol123og1ysK#xo"
	AccessTokenType  = iota + 1
	RefreshTokenType
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

	Info struct {
		Uid  uint64
		Type uint64
	}

	AccessToken interface {
		CreateAccessToken(ctx context.Context, uid uint64) (resp *AccessTokenTicket, err error)
		CheckAccessToken(ctx context.Context, token string) bool
		RefreshAccessToken(ctx context.Context, token string) (resp *AccessTokenTicket, err error)
		DecoderAccessToken(ctx context.Context, token string) (resp *Info, err error)
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
	Aes            = aes.NewAes([]byte(AccessTokenCfg.AccessTokenKey))
	DecryptErr     = errors.New("token error")
)

func (a *AccessTokenTicket) Encode(uid uint64) (err error) {
	accessTokenInfoB, err := xjson.Marshal(Info{
		uid,
		AccessTokenType,
	})
	if !errors.Is(err, nil) {
		return err
	}
	accessToken, err := Aes.Encrypt(accessTokenInfoB)
	if !errors.Is(err, nil) {
		return err
	}
	a.AccessToken = base64.StdEncoding.EncodeToString(accessToken)

	refreshTokenInfoB, err := xjson.Marshal(Info{
		uid,
		RefreshTokenType,
	})
	if !errors.Is(err, nil) {
		return err
	}
	refreshToken, err := Aes.Encrypt(refreshTokenInfoB)
	if !errors.Is(err, nil) {
		return err
	}
	a.RefreshToken = base64.StdEncoding.EncodeToString(refreshToken)
	return err
}

func (a *AccessTokenTicket) Decode(token string) (info *Info, err error) {
	tokenB, _ := base64.StdEncoding.DecodeString(token)
	if !errors.Is(err, nil) {
		return
	}
	infoB, err := Aes.Decrypt(tokenB)
	if !errors.Is(err, nil) {
		return
	}
	err = xjson.Unmarshal(infoB, info)
	return
}
