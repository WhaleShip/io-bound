package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whaleship/io-bound/internal/domain"
)

type taskService struct {
	tasks   map[string]*domain.TaskStatus
	tasksMu sync.RWMutex
}

func NewTaskService() *taskService {
	return &taskService{
		tasks: make(map[string]*domain.TaskStatus),
	}
}

func (s *taskService) CreateTask() string {
	taskID := uuid.New().String()

	s.tasksMu.Lock()
	s.tasks[taskID] = &domain.TaskStatus{Status: "pending"}
	s.tasksMu.Unlock()

	go s.process(taskID)

	return taskID
}

func (s *taskService) GetTask(id string) (*domain.TaskStatus, error) {
	s.tasksMu.RLock()
	defer s.tasksMu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (s *taskService) process(id string) {
	s.tasksMu.Lock()
	s.tasks[id].Status = "in-progress"
	s.tasksMu.Unlock()

	time.Sleep(4 * time.Minute)

	s.tasksMu.Lock()
	s.tasks[id].Status = "completed"
	s.tasks[id].Result = fiber.Map{"message": fmt.Sprintf("Task %s completed successfully", id)}
	s.tasksMu.Unlock()
}
