package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/whaleship/io-bound/internal/domain"
)

type taskService struct {
	redConn *redis.Client
}

func NewTaskService(client *redis.Client) *taskService {
	return &taskService{redConn: client}
}

func (s *taskService) CreateTask(ctx context.Context) (string, error) {
	taskID := uuid.New().String()
	hkey := fmt.Sprintf("task:%s", taskID)

	err := s.redConn.HSet(ctx, hkey, "status", "pending").Err()
	if err != nil {
		return "", err
	}

	go s.process(ctx, taskID)

	return taskID, nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string) (*domain.TaskStatus, error) {
	hkey := fmt.Sprintf("task:%s", taskID)
	vals, err := s.redConn.HGetAll(ctx, hkey).Result()
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	status := vals["status"]
	resultJSON := vals["result"]

	var result interface{}
	if resultJSON != "" {
		if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
			log.Println(err)
		}
	}
	return &domain.TaskStatus{Status: status, Result: result}, nil
}

type taskResult struct {
	Message string `json:"message"`
}

func (s *taskService) process(ctx context.Context, taskID string) {
	hkey := fmt.Sprintf("task:%s", taskID)
	s.redConn.HSet(ctx, hkey, "status", "in-progress")

	time.Sleep(4 * time.Minute)

	res := taskResult{Message: fmt.Sprintf("Task %s completed successfully", taskID)}
	resJSON, _ := json.Marshal(res)
	s.redConn.HMSet(ctx, hkey,
		"status", "completed",
		"result", string(resJSON),
	)
}
