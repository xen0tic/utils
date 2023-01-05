package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func New(conf *redis.Options) *Redis {
	r := new(Redis)
	r.client = redis.NewClient(conf)
	r.ctx = context.Background()
	return r
}

func (r *Redis) SetContext(ctx2 context.Context) {
	r.ctx = ctx2
}

func (r *Redis) Set(key string, value interface{}, expiration ...int) (string, error) {
	exp := 0
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	s, e := r.client.Set(r.ctx, key, value, time.Duration(exp)).Result()
	if e != nil {
		return "", e
	}
	return s, nil
}

func (r *Redis) Get(key string) (string, error) {
	s, e := r.client.Get(r.ctx, key).Result()
	if e != nil {
		return "", e
	}
	return s, nil
}

func (r *Redis) Has(key string) (int64, error) {
	rr, e := r.client.Exists(r.ctx, key).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) Del(key string) (int64, error) {
	rr, e := r.client.Del(r.ctx, key).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) HSet(key string, value interface{}) (int64, error) {
	rr, e := r.client.HSet(r.ctx, key, value).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) HGet(key string, field string, result interface{}) error {
	
	res, err := r.client.HGet(r.ctx, key, field).Result()
	
	if err != nil {
		return err
	}
	
	_ = json.Unmarshal([]byte(res), result)
	
	return nil
}

func (r *Redis) HExists(key string, field string) bool {
	res := r.client.HExists(r.ctx, key, field)
	exist, _ := res.Result()
	return exist
}

func (r *Redis) HGetString(key string, field string) (string, error) {
	
	res, err := r.client.HGet(r.ctx, key, field).Result()
	
	if err != nil {
		return "", err
	}
	
	return res, nil
}

func (r *Redis) HDel(key string, fields ...string) (int64, error) {
	rr, e := r.client.HDel(r.ctx, key, fields...).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	rr, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return rr, nil
}

func (r *Redis) RPush(key string, value string) (int64, error) {
	rr, e := r.client.RPush(r.ctx, key, value).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) RPushArray(key string, value ...interface{}) (int64, error) {
	rr, e := r.client.RPush(r.ctx, key, value...).Result()
	if e != nil {
		return 0, e
	}
	return rr, nil
}

func (r *Redis) LRange(key string) ([]string, error) {
	return r.LRangeWithStartAndMAx(key, 0, -1)
}

func (r *Redis) LPos(key, element string) (int64, error) {
	res := r.client.LPos(r.ctx, key, element, redis.LPosArgs{})
	return res.Result()
}

func (r *Redis) LRangeWithStartAndMAx(key string, start int64, max int64) ([]string, error) {
	res, err := r.client.LRange(r.ctx, key, start, start+max).Result()
	if err != nil {
		return []string{}, err
	}
	return res, err
}

func (r *Redis) ListPopCount(key string, count int) ([]string, error) {
	return r.client.LPopCount(r.ctx, key, count).Result()
}

func (r *Redis) ListRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

func (r *Redis) LLen(key string) (int64, error) {
	res, err := r.client.LLen(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	res := r.client.Keys(r.ctx, pattern)
	return res.Result()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
