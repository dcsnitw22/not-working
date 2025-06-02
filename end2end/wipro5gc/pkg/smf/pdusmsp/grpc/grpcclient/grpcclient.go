package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
)

const (
	K8sDnsResolver string = "dns://10.96.0.10:53/"
)

type GrpcClient struct {
	ConnAddrUpfgw string                         // Upfgw IP
	ClientUpfgw   protos.SendSmContextDataClient // Upfgw grpc Client Conn
	ConnAddrNgap  string                         // Ngap IP
	ClientNgap    protos.N2InfoNgapEncoderClient // Ngap grpc Client Conn
}

// Initialize with IP data
func NewGrpcClient(cfg config.GrpcClientInfoConfig) *GrpcClient {
	clientAddrUpfgw := cfg.ClientIPUpfgw + ":" + cfg.ClientPortUpfgw
	clientAddrNgap := cfg.ClientIPNgap + ":" + cfg.ClientPortNgap
	return &GrpcClient{
		ConnAddrUpfgw: clientAddrUpfgw,
		ConnAddrNgap:  clientAddrNgap,
	}
}

// Start client with dial
func (g *GrpcClient) Start() {
	// Get UPFGW address and then dial
	conn, err := grpc.Dial(K8sDnsResolver+g.ConnAddrUpfgw, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.ClientUpfgw = protos.NewSendSmContextDataClient(conn)

	// Get NGAP address and then dial
	conn, err = grpc.Dial(K8sDnsResolver+g.ConnAddrNgap, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Dial error. Did not connect: %v", err)
	}
	g.ClientNgap = protos.NewN2InfoNgapEncoderClient(conn)
}

// Send SmCreate data request to Upfgw
func (g *GrpcClient) SendSmContextCreateData(createData *protos.SmContextCreateDataRequest) {
	client := g.ClientUpfgw
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// createData := &SmContextCreateDataRequest{
	// 	Pei:          "Pei",
	// 	Dnn:          "Dnn",
	// 	PduSessionId: 11,
	// }

	// grpc protos generated function
	r, err := client.SendSmContextCreateData(ctx, createData)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	klog.Infof("Response from server: %d", r.GetSmContextID())
}

// Send Update data
func (g *GrpcClient) SendSmContextUpdateData(updateData *protos.SmContextUpdateDataRequest) {
	client := g.ClientUpfgw
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// updateData := &SmContextUpdateDataRequest{
	// 	Pei:         "Pei",
	// 	ServingNfId: "ServingNfId",
	// 	N2SmInfo:    &N2SmInformation{},
	// 	Guami: &Guami{
	// 		AmfId: "amfid1",
	// 	},
	// }

	// grpc protos generated function
	res, err := client.SendSmContextUpdateData(ctx, updateData)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	klog.Infof("Response from server: %d", res.GetPduSessionId())
}

// Send Release data
func (g *GrpcClient) SendSmContextReleaseData(releaseData *protos.SmContextReleaseDataRequest) {
	client := g.ClientUpfgw
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// releaseData := &SmContextUpdateDataRequest{
	// 	Pei:         "Pei",
	// 	ServingNfId: "ServingNfId",
	// 	Guami: &Guami{
	// 		AmfId: "amfid1",
	// 	},
	// }

	// grpc protos generated function
	res, err := client.SendSmContextReleaseData(ctx, releaseData)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
	}
	klog.Infof("Response from server: %d", res.GetPduSessionId())
}

func (g *GrpcClient) SendN2Info(n2Info *protos.N2Information) (encodedN2Info *protos.EncodedN2Information) {
	client := g.ClientNgap
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	res, err := client.SendN2Info(ctx, n2Info)
	if err != nil {
		klog.Fatalf("Error. Could not get response: %v", err)
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}
	}
	return res
}
