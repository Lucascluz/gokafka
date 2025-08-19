package cache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/lucas/gokafka/shared/utils"
)

type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		client: redis.NewClient(&redis.Options{
			Addr: utils.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			DB:   0,
		}),
	}
}

func (tb *TokenBlacklist) BlacklistToken(tokenID string, expiration time.Duration) error {
	return tb.client.Set("blacklist:"+tokenID, "revoked", expiration).Err()
}

func (tb *TokenBlacklist) IsTokenBlacklisted(tokenID string) bool {
	result := tb.client.Get("blacklist:" + tokenID)
	return result.Err() == nil
}
