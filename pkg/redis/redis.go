package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// TODO: `CredentialsProvider`を用いた本番を想定したクライアントの提供方法に修正する
func NewClient(host, port, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	})

	return rdb
}
