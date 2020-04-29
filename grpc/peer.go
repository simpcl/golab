package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Master struct {
}

type Worker struct {
	conn *grpc.ClientConn
}

func (m *Master) start(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterHeartbeatServer(s, m)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (m *Master) UpHeartbeat(ctx context.Context, in *HeartbeatRequest) (*HeartbeatReply, error) {
	return &HeartbeatReply{Msg: "ok"}, nil
}

func (w *Worker) connectTo(addr string) HeartbeatClient {
	var err error

	w.conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return NewHeartbeatClient(w.conn)
}

func (w *Worker) callTo(c HeartbeatClient, who string) {
	r, err := c.UpHeartbeat(context.Background(), &HeartbeatRequest{Who: who})
	if err != nil {
		log.Fatalf("call replica rpc error: %v\n", err)
	}
	log.Printf("reply msg: %s\n", r.Msg)
}

func (w *Worker) start(addr string) {
	client := w.connectTo(addr)
	w.callTo(client, "abc")
}

func (w *Worker) destroy() {
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
}
