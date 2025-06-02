package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/protos"
)

const (
	K8sDnsResolver string = "dns://10.96.0.10:53/"
)

type GrpcClient struct {
	connAddr string
	Client   protos.N1N2MessageClient
}

// Initialize with IP data
func NewGrpcClient() *GrpcClient {
	return &GrpcClient{
		connAddr: "127.0.0.1:50051",
	}
}

func (g *GrpcClient) Start() {
	// Get Pdusmsp address and then dial
	conn, err := grpc.Dial(K8sDnsResolver+g.connAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	client := protos.NewN1N2MessageClient(conn)
	g.Client = client
}

func (g *GrpcClient) SendN1N2MessageData(smContextID string, n1n2Data *protos.N1N2MessageTransferDataRequest) {
	client := g.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	r, err := client.SendN1N2MessageTransferData(ctx, n1n2Data)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	klog.Infof("Response from server: %v", r.Status)
}
