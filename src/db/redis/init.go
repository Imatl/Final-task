package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"supportflow/core"
)

var Client *redis.Client

func Init(ctx context.Context) error {
	db, _ := strconv.Atoi(core.GetString("db.redis.db", "0"))

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", core.GetString("db.redis.host", "localhost"), core.GetString("db.redis.port", "6379")),
		Password: core.GetString("db.redis.password", ""),
		DB:       db,
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis: ping: %w", err)
	}

	return nil
}

func Close() {
	if Client != nil {
		Client.Close()
	}
}
