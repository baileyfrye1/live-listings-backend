package session

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once sync.Once
	rdb  *redis.Client
)

func GetClient() *redis.Client {
	once.Do(
		func() {
			rdb = redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_ADDRESS"),
				Password: "",
				DB:       0,
				Protocol: 2,
			})
		},
	)
	return rdb
}
