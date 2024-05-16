package service

import (
	"errors"
	"log"

	"github.com/wansanjou/test-mongo/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskService struct {
	task_repo models.TaskRepository
}

func NewTaskService(task_repo models.TaskRepository) TaskService {
	return &taskService{task_repo: task_repo}
}

func (s taskService) ListTasks(search, sortBy string) ([]TaskResponse, error) {
	tasks, err := s.task_repo.ListTasks(search, sortBy)
	if err != nil {
		return nil, err
	}

	task_responses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		task_response := TaskResponse{
			Title:       task.Title,
			Description: task.Description,
			Date: 			 task.Date,
			Image:       task.Image,
			Status:      task.Status,
		}
		task_responses = append(task_responses, task_response)
	}

	return task_responses, nil
}

func (s taskService) GetTaskByID(id string) (*TaskResponse, error) {
	task, err := s.task_repo.GetTaskByID(id)
	if err != nil {
		log.Printf("Failed to get task by id: %v, error: %v", id, err)
		return nil, err
	}

	task_response := TaskResponse{
		Title:       task.Title,
		Description: task.Description,
		Date:        task.Date,
		Image:       task.Image,
		Status:      task.Status,
	}

	return &task_response, nil
}

func (s taskService) CreateTask(task_res TaskResponse) (*TaskResponse, error) {
	if task_res.Status == "" {
		task_res.Status = "IN_PROGRESS"
	} else if task_res.Status != "IN_PROGRESS" && task_res.Status != "COMPLETED" {
		return nil, errors.New("Error data Status !!")
	}

	task := models.Task{
		ObjectID:    primitive.NewObjectID(),
		Title:       task_res.Title,
		Description: task_res.Description,
		Date:        task_res.Date, 
		Image:       task_res.Image,
		Status:		   task_res.Status,
	}

	insert_task, err := s.task_repo.CreateTask(task)
	if err != nil {
		return nil, err
	}

	response := TaskResponse{
		Title:       insert_task.Title,
		Description: insert_task.Description,
		Date:        task_res.Date,
		Image:       insert_task.Image,
		Status:      insert_task.Status,
	}

	return &response, nil
}


func (s taskService) UpdateTask(id string, task_res TaskResponse) (*TaskResponse, error) {
	if id == "" {
    log.Printf("Invalid ID: %v", id)
    return nil, errors.New("Invalid ID")
	}

	if task_res.Status != "IN_PROGRESS" && task_res.Status != "COMPLETED" {
		log.Printf("Invalid status: %v", task_res.Status)
		return nil, errors.New("Invalid status value")
	}

	task := models.Task{
		Title:       task_res.Title,
		Description: task_res.Description,
		Date:        task_res.Date,
		Image:       task_res.Image,
		Status:      task_res.Status,
	}

	update_task, err := s.task_repo.UpdateTask(id, task)
	if err != nil {
		log.Printf("Failed to update task with id: %v, error: %v", id, err) 
		return nil, err
	}

	response := TaskResponse{
		Title:       update_task.Title,
		Description: update_task.Description,
		Date:        update_task.Date,
		Image:       update_task.Image,
		Status:      update_task.Status,
	}

	return &response, nil
}
