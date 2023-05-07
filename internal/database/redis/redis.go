package redis

import (
	"application_template/internal/config"
	"application_template/internal/database/connect"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func New(conf config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		DB:       conf.RedisDB,
		Password: conf.RedisPassword,
	})
}

func Get(c *gin.Context, k string) (string, error) {
	rdb := connect.RedisDB

	cmd := rdb.Get(c, k)
	if cmd.Err() != nil {
		log.Printf("redis get %s failed: %s\n", k, cmd.Err())
		return "", cmd.Err()
	}

	result, err := cmd.Result()
	if err != nil {
		log.Printf("redis result failed: %s\n", err)
		return "", err
	}

	log.Println("returning from redis cache")
	return result, nil
}

func Set(c *gin.Context, k string, val interface{}) error {
	rdb := connect.RedisDB

	js, err := json.Marshal(gin.H{"data": val})
	if err != nil {
		log.Printf("marshal error: %s\n", err)
		return err
	}

	err = rdb.Set(c, k, js, 1*time.Hour).Err()
	if err != nil {
		log.Printf("redis set %s failed: %s\n", k, err)
		return err
	}

	return nil
}

func Unset(c *gin.Context, k string) error {
	rdb := connect.RedisDB

	if err := rdb.Del(c, k).Err(); err != nil {
		log.Printf("redis del %s failed: %s\n", k, err)
		return err
	}

	return nil
}
