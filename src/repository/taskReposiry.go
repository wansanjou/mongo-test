package repository

import (
	"context"
	"log"

	"github.com/wansanjou/test-mongo/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	Client *mongo.Client
}

func NewTaskRepository(Client *mongo.Client) models.TaskRepository {
	return &TaskRepository{Client: Client} 
}


func (t *TaskRepository) ListTasks(search, sortBy string) ([]models.Task, error) {
	tasks := []models.Task{}

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
			{"status": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	options := options.Find()
	switch sortBy {
	case "Title":
		options.SetSort(bson.D{{Key: "title", Value: 1}})
	case "Date":
		options.SetSort(bson.D{{Key: "date", Value: 1}})
	case "Status":
		options.SetSort(bson.D{{Key: "status", Value: 1}})
	}

	collection := t.Client.Database("task_db").Collection("tasks")
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}



func (t *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
			log.Printf("Failed to convert id: %v, error: %v", id, err)
			return nil, err
	}

	task := models.Task{}
	filter := bson.M{"_id": objectID}

	collection := t.Client.Database("task_db").Collection("tasks")
	err = collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
			log.Printf("Failed to find task with filter: %v, error: %v", filter, err)
			return nil, err
	}

	return &task, nil
}

func (t *TaskRepository) CreateTask(task models.Task) (*models.Task, error) {
	collection := t.Client.Database("task_db").Collection("tasks")
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *TaskRepository) UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
			log.Printf("Failed to convert id: %v, error: %v", id, err)
			return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updatedTask}

	collection := t.Client.Database("task_db").Collection("tasks")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
			log.Printf("Failed to update task with id: %v, error: %v", id, err)
			return nil, err
	}

	// Find the updated task to return
	err = collection.FindOne(context.Background(), filter).Decode(&updatedTask)
	if err != nil {
			log.Printf("Failed to retrieve updated task with id: %v, error: %v", id, err)
			return nil, err
	}

	return &updatedTask, nil
}


