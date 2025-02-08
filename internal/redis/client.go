package redis

import (
	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewClient(cfg Config) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}
