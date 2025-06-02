package grpcnas

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/grpcSmfNas/pb"
)

const (
	//K8sDnsResolver string = "dns://10.96.0.10:53/"
	serverAddr string = "grpcnassmf-service.nassm.svc.cluster.local:50052"
)

// func CreateGRPCNasClient() (pb.SmfNasClient, context.Context) {
func CreateGRPCNasClient() pb.SmfNasClient {
	// Create a new gRPC client
	clientConn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Failed to connect to server: %v", err)
	}
	// defer clientConn.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	client := pb.NewSmfNasClient(clientConn)
	// return client, ctx
	return client
}
