package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// type PostCache interface {
// 	Set(key string, value *Shoe)
// 	Get(key string) *Shoe
// }

type redisCache struct {
	host   string
	db     int
	expire time.Duration
}

func NewRedisCache(host string, db int, expire time.Duration) *redisCache {
	return &redisCache{
		host:   host,
		db:     db,
		expire: expire,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *Shoe) {
	client := cache.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client.Set(key, json, cache.expire*time.Second)
}

func (cache *redisCache) Get(key string) *Shoe {
	client := cache.getClient()
	value, err := client.Get(key).Result()
	if err != nil {
		return nil
	}
	shoe := Shoe{}

	err = json.Unmarshal([]byte(value), &shoe)
	if err != nil {
		panic(err)
	}
	return &shoe
}
