package redisToken

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/myxy99/ndisk/internal/nuser/server/token"
)

type key string

const (
	tokenKey   key = `disk_user_token_%d`
	refreshKey key = `disk_user_refreshToken_%d`
)

func (s key) Format(arg uint64) string {
	return fmt.Sprintf(string(s), arg)
}

func NewAccessToken(c *redis.Client) *AccessToken {
	return &AccessToken{
		c,
	}
}

type AccessToken struct {
	c *redis.Client
}

func (a *AccessToken) CreateAccessToken(ctx context.Context, uid uint64) (resp *token.AccessTokenTicket, err error) {
	resp = new(token.AccessTokenTicket)
	err = resp.Encode(uid)
	if err != nil {
		return nil, err
	}
	err = a.c.Set(ctx, tokenKey.Format(uid), resp.AccessToken, token.AccessTokenCfg.AccessTokenTime).Err()
	if err != nil {
		return nil, err
	}
	err = a.c.Set(ctx, refreshKey.Format(uid), resp.RefreshToken, token.AccessTokenCfg.AccessTokenTime+token.AccessTokenCfg.RefreshTokenTime).Err()
	if err != nil {
		return nil, err
	}
	return
}

func (a *AccessToken) CheckAccessToken(ctx context.Context, tokens string) bool {
	i, _ := new(token.AccessTokenTicket).Decode(tokens)
	if i.Uid <= 0 || i.Type != token.AccessTokenType {
		return false
	}
	rToken := a.c.Get(ctx, tokenKey.Format(i.Uid)).String()
	return rToken == tokens
}

func (a *AccessToken) RefreshAccessToken(ctx context.Context, tokens string) (resp *token.AccessTokenTicket, err error) {
	i, _ := new(token.AccessTokenTicket).Decode(tokens)
	if i.Uid <= 0 || i.Type != token.RefreshTokenType {
		return nil, token.DecryptErr
	}
	rToken := a.c.Get(ctx, refreshKey.Format(i.Uid)).Val()
	if rToken != tokens {
		return nil, token.DecryptErr
	}
	return a.CreateAccessToken(ctx, i.Uid)
}

func (a *AccessToken) ClearAccessToken(ctx context.Context, uid uint64) (err error) {
	return a.c.Del(ctx, refreshKey.Format(uid), tokenKey.Format(uid)).Err()
}

func (a *AccessToken) DecoderAccessToken(ctx context.Context, tokens string) (info *token.Info, err error) {
	return new(token.AccessTokenTicket).Decode(tokens)
}
