package grpc

import (
	"context"
	"time"

	//"net"
	"k8s.io/klog"

	//"w5gc.io/wipro5gcore/openapi"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"
)

/*const (
	GrpcChannelCapacity = 100
)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	grpcMsg *GrpcMessageInfo
}
*/

type Grpc interface {
	Start()
	WatchGrpcChannel() chan *grpcserver.GrpcMessage
	SendDataForCreateSmContext(context.Context, *create_sm_context_grpc.CreateSmContextDataFromNasMod) (*create_sm_context_grpc.CreateSmContextRespToNasMod, error)
}

type GrpcInfo struct {
	grpcStartTime time.Time
	GrpcServer    *grpcserver.GrpcServer
}

func NewGrpc(cfgServer config.GrpcServerInfoConfig) Grpc {
	return &GrpcInfo{
		GrpcServer: grpcserver.NewGrpcServer(cfgServer),
	}
}

func (g *GrpcInfo) Start() {
	klog.Info("Starting csp grpc server")
	go g.GrpcServer.Start()
}

// Implemented functions =====

func (g *GrpcInfo) WatchGrpcChannel() chan *grpcserver.GrpcMessage {
	return g.GrpcServer.WatchGrpcChannel()
}

func (g *GrpcInfo) SendDataForCreateSmContext(ctx context.Context, in *create_sm_context_grpc.CreateSmContextDataFromNasMod) (*create_sm_context_grpc.CreateSmContextRespToNasMod, error) {
	return g.GrpcServer.SendDataForCreateSmContext(ctx, in)
}
