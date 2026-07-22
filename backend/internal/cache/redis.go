package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"cloudpulse/backend/internal/metrics"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) *RedisCache {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{
			Addr: redisURL,
		}
	}

	client := redis.NewClient(opts)
	
	if err := redisotel.InstrumentTracing(client); err != nil {
		slog.Error("Failed to instrument redis with otel", slog.String("error", err.Error()))
	}
	
	return &RedisCache{client: client}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("set").Observe(time.Since(start).Seconds()) }()

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		slog.Error("Redis Error", slog.String("operation", "set"), slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("get").Observe(time.Since(start).Seconds()) }()

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			metrics.CacheMisses.WithLabelValues(key).Inc()
			slog.Info("Cache Miss", slog.String("key", key))
		} else {
			slog.Error("Redis Error", slog.String("operation", "get"), slog.String("key", key), slog.String("error", err.Error()))
		}
		return err
	}
	metrics.CacheHits.WithLabelValues(key).Inc()
	slog.Info("Cache Hit", slog.String("key", key))
	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("delete").Observe(time.Since(start).Seconds()) }()

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		slog.Error("Redis Error", slog.String("operation", "delete"), slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// Queue methods

func (c *RedisCache) Enqueue(ctx context.Context, queueName string, value interface{}) error {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("enqueue").Observe(time.Since(start).Seconds()) }()

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.LPush(ctx, queueName, bytes).Err()
	if err == nil {
		metrics.QueueLength.WithLabelValues(queueName).Inc()
		slog.Info("Redis Queue Push", slog.String("queue", queueName))
	} else {
		slog.Error("Redis Error", slog.String("operation", "enqueue"), slog.String("queue", queueName), slog.String("error", err.Error()))
	}
	return err
}

func (c *RedisCache) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (string, error) {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("dequeue").Observe(time.Since(start).Seconds()) }()

	result, err := c.client.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		if err != redis.Nil {
			slog.Error("Redis Error", slog.String("operation", "dequeue"), slog.String("queue", queueName), slog.String("error", err.Error()))
		}
		return "", err
	}
	if len(result) < 2 {
		return "", redis.Nil 
	}
	metrics.QueueLength.WithLabelValues(queueName).Dec()
	slog.Info("Redis Queue Pop", slog.String("queue", queueName))
	return result[1], nil
}

// Rate Limiting

func (c *RedisCache) RateLimit(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	start := time.Now()
	defer func() { metrics.RedisLatency.WithLabelValues("ratelimit").Observe(time.Since(start).Seconds()) }()

	pipe := c.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if incr.Val() > limit {
		return false, nil // Limit exceeded
	}

	return true, nil
}
