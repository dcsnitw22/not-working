package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"
)

const (
	GrpcChannelCapacity = 100
)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	MsgType string
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	create_sm_context_grpc.UnimplementedSendDataForCreateSmContextServer
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
	create_sm_context_grpc.RegisterSendDataForCreateSmContextServer(server, g)
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
func (g *GrpcServer) SendDataForCreateSmContext(ctx context.Context, in *create_sm_context_grpc.CreateSmContextDataFromNasMod) (*create_sm_context_grpc.CreateSmContextRespToNasMod, error) {
	createSmContextData := GrpcMessageInfo(in)
	go func() {
		// Send N1N2 data into channel
		g.grpcChannel <- &GrpcMessage{
			//msg type hardcoded for now
			MsgType: "create",
			GrpcMsg: &createSmContextData,
		}
	}()
	klog.Infof("create sm context data : %+v", createSmContextData)
	return &create_sm_context_grpc.CreateSmContextRespToNasMod{
		ErrMessage: "no error",
	}, nil
}
