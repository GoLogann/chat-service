package redis

import (
	"context"
	"encoding/json"
	"github.com/GoLogann/chat-service/pkg/models"
	"log"

	"github.com/go-redis/redis/v8"
)

type PubSub struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func (c *Client) NewPubSub(channels ...string) *PubSub {
	return &PubSub{
		client: c.client,
		pubsub: c.client.Subscribe(context.Background(), channels...),
	}
}

func (ps *PubSub) PublishResponse(response models.SQSResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return ps.client.Publish(context.Background(), "responses:"+response.SessionID, data).Err()
}

func (ps *PubSub) ListenResponses(ctx context.Context, sessionID string, handler func(response models.SQSResponse)) {
	channel := ps.pubsub.Channel()
	for msg := range channel {
		var response models.SQSResponse
		if err := json.Unmarshal([]byte(msg.Payload), &response); err != nil {
			log.Printf("Error decoding PubSub message: %v", err)
			continue
		}
		handler(response)
	}
}
