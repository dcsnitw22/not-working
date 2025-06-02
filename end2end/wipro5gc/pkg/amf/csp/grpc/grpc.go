package grpc

import (
	"context"
	"time"

	//"net"
	"k8s.io/klog"

	//"w5gc.io/wipro5gcore/openapi"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcclient"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/ngapNas"
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
	WatchN1N2DataGrpcChannel() chan *ngapNas.DLRequest
	SendDataForCreateSmContext(context.Context, *create_sm_context_grpc.CreateSmContextDataFromNasMod) (*create_sm_context_grpc.CreateSmContextRespToNasMod, error)
	SendN1N2Data(*ngapNas.DLRequest) *ngapNas.DLResponse
}

type GrpcInfo struct {
	grpcStartTime time.Time
	GrpcClient    *grpcclient.GrpcClient
	GrpcServer    *grpcserver.GrpcServer
}

func NewGrpc(cfg config.CspConfig) Grpc {
	return &GrpcInfo{
		GrpcServer: grpcserver.NewGrpcServer(cfg.GrpcServerInfo),
		GrpcClient: grpcclient.NewGrpcClient(cfg.GrpcClientInfo),
	}
}

func (g *GrpcInfo) Start() {
	klog.Info("Starting csp grpc server")
	go g.GrpcServer.Start()
	g.GrpcClient.Start()
}

// Implemented functions =====

func (g *GrpcInfo) WatchGrpcChannel() chan *grpcserver.GrpcMessage {
	return g.GrpcServer.WatchGrpcChannel()
}

func (g *GrpcInfo) WatchN1N2DataGrpcChannel() chan *ngapNas.DLRequest {
	return g.GrpcClient.WatchN1N2DataChannel()
}

func (g *GrpcInfo) SendDataForCreateSmContext(ctx context.Context, in *create_sm_context_grpc.CreateSmContextDataFromNasMod) (*create_sm_context_grpc.CreateSmContextRespToNasMod, error) {
	return g.GrpcServer.SendDataForCreateSmContext(ctx, in)
}

func (g *GrpcInfo) SendN1N2Data(in *ngapNas.DLRequest) *ngapNas.DLResponse {
	return g.GrpcClient.N1N2DataTransfer(in)
}
