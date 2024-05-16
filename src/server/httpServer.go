package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/wansanjou/test-mongo/src/config"
	"github.com/wansanjou/test-mongo/src/handler"

	"go.mongodb.org/mongo-driver/mongo"
)

type IHttpServer interface {
	StartTaskServer(taskHandler *handler.TaskHandler)
}

type httpServer struct {
	app      *fiber.App
	cfg      *config.Config
	dbClient *mongo.Client
}

func NewHttpServer(cfg *config.Config, dbClient *mongo.Client) IHttpServer {
	return &httpServer{
		app:      fiber.New(),
		cfg:      cfg,
		dbClient: dbClient,
	}
}

func (s *httpServer) Listen() {
	log.Printf("Server is starting on %s", s.cfg.App.Url)
	if err := s.app.Listen(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}

func (s *httpServer) gracefullyShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := s.app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}

func (s *httpServer) StartTaskServer(taskHandler *handler.TaskHandler) {
	s.app.Get("/task", taskHandler.ListTasks)
	s.app.Get("/task/:id", taskHandler.GetTaskByID)
	s.app.Post("/task", taskHandler.CreateTask)
	s.app.Put("/task/:id", taskHandler.UpdateTask)

	go s.Listen()
	s.gracefullyShutdown()
}
