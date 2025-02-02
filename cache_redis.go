package cacheredis

import (
	"context"
	"errors"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// RedisCache - an implementation of the CacheInterface to connect to Redis
type RedisCache struct {
	rdb                *redis.Client
	ctx                context.Context
	defExpireMilliSecs int
}

var (
	ErrKeyDoesNotExist error = errors.New(`key does not exist`)
)

// NewRedisCache create a RedisCache object
func NewRedisCache(address string, password string, db, milliSecsExpire int) *RedisCache {
	if strings.LastIndex(address, `:`) == -1 {
		address += ":6379"
	}
	rds := &RedisCache{
		rdb: redis.NewClient(
			&redis.Options{
				Addr:     address,
				Password: password, // no password set
				DB:       db,       // use default DB
			}),
		ctx:                context.Background(),
		defExpireMilliSecs: milliSecsExpire,
	}
	cmd := rds.rdb.Ping(rds.ctx)
	result, err := cmd.Result()
	if err != nil || result != "PONG" {
		return nil
	}
	return rds
}

// NewRedisCacheContext create a RedisCache object with context
func NewRedisCacheContext(ctx context.Context, address string, password string, db, milliSecsExpire int) *RedisCache {
	if strings.LastIndex(address, `:`) == -1 {
		address += ":6379"
	}
	rds := &RedisCache{
		rdb: redis.NewClient(
			&redis.Options{
				Addr:     address,
				Password: password, // no password set
				DB:       db,       // use default DB
			}),
		ctx:                ctx,
		defExpireMilliSecs: milliSecsExpire,
	}
	cmd := rds.rdb.Ping(rds.ctx)
	result, err := cmd.Result()
	if err != nil || result != "PONG" {
		return nil
	}
	return rds
}

// Set a value by key with an expiration in milliseconds
func (rc *RedisCache) Set(key string, value []byte, exp int) error {
	if exp == 0 {
		exp = rc.defExpireMilliSecs
	}
	return rc.rdb.Set(
		rc.ctx,
		key,
		value,
		time.Duration(exp)*time.Millisecond).Err()
}

// Get value by key
func (rc *RedisCache) Get(dst []byte, key string) []byte {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return []byte{}
	}
	return []byte(val)
}

// GetWithErr gets value by key that can indicate an error
func (rc *RedisCache) GetWithErr(key string) ([]byte, error) {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return []byte{}, ErrKeyDoesNotExist
	}

	return []byte(val), nil
}

// Del removes the value by key
func (rc *RedisCache) Del(keyPattern string) error {
	keys := rc.rdb.Keys(rc.ctx, keyPattern)
	for _, v := range keys.Val() {
		if err := rc.rdb.Del(rc.ctx, v).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Has returns true when the key exist
func (rc *RedisCache) Has(key string) bool {
	cnt := rc.rdb.Exists(rc.ctx, key)
	return cnt.Val() > 0
}

// Reset flushes all keys
func (rc *RedisCache) Reset() {
	rc.rdb.FlushAll(rc.ctx)
}

// Ping tests if the connection to server has succeeded
func (rc *RedisCache) Ping() (string, error) {
	return rc.rdb.Ping(rc.ctx).Result()
}

// ListKeys lists all keys
func (rc *RedisCache) ListKeys() []string {
	var (
		cursor  uint64
		allkeys []string
	)
	for {
		keys, cursor, err := rc.rdb.Scan(rc.ctx, cursor, "*", 10).Result()
		if err != nil {
			break
		}
		allkeys = append(allkeys, keys...)
		if cursor == 0 {
			break
		}
	}
	return allkeys
}
