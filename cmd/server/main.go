package main

import (
	"log"

	"github.com/AntonRadchenko/tasks-service/internal/database"
	"github.com/AntonRadchenko/tasks-service/internal/task"
	"github.com/AntonRadchenko/tasks-service/internal/transport/grpc"
)

func main() {
	// 1. Инициализация БД
	database.InitDB()

	// 2. Репозиторий и сервис задач
	repo := &task.TaskRepo{}                   
	svc := task.NewTaskService(repo)              

	// 3. Клиент к Users-сервису
	userClient, conn, err := grpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to users service: %v", err)
	}
	defer conn.Close()

	// 4. Запуск gRPC Tasks-сервиса
	if err := grpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("Tasks gRPC server error: %v", err)
	}
}