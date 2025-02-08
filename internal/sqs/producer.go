package sqs

import (
	"context"
	"encoding/json"

	"github.com/GoLogann/chat-service/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Producer struct {
	client   *sqs.Client
	queueURL string
}

func NewProducer(queueURL string) *Producer {
	return &Producer{
		queueURL: queueURL,
	}
}

func (p *Producer) SendMessage(ctx context.Context, msg models.SQSMessage) error {
	msgBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            aws.String(string(msgBody)),
		QueueUrl:               aws.String(p.queueURL),
		MessageGroupId:         aws.String(msg.SessionID),
		MessageDeduplicationId: aws.String(msg.MessageID),
	})
	return err
}
