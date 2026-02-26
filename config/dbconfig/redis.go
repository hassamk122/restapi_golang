package dbconfig

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	addr := getEnv("REDIS_ADDRESS", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")

	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis : %w", err))
	}

	log.Println("Sucessfully connected to redis")

	return rdb
}

func getEnv(key, defaulValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaulValue
}
