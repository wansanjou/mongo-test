package service

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskResponse struct {
	ObjectID    primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title json:"title" validate:"required,max=100"`
	Description string             `bson:"description json:"description"`
	Date        primitive.Timestamp             `bson:"date json:"date" validate:"required"`
	Image       string             `bson:"image json:"image"`
	Status      string             `bson:"status json:"status" validate:"required,oneof=IN_PROGRESS COMPLETED"`
}

type TaskService interface {
	ListTasks(search, sortBy string) ([]TaskResponse, error)
	GetTaskByID(id string) (*TaskResponse, error)
	CreateTask(task TaskResponse) (*TaskResponse, error)
	UpdateTask(id string, task TaskResponse) (*TaskResponse, error)
}
