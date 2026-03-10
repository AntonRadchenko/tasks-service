package grpc

import (
	"net"

	taskpb "github.com/AntonRadchenko/project-protos/proto/task"
	userpb "github.com/AntonRadchenko/project-protos/proto/user"
	"github.com/AntonRadchenko/tasks-service/internal/task"
	"google.golang.org/grpc"
)

func RunGRPC(svc *task.TaskService, userClient userpb.UserServiceClient) error {
	// 1. net.Listen на ":50052"
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}

	// 2. grpc.NewServer()
	grpcServer := grpc.NewServer()

	// 3. taskpb.RegisterTaskServiceServer(grpcServer, NewHandler(svc, userClient))
	handler := NewHandler(svc, userClient)
	taskpb.RegisterTaskServiceServer(grpcServer, handler)

	// 4. grpcServer.Serve(listener) (блокируется)
	return grpcServer.Serve(listener)
}