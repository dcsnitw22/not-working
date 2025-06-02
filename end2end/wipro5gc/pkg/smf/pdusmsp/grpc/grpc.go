package grpc

import (
	"time"

	//"net"

	//"w5gc.io/wipro5gcore/openapi"

	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/grpcclient"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
)

// Grpc Interface will be used in pdusmsp.go as Grpc
type Grpc interface {
	Start()

	WatchGrpcChannel() chan *grpcserver.GrpcMessage
	SendSmContextCreateData(*protos.SmContextCreateDataRequest)
	SendSmContextUpdateData(*protos.SmContextUpdateDataRequest)
	SendSmContextReleaseData(*protos.SmContextReleaseDataRequest)
	SendN2Info(*protos.N2Information) *protos.EncodedN2Information
}

// GrpcInfo struct will implement Grpc Interface
type GrpcInfo struct {
	grpcStartTime time.Time

	GrpcClient *grpcclient.GrpcClient // Client sends SmContextData requests
	GrpcServer *grpcserver.GrpcServer // Server receives N1N2 Data from upfgw
}

// Initialize with new data
func NewGrpc(cfgServer config.GrpcServerInfoConfig, cfgClient config.GrpcClientInfoConfig) Grpc {
	return &GrpcInfo{
		GrpcClient: grpcclient.NewGrpcClient(cfgClient),
		GrpcServer: grpcserver.NewGrpcServer(cfgServer),
	}
}

// Start client and server
func (g *GrpcInfo) Start() {
	go g.GrpcServer.Start()
	g.GrpcClient.Start()
}

// Implemented functions =====

func (g *GrpcInfo) WatchGrpcChannel() chan *grpcserver.GrpcMessage {
	return g.GrpcServer.WatchGrpcChannel()
}

func (g *GrpcInfo) SendSmContextCreateData(createData *protos.SmContextCreateDataRequest) {
	g.GrpcClient.SendSmContextCreateData(createData)
}

func (g *GrpcInfo) SendSmContextUpdateData(updateData *protos.SmContextUpdateDataRequest) {
	g.GrpcClient.SendSmContextUpdateData(updateData)
}

func (g *GrpcInfo) SendSmContextReleaseData(releaseData *protos.SmContextReleaseDataRequest) {
	g.GrpcClient.SendSmContextReleaseData(releaseData)
}

func (g *GrpcInfo) SendN2Info(n2Info *protos.N2Information) *protos.EncodedN2Information {
	return g.GrpcClient.SendN2Info(n2Info)
}
