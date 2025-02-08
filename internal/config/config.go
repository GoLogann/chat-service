package config

type Config struct {
	Redis RedisConfig
	AWS   AWSConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type AWSConfig struct {
	SQS SQSConfig
}

type SQSConfig struct {
	QueueURL         string
	ResponseQueueURL string
}

func LoadConfig(path string) Config {
	return Config{
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		AWS: AWSConfig{
			SQS: SQSConfig{
				QueueURL:         "https://sqs.us-east-1.amazonaws.com/123456789/chat-requests.fifo",
				ResponseQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/chat-responses.fifo",
			},
		},
	}
}
