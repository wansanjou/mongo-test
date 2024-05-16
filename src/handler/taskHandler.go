package handler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wansanjou/test-mongo/src/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) ListTasks(c *fiber.Ctx) error {
	search := c.Query("search")
	sortBy := c.Query("sortBy")

	tasks, err := h.taskService.ListTasks(search,sortBy)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get tasks"})
	}

	return c.Status(http.StatusOK).JSON(tasks)
}

func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
			log.Printf("Invalid ID: %v", id)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	log.Printf("Received ID: %v", id)

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
			log.Printf("Failed to get task by id: %v, error: %v", id, err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get task"})
	}

	return c.JSON(task)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task service.TaskResponse
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	response, err := h.taskService.CreateTask(task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var task service.TaskResponse
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("Received request to update task with ID: %v", id)

	response, err := h.taskService.UpdateTask(id, task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	log.Printf("Successfully updated task with ID: %v", id) 

	return c.JSON(response)
}
