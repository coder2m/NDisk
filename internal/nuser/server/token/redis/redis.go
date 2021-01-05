package redis

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
	t := new(token.AccessTokenTicket)
	t.AccessToken = tokens
	uid, _ := t.Decode()
	if uid <= 0 {
		return false
	}
	rToken := a.c.Get(ctx, tokenKey.Format(uid)).String()
	return rToken == tokens
}

func (a *AccessToken) RefreshAccessToken(ctx context.Context, tokens string) (resp *token.AccessTokenTicket, err error) {
	t := new(token.AccessTokenTicket)
	t.RefreshToken = tokens
	uid, _ := t.Decode()
	if uid <= 0 {
		return nil, err
	}
	rToken := a.c.Get(ctx, refreshKey.Format(uid)).String()
	if rToken != tokens {

	}
	return
}

func (a *AccessToken) ClearAccessToken(ctx context.Context, uid uint64) (err error) {
	panic("implement me")
}
