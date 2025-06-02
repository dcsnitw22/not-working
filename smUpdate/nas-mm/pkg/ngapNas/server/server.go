package grpcserver

import (
	"context"
	"nasMM/pkg/config"
	grpccspnas "nasMM/pkg/grpcCspNas"
	"nasMM/pkg/grpcCspNas/protos/create_sm_context_grpc"
	"nasMM/pkg/grpcNgap"
	"nasMM/pkg/grpcNgap/protos/n1n2messagetransfer"
	"nasMM/pkg/nas"
	"nasMM/pkg/ngapNas/pb"
	"net"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog"
	//"w5gc.io/wipro5gcore/pkg/amf/ngap/nas/grpc/protos/ngapNas/pb"
	//"w5gc.io/wipro5gcore/pkg/smf/nas/config"
)

type GrpcMessageInfo interface{}

type GrpcMessage struct {
	MsgType int32
	GrpcMsg *GrpcMessageInfo
}

type GrpcServer struct {
	pb.UnimplementedDataServiceServer
	ServerIP   string
	ServerPort string
}

// Initialize server with config data
// Receive data and send to grpcChannel
func NewGrpcServer(cfg config.GrpcNgapServerConfig) *GrpcServer {
	return &GrpcServer{
		ServerIP:   cfg.NgapServerIP,
		ServerPort: cfg.NgapServerPort,
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
	pb.RegisterDataServiceServer(server, g)
	klog.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Server handler function
func (g *GrpcServer) SendData(ctx context.Context, in *pb.DataRequest) (*pb.DataResponse, error) {
	klog.Info("Inside Send Data function")
	response := pb.DataResponse{}
	if in.ReqType == "Uplink NAS" {
		klog.Info("Received Uplink NAS request")
		uplinkRequest, err := handleUplinkRequest(in)
		//TO DO: save the information in Redis of amf
		klog.Info(uplinkRequest)
		response.Type = true
		response.NasPdu = nil
		if err != nil {
			response.Error = err.Error()

		}
		response.Error = ""

	} else if in.ReqType == "Registration Request" {
		klog.Info("Received Registration Request")
		registrationRequest, registrationAccept, err := handleRegistrationRequest(in)
		klog.Info(registrationAccept)

		response.Type = false
		response.NasPdu = registrationAccept

		if err != nil {
			response.Error = err.Error()

		}
		response.Error = ""
		supi := registrationRequest.MobileIdentity.MCC + registrationRequest.MobileIdentity.MNC + registrationRequest.MobileIdentity.SchemeOutput
		response.Supi = supi
		//TODO : Save the registration request data in redis of amf
		klog.Info(registrationRequest)

	}
	klog.Info(&response)
	return &response, nil
}

func handleRegistrationRequest(in *pb.DataRequest) (nas.RegistrationRequestModel, []byte, error) {
	var decodedNasMsg nas.RegistrationRequestModel
	var inputData pb.NasMessage
	in.Data.UnmarshalTo(&inputData)

	// klog.Info(inputData.NasPdu)

	epd, msgType, err := nas.Classify(inputData.NasPdu)
	if err != nil {
		return decodedNasMsg, nil, err
	}

	decodedMsg, err := nas.ReRouteDecode(epd, msgType, inputData.NasPdu)
	if err != nil {
		return decodedNasMsg, nil, err
	}

	decodedNasMsg = decodedMsg.(nas.RegistrationRequestModel)

	//Construct Registration Accept with dummy data
	registrationAccept := nas.RegistrationAcceptModel{ExtendedProtocolDiscriminator: decodedNasMsg.ExtendedProtocolDiscriminator, SecurityHeaderType: "Plain 5GS NAS message", MessageType: "REGISTRATION_ACCEPT", RegResult: "3GPP access", Sms: "SMS over NAS not allowed", NssaPerformed: "Network slice-specific authentication and authorization is not to be performed", EmergencyReg: "Not registered for emergency services", RoamingReg: "No additional information"}

	registrationAcceptByteArray, err := nas.ReRouteEncode(registrationAccept.ExtendedProtocolDiscriminator, registrationAccept.MessageType, registrationAccept)
	if err != nil {
		return decodedNasMsg, nil, err
	}

	return decodedNasMsg, registrationAcceptByteArray, nil

}

func handleUplinkRequest(in *pb.DataRequest) (nas.ULNasModel, error) {
	var inputData pb.CreateSmContext
	var uplinkRequest nas.ULNasModel
	in.Data.UnmarshalTo(&inputData)

	klog.Info(inputData.N1SmContainer)

	epd, msgType, err := nas.Classify(inputData.N1SmContainer)
	if err != nil {
		return uplinkRequest, nil
	}

	decodedMsg, err := nas.ReRouteDecode(epd, msgType, inputData.N1SmContainer)
	if err != nil {
		return uplinkRequest, err
	}
	uplinkRequest = decodedMsg.(nas.ULNasModel)

	klog.Info("Finished decoding uplink mm message")

	//Create CSP grpc client & send info
	// pduSessionIDIEI := uplinkRequest.PduSessionIdIEI
	// pduSessionID := uplinkRequest.PduSessionId

	klog.Info("Sending data to CSP")
	payloadContainerToCsp := inputData.N1SmContainer[6:]

	cspNssai := create_sm_context_grpc.Snssai{Sst: inputData.Snssai.Sst, Sd: inputData.Snssai.Sd}
	cspPlmnID := create_sm_context_grpc.PlmnId{Mcc: inputData.NrLocation.Tai.PlmnId.Mcc, Mnc: inputData.NrLocation.Tai.PlmnId.Mnc}
	cspTai := create_sm_context_grpc.Tai{PlmnId: &cspPlmnID, Tac: inputData.NrLocation.Tai.Tac}
	cspNcgi := create_sm_context_grpc.Ncgi{PlmnId: cspTai.PlmnId, NrCellId: inputData.NrLocation.Ncgi.NrCellId}

	cspNrLoc := create_sm_context_grpc.NrLocation{Tai: &cspTai, Ncgi: &cspNcgi, AgeOfLocationInformation: inputData.NrLocation.AgeOfLocationInformation, UeLocationTimestamp: inputData.NrLocation.UeLocationTimestamp, GeographicalInformation: inputData.NrLocation.GeographicalInformation, GeodeticInformation: inputData.NrLocation.GeodeticInformation, GlobalGnbId: inputData.NrLocation.GlobalGnbId}
	cspMsg := create_sm_context_grpc.CreateSmContextDataFromNasMod{AnType: inputData.AnType, Snssai: &cspNssai, NrLocation: &cspNrLoc, PduSessionId: int32(uplinkRequest.PduSessionId), N1SmContainer: payloadContainerToCsp, Supi: inputData.Supi}

	klog.Info(cspMsg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//client, ctx := grpccspnas.CreateNASCSPClient()
	client := grpccspnas.CreateNASCSPClient()
	resp, err := client.SendDataForCreateSmContext(ctx, &cspMsg)

	klog.Info("Sent data to CSP")

	klog.Info(err)

	if err != nil {
		return uplinkRequest, err
	}

	klog.Info(resp)

	return uplinkRequest, nil

}

func (g *GrpcServer) HandleDownlink(ctx context.Context, in *pb.DLRequest) (*pb.DLResponse, error) {
	klog.Info("Inside Handle Downlink Request")

	var dlResp pb.DLResponse

	downlinkMsg := nas.DLNasModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "DL_NAS_TRANSPORT", PayLoadContainerType: "N1 SM information", PayLoadContainer: in.N1DataBytes, PduSessionIdIEI: 18, PduSessionId: int(in.PduSessionId)}
	encodedDLmsg, err := nas.EncodeDLNAS(downlinkMsg)
	if err != nil {
		dlResp.Error = err.Error()
		return &dlResp, nil
	}

	klog.Info(encodedDLmsg)

	n1n2Data := n1n2messagetransfer.N1N2Data{N1DataBytes: encodedDLmsg, N2DataBytes: in.N2DataBytes, UeContextId: in.UeContextId}

	ctxNgap, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := grpcNgap.CreateNgapClient()
	resp, err := client.SendN1N2DataTransfer(ctxNgap, &n1n2Data)

	klog.Info("Sent data to Ngap")

	if err != nil {
		klog.Fatalf("Error sending data to Ngap: %v", err)
	}

	klog.Info(resp)
	if resp.Err != "" {
		dlResp.Error = resp.Err
	}

	dlResp.Error = ""

	return &dlResp, nil
}
