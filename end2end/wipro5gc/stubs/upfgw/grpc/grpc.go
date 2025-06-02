package grpc

import (
	"time"

	//"net"

	//"w5gc.io/wipro5gcore/openapi"

	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/grpcclient"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/grpcserver"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/protos"
)

type Grpc interface {
	Start()
	WatchGrpcChannel() chan *grpcserver.GrpcMessage
	SendN1N2MessageData(string, *protos.N1N2MessageTransferDataRequest)
}

type GrpcInfo struct {
	grpcStartTime time.Time
	grpcServer    grpcserver.GrpcServer
	grpcClient    grpcclient.GrpcClient
}

// Initialize Grpc with Server and Client
func NewGrpc() Grpc {
	return &GrpcInfo{
		grpcServer: *grpcserver.NewGrpcServer(),
		grpcClient: *grpcclient.NewGrpcClient(),
	}
}

func (g *GrpcInfo) Start() {
	go g.grpcServer.Start()
	g.grpcClient.Start()

	// Stub to listen for incoming channel data
	go func() {
		for {
			select {
			case ch := <-g.WatchGrpcChannel():
				smContextID := (*ch.GrpcMsg).(string)
				// Stub dummy data
				n1n2Data := protos.N1N2MessageTransferDataRequest{
					SmContextID:       smContextID,
					NgapleType:        "PDU_RES_SETUP_REQ",
					TunnelEndpointAdr: "1.2.3.4",
					Gtpteid:           "123",

					// N1Message: &protos.N1Message{
					// 	N1MsgClass:                    "5GMM",
					// 	ExtendedProtocolDiscriminator: "String",
					// },
					// N2Info: &protos.N2Info{
					// 	N2InformationClass:         "SM",
					// 	PduSessionId:               pduSessionId,
					// 	PduSessionType:             "T1",
					// 	NgapIeType:                 "PDU_RES_SETUP_REQ",
					// 	GtpTunnelEndpointIpAddress: "1.2.3.4",
					// 	GtpTeid:                    "123",
					// 	Qfi:                        1,
					// 	QosFlowLevelQosParameters: &protos.QosFlowLevelQosParameters{
					// 		Fqi: 1,
					// 		Arp: &protos.Arp{
					// 			PriorityLevel:           "100",
					// 			PreemptionCapability:    "NOT_PREEMPT",
					// 			PreemptionVulnerability: "NOT_PREEMPTABLE",
					// 		},
					// 		GbrQosFlowInformation: &protos.GbrQoSFlowInformation{
					// 			MaximumFlowBitrateDownlink:    100,
					// 			MaximumFlowBitrateUplink:      100,
					// 				GuaranteedFlowBitrateDownlink: 100,
					// 				GuaranteedFlowBitrateUplink:   100,
					// 			},
					// 		},
					// 	},
					// 	OldGuami: &protos.Guami{
					// 		PlmnId: &protos.PlmnId{
					// 			Mcc: "404",
					// 			Mnc: "10",
					// 		},
					// 		AmfId: "218A9E",
					// 	},
					// 	Arp: &protos.Arp{
					// 		PriorityLevel:           "100",
					// 		PreemptionCapability:    "NOT_PREEMPT",
					// 		PreemptionVulnerability: "NOT_PREEMPTABLE",
					// 	},
				}
				// Send to pdusmsp server
				g.SendN1N2MessageData(smContextID, &n1n2Data)
			}
		}
	}()
}

func (g *GrpcInfo) WatchGrpcChannel() chan *grpcserver.GrpcMessage {
	return g.grpcServer.WatchGrpcChannel()
}

func (g *GrpcInfo) SendN1N2MessageData(smContextID string, n1n2Data *protos.N1N2MessageTransferDataRequest) {
	g.grpcClient.SendN1N2MessageData(smContextID, n1n2Data)
}
