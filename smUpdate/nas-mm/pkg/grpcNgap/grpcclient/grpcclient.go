package grpcNgapclient

import (
	"context"
	"nasMM/pkg/config"
	"nasMM/pkg/grpcNgap/protos/n1n2messagetransfer"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

type GrpcNgapClient struct {
	ConnAddr string
	Client   n1n2messagetransfer.N1N2DataTransferClient
}

// Initialize with IP data
func NewGrpcClient(cfg config.GrpcNgapDLServerConfig) *GrpcNgapClient {
	clientAddr := cfg.DLNgapServerIP + ":" + cfg.DLNgapServerPort
	return &GrpcNgapClient{
		ConnAddr: clientAddr,
	}
}

// Start client with dial
func (g *GrpcNgapClient) Start() {
	// Get CSP pod IP address and then dial
	conn, err := grpc.Dial(g.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.Client = n1n2messagetransfer.NewN1N2DataTransferClient(conn)
}

func (g *GrpcNgapClient) SendN1N2DataToNgap(req *n1n2messagetransfer.N1N2Data) *n1n2messagetransfer.Error {
	client := g.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// grpc protos generated function
	res, err := client.SendN1N2DataTransfer(ctx, req)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	//call appropriate function if there is no error
	klog.Infoln("Response from server: ", res)
	return res

}
