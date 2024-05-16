package main

import (
	"os"

	"github.com/wansanjou/test-mongo/pkg/db"
	"github.com/wansanjou/test-mongo/src/config"
	"github.com/wansanjou/test-mongo/src/handler"
	"github.com/wansanjou/test-mongo/src/repository"
	"github.com/wansanjou/test-mongo/src/server"
	"github.com/wansanjou/test-mongo/src/service"
)

func main() {
	cfg := config.NewConfig(func() string {
		if len(os.Args) > 1 {
			return os.Args[1]
		}
		return "./.env.http.task"
	}())

	dbClient := db.DBConn(cfg)

	httpServer := server.NewHttpServer(cfg, dbClient)

	taskRepository := repository.NewTaskRepository(dbClient)
	taskService := service.NewTaskService(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)

	httpServer.StartTaskServer(taskHandler)
}