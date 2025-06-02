package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb"
)

const (
	K8sDnsResolver string = "dns://10.96.0.10:53/"
)

type GrpcClient struct {
	ConnAddr string               // AMF IP
	Client   pb.DataServiceClient // grpc Client Conn
}

// Initialize with IP data
func NewGrpcClient(cfg config.GrpcClientInfoConfig) *GrpcClient {
	clientAddr := cfg.ClientIP + ":" + cfg.ClientPort
	return &GrpcClient{
		ConnAddr: clientAddr,
	}
}

// Start client with dial
func (g *GrpcClient) Start() {
	// Get NAS pod IP address and then dial
	conn, err := grpc.NewClient(g.ConnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fmt.Println("Server address: ", conn.Target())
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.Client = pb.NewDataServiceClient(conn)
}

// which struct implements the SendNasPduFileClient interface?
// Send NAS PDU File Request to NAS pod IP address
// todo: modify
func (g *GrpcClient) SendData(req *pb.DataRequest) *pb.DataResponse { // *create_sm_context_grpc.CreateSmContextRespToNgapMod {
	client := g.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// grpc protos generated function
	res, err := client.SendData(ctx, req)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	//call appropriate function if there is no error
	klog.Infoln("Response from server: ", res)
	return res
}
