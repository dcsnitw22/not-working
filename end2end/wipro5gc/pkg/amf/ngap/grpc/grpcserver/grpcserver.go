package grpcserver

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"strconv"

	"google.golang.org/grpc"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos"
)

/*const (
	GrpcChannelCapacity = 100
)*/

var N1Message = make(chan []byte, 1)
var N2Info = make(chan []byte, 1)
var UeContextId = make(chan []byte, 1)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	protos.UnimplementedN2InfoNgapEncoderServer
	protos.UnimplementedN1N2DataTransferServer
	grpcChannel       chan *GrpcMessage
	ServerIPPdusmsp   string
	ServerPortPdusmsp string
	ServerIPNas       string
	ServerPortNas     string
}

// Initialize server with config data
// Receive data and send to grpcChannel
func NewGrpcServer(cfg config.GrpcServerInfoConfig) *GrpcServer {
	klog.Infof("ngap grpc server config +%v", cfg)
	return &GrpcServer{
		grpcChannel:       make(chan *GrpcMessage), //, GrpcChannelCapacity),
		ServerIPPdusmsp:   cfg.ServerIPPdusmsp,
		ServerPortPdusmsp: cfg.ServerPortPdusmsp,
		ServerIPNas:       cfg.ServerIPNas,
		ServerPortNas:     cfg.ServerPortNas,
	}
}

// Start server for Pdusmsp
func (g *GrpcServer) Start() {
	servAddrPdusmsp := g.ServerIPPdusmsp + ":" + g.ServerPortPdusmsp
	klog.Info("grpc server config : ", servAddrPdusmsp)
	lis, err := net.Listen("tcp", servAddrPdusmsp)
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	serverPdusmsp := grpc.NewServer()
	protos.RegisterN2InfoNgapEncoderServer(serverPdusmsp, g)
	klog.Infof("server listening at %v", lis.Addr())
	if err := serverPdusmsp.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Start server for Nas
func (g *GrpcServer) Start2() {
	servAddrNas := g.ServerIPNas + ":" + g.ServerPortNas
	klog.Info("grpc server config : ", servAddrNas)
	lis2, err := net.Listen("tcp", servAddrNas)
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	serverNas := grpc.NewServer()
	protos.RegisterN1N2DataTransferServer(serverNas, g)
	klog.Infof("server listening at %v", lis2.Addr())
	if err := serverNas.Serve(lis2); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Watch channel data
func (g *GrpcServer) WatchGrpcChannel() chan *GrpcMessage {
	return g.grpcChannel
}

func (g *GrpcServer) SendN2Info(ctx context.Context, in *protos.N2Information) (*protos.EncodedN2Information, error) {
	ip := []byte(in.GtpTunnelEndpointIpAddress)
	klog.Info("gtp tunnel endpoint id : ", in.GtpTeid, " and its type : ", reflect.TypeOf(in.GtpTeid))
	teid := asn1rt.OctetString([]byte(in.GtpTeid))
	klog.Info("gtp tunnel endpoint id bytes : ", teid, " and its type : ", reflect.TypeOf(teid))
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
				GTPTEID: asn1gen.GTPTEID(teid),
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
	var pCap asn1gen.PreEmptionCapability
	if in.QosFlowLevelQosParameters.Arp.PreemptionCapability == "NOT_PREEMPT" {
		pCap = asn1gen.PreEmptionCapability(asn1gen.PreEmptionCapabilityShallNotTriggerPreEmption)
	} else {
		e := "check server side code or check preemption capability value"
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       e,
		}, errors.New(e)
	}

	var pVul asn1gen.PreEmptionVulnerability
	if in.QosFlowLevelQosParameters.Arp.PreemptionVulnerability == "NOT_PREEMPTABLE" {
		pVul = asn1gen.PreEmptionVulnerability(asn1gen.PreEmptionVulnerabilityNotPreEmptable)
	} else {
		e := "check server side code or check preemption vulnerability value"
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       e,
		}, errors.New(e)
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
					PreEmptionCapability:    pCap,
					PreEmptionVulnerability: pVul,
				},
				GBRQosInformation: &asn1gen.GBRQosInformation{
					MaximumFlowBitRateDL:    asn1gen.BitRate(in.QosFlowLevelQosParameters.GbrQosFlowInformation.MaximumFlowBitrateDownlink),
					MaximumFlowBitRateUL:    asn1gen.BitRate(in.QosFlowLevelQosParameters.GbrQosFlowInformation.MaximumFlowBitrateUplink),
					GuaranteedFlowBitRateDL: asn1gen.BitRate(in.QosFlowLevelQosParameters.GbrQosFlowInformation.GuaranteedFlowBitrateDownlink),
					GuaranteedFlowBitRateUL: asn1gen.BitRate(in.QosFlowLevelQosParameters.GbrQosFlowInformation.GuaranteedFlowBitrateUplink),
				},
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
	/*encodedPduSessionResourceSetupRequestTransfer, err := asn1gen.Marshal(pduSessionResourceSetupRequestTransfer)
	if err != nil {
		return &protos.EncodedN2Information{
			EncodedData: nil,
			Error:       err.Error(),
		}, err
	}*/
	pduSessionResourceSetupRequestTransferBytes, err := json.Marshal(pduSessionResourceSetupRequestTransfer)
	if err != nil {
		klog.Error("error in unmarshaling pdu session resource setup request transfer bytes : ", err)
	} else {
		klog.Info("Pdu session resource setup request transfer bytes: ", pduSessionResourceSetupRequestTransferBytes)
	}
	return &protos.EncodedN2Information{
		EncodedData: pduSessionResourceSetupRequestTransferBytes,
		Error:       "nil",
	}, nil
}

func (g *GrpcServer) SendN1N2DataTransfer(ctx context.Context, in *protos.N1N2Data) (*protos.Error, error) {
	klog.Info("Inside sendn1n2data grpc server function")
	go func() {
		klog.Info("Sending n1n2data in grpc channel in ngap")
		x := GrpcMessageInfo(in)
		g.grpcChannel <- &GrpcMessage{GrpcMsg: &x}
	}()
	// N1Message <- in.N1DataBytes
	// N2Info <- in.N2DataBytes
	return &protos.Error{
		Err: "nil",
	}, nil
}
