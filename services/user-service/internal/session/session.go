package session

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type SessionData struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(addr, password string, db int) *RedisSessionStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test the connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	log.Println("Connected to Redis")

	return &RedisSessionStore{client: rdb}
}

func (r *RedisSessionStore) CreateSession(userID, email, role string) (string, error) {
	sessionID := uuid.New().String()
	sessionData := SessionData{
		UserID:    userID,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
	}

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	// Store session for 24 hours
	err = r.client.Set(context.Background(), "session:"+sessionID, sessionJSON, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (r *RedisSessionStore) GetSession(sessionID string) (*SessionData, error) {
	sessionJSON, err := r.client.Get(context.Background(), "session:"+sessionID).Result()
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	err = json.Unmarshal([]byte(sessionJSON), &sessionData)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (r *RedisSessionStore) DeleteSession(sessionID string) error {
	return r.client.Del(context.Background(), "session:"+sessionID).Err()
}
