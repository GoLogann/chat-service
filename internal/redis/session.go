package redis

import (
	"context"
	"time"

	"github.com/GoLogann/chat-service/pkg/utils"
)

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (c *Client) CreateSession(userID string) (*Session, error) {
	sessionID := utils.GenerateUUID()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err := c.client.HSet(context.Background(), "session:"+sessionID,
		"user_id", userID,
		"created_at", session.CreatedAt.Format(time.RFC3339),
		"expires_at", session.ExpiresAt.Format(time.RFC3339),
	).Err()

	if err != nil {
		return nil, err
	}

	c.client.Expire(context.Background(), "session:"+sessionID, time.Hour)
	return session, nil
}

func (c *Client) GetSession(sessionID string) *Session {
	result, err := c.client.HGetAll(context.Background(), "session:"+sessionID).Result()
	if err != nil || len(result) == 0 {
		return nil
	}

	createdAt, _ := time.Parse(time.RFC3339, result["created_at"])
	expiresAt, _ := time.Parse(time.RFC3339, result["expires_at"])

	return &Session{
		ID:        sessionID,
		UserID:    result["user_id"],
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

func (s *Session) IsActive() bool {
	return time.Now().Before(s.ExpiresAt)
}
