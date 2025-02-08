package websocket

import (
	"context"
	"encoding/json"
	"github.com/GoLogann/chat-service/internal/redis"
	"github.com/GoLogann/chat-service/internal/shared"
	"log"
	"net/http"
	"sync"

	"github.com/GoLogann/chat-service/pkg/models"
	"github.com/GoLogann/chat-service/pkg/utils"
	"github.com/GoLogann/chat-service/pkg/workerpool"
	"github.com/gorilla/websocket"
)

type Manager struct {
	upgrader    websocket.Upgrader
	redisClient *redis.Client
	sqsProducer shared.MessageSender
	workerPool  *workerpool.WorkerPool
	connections map[string]*websocket.Conn
	mu          sync.Mutex
}

func NewManager(rc *redis.Client, sp shared.MessageSender, wp *workerpool.WorkerPool) *Manager {
	return &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		redisClient: rc,
		sqsProducer: sp,
		workerPool:  wp,
		connections: make(map[string]*websocket.Conn),
	}
}

func (m *Manager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		m.workerPool.Submit(func() {
			m.processMessage(conn, msg)
		})
	}
}

func (m *Manager) processMessage(conn *websocket.Conn, rawMsg []byte) {
	var userMsg models.UserMessage
	if err := json.Unmarshal(rawMsg, &userMsg); err != nil {
		log.Printf("Error decoding message: %v", err)
		return
	}

	session, err := m.ensureSession(userMsg.UserID, userMsg.SessionID)
	if err != nil {
		log.Printf("Session error: %v", err)
		return
	}

	if cached := m.redisClient.GetMessageCache(session.ID, userMsg.Content); cached != nil {
		if err := conn.WriteJSON(cached); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
		return
	}

	sqsMsg := models.SQSMessage{
		SessionID: session.ID,
		UserID:    userMsg.UserID,
		Content:   userMsg.Content,
		MessageID: utils.GenerateUUID(),
	}

	if err := m.sqsProducer.SendMessage(context.Background(), sqsMsg); err != nil {
		log.Printf("SQS send error: %v", err)
	}
}

func (m *Manager) ensureSession(userID, sessionID string) (*redis.Session, error) {
	if sessionID != "" {
		if session := m.redisClient.GetSession(sessionID); session != nil && session.IsActive() {
			return session, nil
		}
	}
	return m.redisClient.CreateSession(userID)
}

func (m *Manager) SendResponse(sessionID string, response models.SQSResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, ok := m.connections[sessionID]; ok {
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}
