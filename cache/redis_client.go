package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Get(key string) (string, error) {

	val, err := client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return val, err

}

func Set(key string, value string) {

	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}

}
