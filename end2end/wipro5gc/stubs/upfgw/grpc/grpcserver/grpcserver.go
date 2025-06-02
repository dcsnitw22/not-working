package grpcserver

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/protos"
	"w5gc.io/wipro5gcore/stubs/upfgw/sm"
)

const (
	GrpcChannelCapacity = 100
)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	MsgType sm.MessageType
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	grpcChannel chan *GrpcMessage
	protos.UnimplementedSendSmContextDataServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		grpcChannel: make(chan *GrpcMessage, GrpcChannelCapacity),
	}
}

func (g *GrpcServer) Start() {
	klog.Infof("Started upfgw gRPC server")
	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	protos.RegisterSendSmContextDataServer(server, g)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SmContext server handler function
func (g *GrpcServer) SendSmContextCreateData(ctx context.Context, in *protos.SmContextCreateDataRequest) (*protos.SmContextCreateDataResponse, error) {
	smContextID := GrpcMessageInfo(in.SmContextID)
	go func() {
		g.grpcChannel <- &GrpcMessage{
			MsgType: 13,
			GrpcMsg: &smContextID,
		}
	}()
	return &protos.SmContextCreateDataResponse{SmContextID: in.GetSmContextID()}, nil
}

// SmContext server handler function
func (g *GrpcServer) SendSmContextUpdateData(ctx context.Context, in *protos.SmContextUpdateDataRequest) (*protos.SmContextUpdateDataResponse, error) {
	pduSessionId := GrpcMessageInfo(in.PduSessionId)
	go func() {
		g.grpcChannel <- &GrpcMessage{
			MsgType: 13,
			GrpcMsg: &pduSessionId,
		}
	}()
	return &protos.SmContextUpdateDataResponse{PduSessionId: in.GetPduSessionId()}, nil
}

// SmContext server handler function
func (g *GrpcServer) SendSmContextReleaseData(ctx context.Context, in *protos.SmContextReleaseDataRequest) (*protos.SmContextReleaseDataResponse, error) {
	pduSessionId := GrpcMessageInfo(in.PduSessionId)
	go func() {
		g.grpcChannel <- &GrpcMessage{
			MsgType: 13,
			GrpcMsg: &pduSessionId,
		}
	}()
	return &protos.SmContextReleaseDataResponse{PduSessionId: in.GetPduSessionId()}, nil
}

// SmContext return channel data
func (g *GrpcServer) WatchGrpcChannel() chan *GrpcMessage {
	return g.grpcChannel
}
