package token

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/coder2z/component/xcfg"
	"github.com/coder2z/g-saber/xjson"
	"github.com/coder2z/ndisk/pkg/aes"
)

const (
	AccessTokenType = iota + 1
	RefreshTokenType
)

const DefaultAccessKey = "ecol123og1ysK#xo"

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
		Time int64
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
		time.Now().Unix(),
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
		time.Now().Unix(),
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
	info = new(Info)
	tokenB, _ := base64.StdEncoding.DecodeString(token)
	if !errors.Is(err, nil) || len(tokenB) < 32 {
		return
	}
	infoB, err := Aes.Decrypt(tokenB)
	if !errors.Is(err, nil) {
		return
	}
	err = xjson.Unmarshal(infoB, info)
	return
}
