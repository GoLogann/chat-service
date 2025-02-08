package shared

import (
	"context"
	"github.com/GoLogann/chat-service/pkg/models"
)

type ResponseSender interface {
	SendResponse(sessionID string, response models.SQSResponse)
}

type MessageSender interface {
	SendMessage(ctx context.Context, msg models.SQSMessage) error
}
