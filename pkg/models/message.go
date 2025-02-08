package models

type UserMessage struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Content   string `json:"content"`
}

type SQSMessage struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	MessageID string `json:"message_id"`
}

type SQSResponse struct {
	SessionID string `json:"session_id"`
	Content   string `json:"content"`
	MessageID string `json:"message_id"`
}
