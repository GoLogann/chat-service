package redis

import (
	"context"
	"encoding/json"
	"github.com/GoLogann/chat-service/pkg/models"
	"time"
)

const cacheTTL = 5 * time.Minute

func (c *Client) GetMessageCache(sessionID, message string) *models.SQSResponse {
	data, err := c.client.Get(context.Background(), cacheKey(sessionID, message)).Bytes()
	if err != nil {
		return nil
	}

	var response models.SQSResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil
	}
	return &response
}

func (c *Client) SetMessageCache(sessionID string, content string, response models.SQSResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), cacheKey(sessionID, content), data, cacheTTL).Err()
}

func cacheKey(sessionID, message string) string {
	return "cache:" + sessionID + ":" + message
}
