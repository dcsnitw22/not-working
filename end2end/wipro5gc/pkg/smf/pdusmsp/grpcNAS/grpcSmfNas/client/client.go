package grpcSMclient

import (
	"context"
	"log"

	// "nasMain/grpcSmfNas/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/grpcSmfNas/pb"
)

// SMClient is the gRPC client for sending SM data
type SMClient struct {
	client pb.SmfNasClient
	conn   *grpc.ClientConn
}

// NewSMClient creates a new SMClient
func NewSMClient(serverAddr string) (*SMClient, error) {
	clientConn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewSmfNasClient(clientConn)
	return &SMClient{client: client, conn: clientConn}, nil
}

// Close closes the gRPC connection
func (c *SMClient) Close() error {
	return c.conn.Close()
}

// SendSMData sends SM data to the server
func (c *SMClient) SendSMData(ctx context.Context, nas *anypb.Any, reqType string) (nasResponse *anypb.Any, erro error) {
	req := &pb.SMDataRequest{
		NasMessage: nas,
		TypeReq:    reqType,
	}

	resp, err := c.client.SendSMData(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("Response from server: %s", resp)
	return resp.NasResponse, nil
}
