package grpcserver

import (
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos"
)

/*const (
	GrpcChannelCapacity = 100
)*/

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	protos.UnimplementedN2InfoNgapEncoderServer
	//grpcChannel chan *GrpcMessage
	ServerIP   string
	ServerPort string
}

// Initialize server with config data
// Receive data and send to grpcChannel
func NewGrpcServer(cfg config.GrpcServerInfoConfig) *GrpcServer {
	return &GrpcServer{
		//grpcChannel: make(chan *GrpcMessage, GrpcChannelCapacity),
		ServerIP:   cfg.ServerIP,
		ServerPort: cfg.ServerPort,
	}
}

// Start server
func (g *GrpcServer) Start() {
	servAddr := g.ServerIP + ":" + g.ServerPort
	lis, err := net.Listen("tcp", servAddr)
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	protos.RegisterN2InfoNgapEncoderServer(server, g)
	klog.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Watch channel data
/*func (g *GrpcServer) WatchGrpcChannel() chan *GrpcMessage {
	return g.grpcChannel
}*/

func (g *GrpcServer) SendN2Info(ctx context.Context, in *protos.N2Information) (*protos.EncodedN2Information, error) {
	ip := []byte(in.GtpTunnelEndpointIpAddress)
	ulNguUpTnlInformation := asn1gen.UPTransportLayerInformation{
		T: 1,
		U: struct {
			GTPTunnel        *asn1gen.GTPTunnel
			ChoiceExtensions *asn1gen.ProtocolIESingleContainer
		}{
			GTPTunnel: &asn1gen.GTPTunnel{
				TransportLayerAddress: asn1gen.TransportLayerAddress{
					Bytes:     ip,
					BitLength: len(ip) * 8,
				},
				GTPTEID: asn1gen.GTPTEID(in.GtpTeid),
			},
			ChoiceExtensions: nil,
		},
	}
	encodedUlNguUpTnlInformation, err := asn1gen.Marshal(ulNguUpTnlInformation)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	//PDU Session Type
	pduSessionType := asn1gen.PDUSessionType(in.PduSessionType)
	encodedPduSessionType, err := asn1gen.Marshal(pduSessionType)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	pLevel, err := strconv.Atoi(in.QosFlowLevelQosParameters.Arp.PriorityLevel)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	pCap, err := strconv.Atoi(in.QosFlowLevelQosParameters.Arp.PreemptionCapability)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	pVul, err := strconv.Atoi(in.QosFlowLevelQosParameters.Arp.PreemptionVulnerability)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	//QoS Flow Setup Request List
	qosFlowSetupRequestList := asn1gen.QosFlowSetupRequestList{
		asn1gen.QosFlowSetupRequestItem{
			QosFlowIdentifier: asn1gen.QosFlowIdentifier(in.Qfi),
			QosFlowLevelQosParameters: asn1gen.QosFlowLevelQosParameters{
				QosCharacteristics: asn1gen.QosCharacteristics{
					T: 1,
					U: struct {
						NonDynamic5QI    *asn1gen.NonDynamic5QIDescriptor
						Dynamic5QI       *asn1gen.Dynamic5QIDescriptor
						ChoiceExtensions *asn1gen.ProtocolIESingleContainer
					}{
						NonDynamic5QI: &asn1gen.NonDynamic5QIDescriptor{
							FiveQI: asn1gen.FiveQI(in.QosFlowLevelQosParameters.Fqi),
						},
						Dynamic5QI:       nil,
						ChoiceExtensions: nil,
					},
				},
				AllocationAndRetentionPriority: asn1gen.AllocationAndRetentionPriority{
					PriorityLevelARP:        asn1gen.PriorityLevelARP(pLevel),
					PreEmptionCapability:    asn1gen.PreEmptionCapability(pCap),
					PreEmptionVulnerability: asn1gen.PreEmptionVulnerability(pVul),
				},
				GBRQosInformation: &asn1gen.GBRQosInformation{},
			},
		},
	}
	encodedQosFlowSetupRequestList, err := asn1gen.Marshal(qosFlowSetupRequestList)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}

	pduSessionResourceSetupRequestTransfer := asn1gen.PDUSessionResourceSetupRequestTransfer{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdULNGUUPTNLInformation,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedUlNguUpTnlInformation,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdPDUSessionType,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedPduSessionType,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdQosFlowSetupRequestList,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedQosFlowSetupRequestList,
			},
		},
	}
	encodedPduSessionResourceSetupRequestTransfer, err := asn1gen.Marshal(pduSessionResourceSetupRequestTransfer)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}
	return &protos.EncodedN2Information{
		EncodedData: encodedPduSessionResourceSetupRequestTransfer,
		Error:       "nil",
	}, nil
}
