package grpc_test

import (
	"fmt"
	"log"
	"net"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	grpc "google.golang.org/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	grpcpkg "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
)

type server grpcserver.GrpcServer

var _ = Describe("Grpc", func() {
	var (
		// conn   *grpc.ClientConn
		// client protos.SendSmContextDataClient
		g grpcpkg.Grpc
		s *grpc.Server
	)
	go func() {
		lis, err := net.Listen("tcp", "127.0.0.1:50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s = grpc.NewServer()
		// pb.RegisterGreeterServer(s, &server{})
		protos.RegisterN1N2MessageServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// go func() {
	// 	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	// 	if err != nil {
	// 		log.Fatalf("failed to listen: %v", err)
	// 	}
	// 	s := grpc.NewServer()
	// 	// pb.RegisterGreeterServer(s, &server{})
	// 	protos.RegisterSendSmContextDataServer(s, &g.grpcServer)
	// 	log.Printf("server listening at %v", lis.Addr())
	// 	if err := s.Serve(lis); err != nil {
	// 		log.Fatalf("failed to serve: %v", err)
	// 	}
	// }()

	Context("Test Server", func() {
		// var err error
		// conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		// Expect(err).To(BeNil())
		// client = protos.NewSendSmContextDataClient(conn)

		It("Creates new GrpcInfo", func() {
			g = grpcpkg.NewGrpc(config.GrpcServerInfoConfig{}, config.GrpcClientInfoConfig{})
			Expect(fmt.Sprintf("%v", reflect.TypeOf(g))).To(Equal("*grpc.GrpcInfo"))
		})
		It("Starts Grcp server and client", func() {
			// g.Start()
		})
		It("Sends SmContextCreateData", func() {
			// customData := protos.SmContextCreateDataRequest{
			// 	Supi:         "supi11",
			// 	PduSessionId: 113,
			// }
			// g.SendSmContextCreateData(&customData)
		})
		It("Sends SmContextUpdateData", func() {

		})
		It("Sends SmContextReleaseData", func() {

		})
	})
})
