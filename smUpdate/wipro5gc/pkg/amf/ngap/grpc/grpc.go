package grpc

import (
	"context"
	"time"

	//"net"

	//"w5gc.io/wipro5gcore/openapi"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/grpcclient"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb"
)

// Grpc Interface will be used in pdusmsp.go as Grpc
type Grpc interface {
	Start()

	// WatchGrpcChannel() chan *grpcserver.GrpcMessage
	SendData(*pb.DataRequest) *pb.DataResponse

	N2InfoNgapEncoder(context.Context, *protos.N2Information) (*protos.EncodedN2Information, error)
}

// GrpcInfo struct will implement Grpc Interface
type GrpcInfo struct {
	grpcStartTime time.Time

	GrpcClient *grpcclient.GrpcClient // Client sends SmContextData requests
	GrpcServer *grpcserver.GrpcServer // Server receives N2 data from SMF
}

// Initialize with new data
func NewGrpc(cfg *config.NgapConfig) Grpc {
	return &GrpcInfo{
		GrpcClient: grpcclient.NewGrpcClient(cfg.GrpcClientInfo),
		GrpcServer: grpcserver.NewGrpcServer(cfg.GrpcServerInfo),
	}
}

// Start client and server
func (g *GrpcInfo) Start() {
	go g.GrpcServer.Start()
	g.GrpcClient.Start()
}

// Implemented functions =====

/*func (g *GrpcInfo) WatchGrpcChannel() chan *grpcserver.GrpcMessage {
	return g.GrpcServer.WatchGrpcChannel()
}*/

func (g *GrpcInfo) SendData(data *pb.DataRequest) *pb.DataResponse {
	return g.GrpcClient.SendData(data)
}

func (g *GrpcInfo) N2InfoNgapEncoder(ctx context.Context, in *protos.N2Information) (*protos.EncodedN2Information, error) {
	return g.GrpcServer.SendN2Info(ctx, in)
}
