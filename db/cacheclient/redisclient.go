package cacheclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func (rc *RedisClient) Init() {
	redisAddr := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	log.Println("connecting to redis...")

	rc.client = redis.NewClient(&redis.Options{Addr: redisAddr})

	_, err := rc.client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("failed to ping to redis", err)
		return
	}
	log.Println("successfully connected and pinged redis! :)")
}

func (rc *RedisClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", customerrors.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (rc *RedisClient) Set(key, val string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rc.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) Exists(key string) (bool, error) {
	_, err := rc.Get(key)
	if err == customerrors.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (rc *RedisClient) Cleanup() {
	log.Println("closing redis...")
	rc.client.Close()
}
