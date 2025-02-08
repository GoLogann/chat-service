package main

import (
	"context"
	"log"
	"net/http"

	"github.com/GoLogann/chat-service/internal/config"
	"github.com/GoLogann/chat-service/internal/redis"
	"github.com/GoLogann/chat-service/internal/sqs"
	"github.com/GoLogann/chat-service/internal/websocket"
	"github.com/GoLogann/chat-service/pkg/workerpool"
)

func main() {
	cfg := config.LoadConfig("config.yaml")

	redisClient := redis.NewClient(redis.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(redisClient)

	wsManager := websocket.NewManager(
		redisClient,
		sqs.NewProducer(cfg.AWS.SQS.QueueURL),
		workerpool.NewWorkerPool(10),
	)

	sqsConsumer := sqs.NewConsumer(
		cfg.AWS.SQS.ResponseQueueURL,
		redisClient,
		wsManager,
	)

	go func() {
		if err := sqsConsumer.Start(context.Background()); err != nil {
			log.Fatal("SQS Consumer failed:", err)
		}
	}()

	http.HandleFunc("/ws", wsManager.HandleConnections)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
