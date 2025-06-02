package grpcclient

import (
	"context"
	"nasMM/pkg/config"
	"nasMM/pkg/grpcCspNas/protos/create_sm_context_grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

type GrpcClient struct {
	ConnAddr string                                                  // CSP IP
	Client   create_sm_context_grpc.SendDataForCreateSmContextClient // grpc Client Conn
}

// Initialize with IP data
func NewGrpcClient(cfg config.GrpcCspServerConfig) *GrpcClient {
	clientAddr := cfg.CspServerIP + ":" + cfg.CspServerPort
	return &GrpcClient{
		ConnAddr: clientAddr,
	}
}

// Start client with dial
func (g *GrpcClient) Start() {
	// Get CSP pod IP address and then dial
	conn, err := grpc.Dial(g.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.Client = create_sm_context_grpc.NewSendDataForCreateSmContextClient(conn)
}

func (g *GrpcClient) SendDataForCreateSmContext(req *create_sm_context_grpc.CreateSmContextDataFromNasMod) *create_sm_context_grpc.CreateSmContextRespToNasMod {
	client := g.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// grpc protos generated function
	res, err := client.SendDataForCreateSmContext(ctx, req)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	//call appropriate function if there is no error
	klog.Infoln("Response from server: ", res)
	return res
}
