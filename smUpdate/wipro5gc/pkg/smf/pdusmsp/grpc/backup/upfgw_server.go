package grpc

import (
	"flag"
	// pb "w5gc.io/wipro5gcore/pkg/protos"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// // server is used to implement SendData server.
// type server struct {
// 	UnimplementedSendSmContextDataServer
// }

// func (s *server) SendSmContextCreateData(ctx context.Context, in *SmContextCreateDataRequest) (*SmContextCreateDataResponse, error) {
// 	return &SmContextCreateDataResponse{PduSessionId: in.GetPduSessionId()}, nil
// }

// func (s *server) SendSmContextUpdateData(ctx context.Context, in *SmContextUpdateDataRequest) (*SmContextUpdateDataResponse, error) {
// 	return &SmContextUpdateDataResponse{Resp: 200}, nil
// }

// func UpfGwServer() {
// 	flag.Parse()
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	// pb.RegisterGreeterServer(s, &server{})
// 	RegisterSendSmContextDataServer(s, &server{})
// 	log.Printf("server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
