package grpccspnas

import (
	// "context"
	"nasMM/pkg/config"
	"nasMM/pkg/grpcCspNas/grpcclient"
	"nasMM/pkg/grpcCspNas/protos/create_sm_context_grpc"

	// "time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

func CreateNASCSPClient() create_sm_context_grpc.SendDataForCreateSmContextClient { //, context.Context) {
	klog.Infof("Creating GRPC client for NAS to CSP")

	cfg, err := config.InitConfig()
	if err != nil {
		klog.Fatalf("Failed to fetch config: %v", err)
	}

	clientInfo := grpcclient.NewGrpcClient(cfg.CspGrpc)

	// Create a new gRPC client
	clientConn, err := grpc.NewClient(clientInfo.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Failed to connect to server: %v", err)
	}
	//	defer clientConn.Close()

	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()

	client := create_sm_context_grpc.NewSendDataForCreateSmContextClient(clientConn)

	return client //, ctx

}
