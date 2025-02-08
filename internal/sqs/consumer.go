package sqs

import (
	"context"
	"encoding/json"
	"github.com/GoLogann/chat-service/internal/shared"
	"log"
	"time"

	"github.com/GoLogann/chat-service/internal/redis"
	"github.com/GoLogann/chat-service/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
	redis    *redis.Client
	sender   shared.ResponseSender
}

func NewConsumer(queueURL string, rc *redis.Client, sender shared.ResponseSender) *Consumer {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("AWS config error:", err)
	}

	return &Consumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
		redis:    rc,
		sender:   sender,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := c.pollMessages(ctx)
			if err != nil {
				log.Printf("Polling error: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (c *Consumer) pollMessages(ctx context.Context) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       20,
		VisibilityTimeout:     30,
		AttributeNames:        []types.QueueAttributeName{"All"},
		MessageAttributeNames: []string{"All"},
	}

	output, err := c.client.ReceiveMessage(ctx, input)
	if err != nil {
		return err
	}

	for _, msg := range output.Messages {
		c.processMessage(ctx, msg)
	}

	return nil
}

func (c *Consumer) processMessage(ctx context.Context, msg types.Message) {
	var response models.SQSResponse
	if err := json.Unmarshal([]byte(*msg.Body), &response); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	if err := c.redis.SetMessageCache(response.SessionID, response.Content, response); err != nil {
		log.Printf("Cache update error: %v", err)
	}

	c.sender.SendResponse(response.SessionID, response)

	if _, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	}); err != nil {
		log.Printf("Error deleting message: %v", err)
	}
}
