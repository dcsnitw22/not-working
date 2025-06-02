package grpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
	// pb "w5gc.io/wipro5gcore/pkg/protos"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement SendData server.
type server struct {
	protos.UnimplementedSendSmContextDataServer
}

func (s *server) SendSmContextCreateData(ctx context.Context, in *protos.SmContextCreateDataRequest) (*protos.SmContextCreateDataResponse, error) {
	return &protos.SmContextCreateDataResponse{PduSessionId: in.GetPduSessionId()}, nil
}

func (s *server) SendSmContextUpdateData(ctx context.Context, in *protos.SmContextUpdateDataRequest) (*protos.SmContextUpdateDataResponse, error) {
	return &protos.SmContextUpdateDataResponse{PduSessionId: 200}, nil
}

func UpfGwServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	protos.RegisterSendSmContextDataServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
