package server

import (
	"context"
	"log"
	"net"

	pb "peer/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Master struct {
	pb.UnimplementedHeartbeatServer
}

type Worker struct {
	conn *grpc.ClientConn
}

func (m *Master) Start(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHeartbeatServer(s, m)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (m *Master) UpHeartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	return &pb.HeartbeatReply{Msg: "ok"}, nil
}
