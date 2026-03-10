package grpc

import (
	"log"

	userpb "github.com/AntonRadchenko/project-protos/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(addr string) (userpb.UserServiceClient, *grpc.ClientConn, error) {
	// 1. grpc.Dial(addr, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to users service: %v", err)
		return nil, nil, err
	}

	// 2. userpb.NewUserServiceClient(conn)
	client := userpb.NewUserServiceClient(conn)

	// 3. вернуть client, conn, err
	return client, conn, nil
}