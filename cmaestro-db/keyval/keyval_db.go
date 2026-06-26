package keyval

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func New(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Close closes the Redis connection.
func (r *Client) Close() error {
	return r.client.Close()
}

// Set stores a string value.
func (r *Client) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get retrieves a string value.
func (r *Client) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Delete removes a key.
func (r *Client) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists.
func (r *Client) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// SetJSON stores a struct/object as JSON.
func (r *Client) SetJSON(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// GetJSON retrieves JSON into a struct/object.
func (r *Client) GetJSON(key string, dest any) error {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Increment increments an integer key.
func (r *Client) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Expire sets expiration for a key.
func (r *Client) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// FlushDB clears the current database.
func (r *Client) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}
