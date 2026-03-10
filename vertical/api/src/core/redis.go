package core

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		opts := &redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		}
		useTLS := strings.ToLower(os.Getenv("REDIS_USE_TLS")) == "true"
		if useTLS {
			opts.TLSConfig = &tls.Config{}
		}

		redisClient = redis.NewClient(opts)
	}
	return redisClient
}
