package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/whaleship/io-bound/internal/domain"
)

type taskService interface {
	CreateTask(ctx context.Context) (string, error)
	GetTask(ctx context.Context, taskID string) (*domain.TaskStatus, error)
}
type taskHandler struct {
	taskSvc taskService
}

func NewTaskHandler(s taskService) *taskHandler {
	return &taskHandler{taskSvc: s}
}

func (h *taskHandler) HandleCreateTask(c *fiber.Ctx) error {
	taskID, err := h.taskSvc.CreateTask(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"task_id": taskID})
}

func (h *taskHandler) HandleGetTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	task, err := h.taskSvc.GetTask(c.UserContext(), taskID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(task)
}
