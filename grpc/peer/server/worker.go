package server

import (
	"context"
	"log"

	pb "peer/pb"

	"google.golang.org/grpc"
)

func (w *Worker) connectTo(addr string) pb.HeartbeatClient {
	var err error

	w.conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewHeartbeatClient(w.conn)
}

func (w *Worker) callTo(c pb.HeartbeatClient, who string) {
	r, err := c.UpHeartbeat(context.Background(), &pb.HeartbeatRequest{Who: who})
	if err != nil {
		log.Fatalf("call replica rpc error: %v\n", err)
	}
	log.Printf("reply msg: %s\n", r.Msg)
}

func (w *Worker) Start(addr string) {
	client := w.connectTo(addr)
	w.callTo(client, "abc")
}

func (w *Worker) Destroy() {
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
}
