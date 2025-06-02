package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/ngapNas"
)

var N1N2Data = make(chan (*ngapNas.DLRequest))

type GrpcClient struct {
	ConnAddr string                    // NGAP IP
	Client   ngapNas.DataServiceClient // grpc Client Conn
}

// Initialize with IP data
func NewGrpcClient(cfgClient config.GrpcClientInfoConfig) *GrpcClient {
	clientAddr := cfgClient.ClientIP + ":" + cfgClient.ClientPort
	return &GrpcClient{
		ConnAddr: clientAddr,
	}
}

// Start client with dial
func (g *GrpcClient) Start() {
	// Get CSP pod IP address and then dial
	conn, err := grpc.NewClient(g.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.Client = ngapNas.NewDataServiceClient(conn)
}

func (g *GrpcClient) WatchN1N2DataChannel() chan *ngapNas.DLRequest {
	return N1N2Data
}

func (g *GrpcClient) N1N2DataTransfer(req *ngapNas.DLRequest) *ngapNas.DLResponse {
	client := g.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// grpc protos generated function
	res, err := client.HandleDownlink(ctx, req)
	if err != nil {
		klog.V(10)
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	//call appropriate function if there is no error
	klog.Infoln("Response from server: ", res)
	return res
}
