package warp

import (
	"context"
	"github.com/cskr/pubsub"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	UnimplementedControlServiceServer
	ps *pubsub.PubSub
}

func (s *server) UpdateLiveData(_ context.Context, req *UpdateLiveDataRequest) (*AcknowledgedResponse, error) {
	s.ps.Pub(req, "newData")
	return &AcknowledgedResponse{}, nil
}

func NewControlServiceServer(ps *pubsub.PubSub) {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	RegisterControlServiceServer(s, &server{ps: ps})
	log.Printf("ControlServiceServer listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
