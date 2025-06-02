package grpcNgap

import (
	"nasMM/pkg/config"

	grpcNgapclient "nasMM/pkg/grpcNgap/grpcclient"
	"nasMM/pkg/grpcNgap/protos/n1n2messagetransfer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

func CreateNgapClient() n1n2messagetransfer.N1N2DataTransferClient {
	klog.Infof("Creating GRPC client for NAS to Ngap")

	cfg, err := config.InitConfig()
	if err != nil {
		klog.Fatalf("Failed to fetch config: %v", err)
	}

	clientInfo := grpcNgapclient.NewGrpcClient(cfg.DLNgap)

	// Create a new gRPC client
	clientConn, err := grpc.NewClient(clientInfo.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Failed to connect to server: %v", err)
	}

	client := n1n2messagetransfer.NewN1N2DataTransferClient(clientConn)

	return client
}
