package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
)

const (
	GrpcChannelCapacity = 100
)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	MsgType int32
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	protos.UnimplementedN1N2MessageServer
	grpcChannel chan *GrpcMessage
	ServerIP    string
	ServerPort  string
}

// Initialize server with config data
// Receive data and send to grpcChannel
func NewGrpcServer(cfg config.GrpcServerInfoConfig) *GrpcServer {
	return &GrpcServer{
		grpcChannel: make(chan *GrpcMessage, GrpcChannelCapacity),
		ServerIP:    cfg.ServerIP,
		ServerPort:  cfg.ServerPort,
	}
}

// Start server
func (g *GrpcServer) Start() {
	servAddr := g.ServerIP + ":" + g.ServerPort
	lis, err := net.Listen("tcp", servAddr)
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	protos.RegisterN1N2MessageServer(server, g)
	klog.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Watch channel data
func (g *GrpcServer) WatchGrpcChannel() chan *GrpcMessage {
	return g.grpcChannel
}

// Server handler function for SendN1N2 call from client
func (g *GrpcServer) SendN1N2MessageTransferData(ctx context.Context, in *protos.N1N2MessageTransferDataRequest) (*protos.N1N2MessageTransferDataResponse, error) {
	n1n2Data := GrpcMessageInfo(in)
	go func() {
		// Send N1N2 data into channel
		g.grpcChannel <- &GrpcMessage{
			// Change msgType as per sm.go
			MsgType: 11,
			GrpcMsg: &n1n2Data,
		}
	}()
	return &protos.N1N2MessageTransferDataResponse{
		PduSessionId: in.PduSessionId,
	}, nil
}
