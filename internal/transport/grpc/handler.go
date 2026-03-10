package grpc

import (
	"context"
	"fmt"

	taskpb "github.com/AntonRadchenko/project-protos/proto/task"
	userpb "github.com/AntonRadchenko/project-protos/proto/user"
	"github.com/AntonRadchenko/tasks-service/internal/task"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc        *task.TaskService
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.TaskService, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.Task, error) {
	// 1. Проверить пользователя:
	_, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId})
	if err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	params := task.CreateTaskParams{
		Task:   req.Task,
		IsDone: req.IsDone,
		UserID: uint(req.UserId),
	}

	// 2. Внутренняя логика
	createdTask, err := h.svc.CreateTask(params)
	if err != nil {
		return nil, err
	}

	// 3. Ответ
	return &taskpb.Task{
		Id:     uint32(createdTask.ID),
		Task:   createdTask.Task,
		IsDone: createdTask.IsDone,
		UserId: uint32(createdTask.UserID),
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	// Вызвать svc.GetTaskByID(req.Id)
	task, err := h.svc.GetTaskByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	// Вернуть &taskpb.Task{...}
	return &taskpb.Task{
		Id:     uint32(task.ID),
		Task:   task.Task,
		IsDone: task.IsDone,
		UserId: uint32(task.UserID),
	}, nil
}

func (h *Handler) ListTasks(ctx context.Context, _ *emptypb.Empty) (*taskpb.TaskList, error) {
	// Вызвать svc.GetTasks()
	tasks, err := h.svc.GetTasks()
	if err != nil {
		return nil, err
	}

	// Преобразовать срез task.Task → []*taskpb.Task
	pbTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Task:   t.Task,
			IsDone: t.IsDone,
			UserId: uint32(t.UserID),
		})
	}

	// Вернуть &taskpb.TaskList{Tasks: ...}
	return &taskpb.TaskList{Tasks: pbTasks}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.TaskList, error) {
	// Проверить пользователя через userClient.GetUser()
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	// Вызвать svc.GetTasksByUser(req.UserId)
	tasks, err := h.svc.GetTasksByUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	// Преобразовать срез task.Task в []*taskpb.Task
	protoTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		protoTasks = append(protoTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Task:   t.Task,
			IsDone: t.IsDone,
			UserId: uint32(t.UserID),
		})
	}

	return &taskpb.TaskList{Tasks: protoTasks}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.Task, error) {
	// Если обновляется user_id, проверить нового пользователя через userClient.GetUser()
	if req.UserId != nil {
		if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: *req.UserId}); err != nil {
			return nil, fmt.Errorf("user %d not found: %w", *req.UserId, err)
		}
	}

	// параметры
	var userID *uint
	if req.UserId != nil {
		uid := uint(*req.UserId)
		userID = &uid
	}

	params := task.UpdateTaskParams{
		Task:   req.Task,
		IsDone: req.IsDone,
		UserID: userID,
	}

	// Вызвать svc.UpdateTask(req.Id, params)
	updatedTask, err := h.svc.UpdateTask(uint(req.Id), params)
	if err != nil {
		return nil, err
	}

	return &taskpb.Task{
		Id:     uint32(updatedTask.ID),
		Task:   updatedTask.Task,
		IsDone: updatedTask.IsDone,
		UserId: uint32(updatedTask.UserID),
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*emptypb.Empty, error) {
	// Вызвать svc.DeleteTask(req.Id)
	err := h.svc.DeleteTask(uint(req.Id))
	if err != nil {
		return nil, err
	}

	// Вернуть &emptypb.Empty{}
	return &emptypb.Empty{}, nil
}