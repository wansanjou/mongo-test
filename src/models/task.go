package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Task struct {
	ObjectID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title string `bson:"title" json:"title "`
	Description string `bson:"description" json:"description "`
	Date primitive.Timestamp `bson:"date" json:"date "`
	Image string `bson:"image" json:"image "`
	Status string `bson:"status" json:"status "`
}

type TaskRepository interface {
	ListTasks(search, sortBy string) ([]Task, error)
	GetTaskByID(id string) (*Task, error)
	CreateTask(task Task) (*Task, error)
	UpdateTask(id string, task Task) (*Task, error)
}
