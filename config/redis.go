package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	ctx context.Context
	rdb *redis.Client
}

func (rc *RedisCache) GetClient() *redis.Client {
	return rc.rdb
}

// rc RedisCache
func (rc *RedisCache) GetObject(key string, dest any) (bool, error) {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rc *RedisCache) GetValue(key string) (string, bool, error) {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (rc *RedisCache) SetObject(key string, obj any, exp time.Duration) error {
	objInByte, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if err = rc.rdb.Set(rc.ctx, key, objInByte, exp).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) SetValue(key string, value string, exp time.Duration) error {
	return rc.rdb.Set(rc.ctx, key, value, exp).Err()
}

func (rc *RedisCache) RemoveKey(ck string) error {
	_, err := rc.rdb.Del(rc.ctx, string(ck)).Result()
	// if err == redis.Nil
	return err
}

func (rc *RedisCache) RemoveKeyWithCount(ck string) (int64, error) {
	return rc.rdb.Del(rc.ctx, ck).Result()
}

func (rc *RedisCache) RemoveKeysWithCount(cks []string) (int64, error) {
	return rc.rdb.Del(rc.ctx, cks...).Result()
}

func (rc *RedisCache) AddSet(setKey string, member string) error {
	if err := rc.rdb.SAdd(rc.ctx, setKey, member).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) GetSetMembers(setKey string) ([]string, error) {
	return rc.rdb.SMembers(rc.ctx, setKey).Result()
}

func (rc *RedisCache) RemoveSetMember(setKey string, member string) error {
	return rc.rdb.SRem(rc.ctx, setKey, member).Err()
}

func (rc *RedisCache) RemoveKeys(keys ...string) error {
	k := make([]string, 0, len(keys))
	for _, ck := range keys {
		k = append(k, ck)
	}
	_, err := rc.rdb.Del(rc.ctx, k...).Result()
	return err
}

func ConnectRedis(ctx context.Context) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0, // use default DB
		PoolSize: 100,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Fail to Connect Redis")
	}
	return &RedisCache{
		rdb: rdb,
		ctx: ctx,
	}
}
