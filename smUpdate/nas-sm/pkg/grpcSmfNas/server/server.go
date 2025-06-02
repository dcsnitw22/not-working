package grpcSMserver

import (
	"fmt"
	"nasMain/pkg/grpcSmfNas/pb"
	"nasMain/pkg/nas"

	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog"
)

type SmfNasGrpcServer struct {
	pb.UnimplementedSmfNasServer
}

func StartSmfNasGrpc() {
	fmt.Println("Started SMF NAS GRPC SERVER")
	klog.Infof("Started SMF-NAS gRPC server")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	pb.RegisterSmfNasServer(server, &SmfNasGrpcServer{})
	klog.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

// Function to Handle Release Requests
func (s *SmfNasGrpcServer) HandleUpdateRelease(ctx context.Context, req *pb.UpRelRequest) (*pb.UpRelRespone, error) {
	klog.Info("Handling Update and Release NAS Requests")
	res := &pb.UpRelRespone{}

	if req.ReqType == "Decode" {
		var byteWrapper pb.ByteDataWrapper
		//Unmarshal byte array
		err := proto.Unmarshal(req.NasRelMsg.Value, &byteWrapper)
		if err != nil {
			return nil, err
		}
		rel := make([]byte, len(byteWrapper.ByteArray))
		copy(rel, byteWrapper.ByteArray)

		switch rel[3] {
		case nas.PduSReleaseRequest:
			req.ReqType = "PDU_SESSION_RELEASE_REQUEST"
		case nas.PduSReleaseComplete:
			req.ReqType = "PDU_SESSION_RELEASE_COMPLETE"
		case nas.PduSModificationRequest:
			req.ReqType = "PDU_SESSION_MODIFICATION_REQUEST"
		case nas.PduSModificationComplete:
			req.ReqType = "PDU_SESSION_MODIFICATION_COMPLETE"
		}

	}

	switch req.ReqType {
	case "PDU_SESSION_RELEASE_REQUEST":
		klog.Info("Received Release Request")
		var byteWrapper pb.ByteDataWrapper
		//Unmarshal byte array
		err := proto.Unmarshal(req.NasRelMsg.Value, &byteWrapper)
		if err != nil {
			return nil, err
		}
		reqRel := make([]byte, len(byteWrapper.ByteArray))
		copy(reqRel, byteWrapper.ByteArray)

		//Decode the byte array
		nasMsg, err := nas.DecodePduSessionReleaseRequest(reqRel)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil
		}
		//Convert result to proto data model
		pbNasMsg := pb.PDUSRelReqModel{Epd: nasMsg.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsg.PDUsessionId)), Pti: proto.Int32(int32(nasMsg.PTI)), MsgType: nasMsg.MessageType, SmCause: nasMsg.SMCause}
		//Convert decoded message to "any" type
		resNasMsg, err := anypb.New(&pbNasMsg)
		if err != nil {
			return nil, err
		}
		//Send back response
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil

	case "PDU_SESSION_RELEASE_COMMAND":
		klog.Info("Received Release Command")
		var model pb.PDUSRelCommandModel
		//Unmarshal to Release Command Model
		err := req.NasRelMsg.UnmarshalTo(&model)
		if err != nil {
			return nil, err
		}
		//Convert to nas model
		commandData := nas.ReleaseCommandModel{ExtendedProtocolDiscriminator: model.Epd, PDUsessionId: int(*model.PdusessionID), PTI: int(*model.Pti), MessageType: model.MsgType, SMCause: model.SmCause}
		//Encode the input data
		nasRes, err := nas.EncodeReleaseCommand(commandData)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil
		}
		//Convert byte array to proto type
		resNasRes := &pb.ByteDataWrapper{ByteArray: nasRes}
		resNasMsg, err := anypb.New(resNasRes)
		if err != nil {
			return nil, err
		}
		//Create response and send back
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil

	case "PDU_SESSION_RELEASE_COMPLETE":
		klog.Info("Received Release Complete")
		var byteWrapper pb.ByteDataWrapper
		//Unmarshal byte array
		err := proto.Unmarshal(req.NasRelMsg.Value, &byteWrapper)
		if err != nil {
			return nil, err
		}
		comRel := make([]byte, len(byteWrapper.ByteArray))
		copy(comRel, byteWrapper.ByteArray)

		//Decode the byte array
		nasMsg, err := nas.DecodeReleaseComplete(comRel)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil
		}
		//Convert result to proto data model
		pbNasMsg := pb.PDUSRelCompleteModel{Epd: nasMsg.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsg.PDUsessionId)), Pti: proto.Int32(int32(nasMsg.PTI)), MsgType: nasMsg.MessageType, SmCause: nasMsg.SMCause}
		//Convert decoded message to any type
		resNasMsg, err := anypb.New(&pbNasMsg)
		if err != nil {
			return nil, err
		}
		//Send back response
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil

	case "PDU_SESSION_RELEASE_REJECT":
		klog.Info("Received Release Reject")
		var model pb.PDUSRelRejModel
		//Unmarshal to Release Reject Model
		err := req.NasRelMsg.UnmarshalTo(&model)
		if err != nil {
			return nil, err
		}
		//Convert to nas model
		rejectData := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: model.Epd, PDUsessionId: int(*model.PdusessionID), PTI: int(*model.Pti), MessageType: model.MsgType, SMcause: model.SmCause}
		//Encode the input data
		nasRes, err := nas.EncodePduSessionReleaseReject(rejectData)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil

		}
		//Convert byte array to proto type
		resNasRes := &pb.ByteDataWrapper{ByteArray: nasRes}
		//Convert to "any" type
		resNasMsg, err := anypb.New(resNasRes)
		if err != nil {
			return nil, err
		}
		//Send back the response
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil
	case "PDU_SESSION_MODIFICATION_REQUEST":
		klog.Info("Received Modification Request")
		var byteWrapper pb.ByteDataWrapper
		//Unmarshal byte array
		err := proto.Unmarshal(req.NasRelMsg.Value, &byteWrapper)
		if err != nil {
			return nil, err
		}
		reqMod := make([]byte, len(byteWrapper.ByteArray))
		copy(reqMod, byteWrapper.ByteArray)

		//Decode the byte array
		nasMsg, err := nas.DecodePduSessionModificationRequest(reqMod)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil

		}
		//Convert result to proto data model
		pbNasMsg := pb.PDUSModReqModel{Epd: nasMsg.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsg.PDUsessionId)), Pti: proto.Int32(int32(nasMsg.PTI)), MsgType: nasMsg.MessageType}
		//Convert decoded message to any type
		resNasMsg, err := anypb.New(&pbNasMsg)
		if err != nil {
			return nil, err
		}
		//Send back response
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil
	case "PDU_SESSION_MODIFICATION_REJECT":
		klog.Info("Received Modification Reject")
		var msg pb.PDUSModRejModel
		err := req.NasRelMsg.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		nasMsg := nas.PduSessionModificationReject{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, SMcause: msg.SmCause}
		byteArray, err := nas.EncodePduSessionModificationReject(nasMsg)
		if err != nil {
			res.Error = err.Error()
			res.NasResponse = nil
			res.ReqType = req.ReqType
			return res, nil

		}
		nasRes := &pb.ByteDataWrapper{ByteArray: byteArray}
		//Convert to "any" type
		resNasMsg, err := anypb.New(nasRes)
		if err != nil {
			return nil, err
		}
		//Send back the response
		res.NasResponse = resNasMsg
		res.Error = ""
		res.ReqType = req.ReqType
		klog.Info(res)
		return res, nil

	}
	//If request type doesnt match send back the error
	res.NasResponse = nil
	res.Error = "Unknown Request Type"
	res.ReqType = req.ReqType
	return res, nil
}

func (s *SmfNasGrpcServer) HandleEstablishment(ctx context.Context, req *pb.EstRequest) (*pb.EstResponse, error) {
	var res *pb.EstResponse
	if req.TypeReq == "" {
		nm, err := decodeRequest(req)
		if err != nil {
			return nil, err
		}

		nasMsg, err := anypb.New(nm)
		if err != nil {
			return nil, err
		}

		res = &pb.EstResponse{NasResponse: nasMsg}
	} else {
		nm, err := encodeRequest(req)
		if err != nil {
			return nil, err
		}

		nasMsg, err := anypb.New(nm)
		if err != nil {
			return nil, err
		}

		res = &pb.EstResponse{NasResponse: nasMsg}

	}
	return res, nil
}

func decodeRequest(req *pb.EstRequest) (proto.Message, error) {
	klog.Info("Received Decode Request")
	var res proto.Message
	// fmt.Println("Inside Decode function")
	var byteWrapper pb.ByteDataWrapper
	err := proto.Unmarshal(req.NasMessage.Value, &byteWrapper)
	if err != nil {
		return nil, err
	}
	byteArray := make([]byte, len(byteWrapper.ByteArray))
	copy(byteArray, byteWrapper.ByteArray)
	epd, msgType, err := nas.Classify(byteArray)
	if err != nil {
		return nil, err

	}
	nasMsg, err := nas.ReRouteDecode(epd, msgType, byteArray)
	if err != nil {
		return nil, err

	}

	switch msgType {
	case "PDU_SESSION_ESTABLISHMENT_REQUEST":
		klog.Info("Received Pdu Session Establishment Request")
		nasMsgs := nasMsg.(nas.PduSessionEstablishmentRequest)
		res = &pb.PDUSEstReqModel{Epd: nasMsgs.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsgs.PDUsessionId)), Pti: proto.Int32(int32(nasMsgs.PTI)), MsgType: nasMsgs.MessageType, DatarateUL: nasMsgs.MaxIntegrityProtectedDataRateUl, DatarateDL: nasMsgs.MaxIntegrityProtectedDataRateDl}
	}

	return res, nil

}

func encodeRequest(req *pb.EstRequest) (proto.Message, error) {
	klog.Info("Receieved Encode Request")
	var res proto.Message

	switch req.TypeReq {
	case "PDU_SESSION_ESTABLISHMENT_ACCEPT":
		klog.Info("Receieved Pdu Session Establishment Accept")
		var msg pb.PDUSEstAccModel
		err := req.NasMessage.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		qosRules := addQoSRules(msg.Qosrules)
		// fmt.Println(msg.Sessionambr)
		sessionAMBR := nas.SessionAMBR{IEI: int(msg.Sessionambr.Iei), UnitUL: msg.Sessionambr.UnitUL, RateUL: int(msg.Sessionambr.RateUL), UnitDL: msg.Sessionambr.UnitDL, RateDL: int(msg.Sessionambr.RateDL)}
		// fmt.Println("After session AMBR")
		nasMsg := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, PduSessionType: msg.PdusType, SSCmode: msg.SscMode, QosRuleIEI: int(msg.QosIEI), AuthorizedQoSRules: qosRules, SessionAmbr: sessionAMBR, Si6lla: msg.Si6Lla, PduSessionTypeVal: msg.PduSessionTypeVal, PduAddrInfo: msg.PduAddrInfo}
		klog.Info(nasMsg)
		// fmt.Println(nasMsg)
		byteArray, err := nas.ReRouteEncode(nasMsg.ExtendedProtocolDiscriminator, nasMsg.MessageType, nasMsg)
		if err != nil {
			return nil, err
		}

		res = &pb.ByteDataWrapper{ByteArray: byteArray}
	case "PDU_SESSION_ESTABLISHMENT_REJECT":
		klog.Info("Received Pdu Session EStablishment Reject")
		var msg pb.PDUSEstRejModel
		err := req.NasMessage.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		nasMsg := nas.PduSessionEstablishmentReject{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, SMcause: msg.SmCause}
		byteArray, err := nas.ReRouteEncode(nasMsg.ExtendedProtocolDiscriminator, nasMsg.MessageType, nasMsg)

		if err != nil {
			return nil, err

		}
		res = &pb.ByteDataWrapper{ByteArray: byteArray}
	}

	return res, nil

}

func addQoSRules(pbQoSRules []*pb.QosRules) []nas.QoSRule {
	qos := make([]nas.QoSRule, len(pbQoSRules))
	for i := range pbQoSRules {
		pf := addPFFilter(pbQoSRules[i].Pf)
		qos[i] = nas.QoSRule{QoSIdentifier: pbQoSRules[i].Qosidentifier, Operation: pbQoSRules[i].Operation, DQR: pbQoSRules[i].Dqr, PacketFilterList: pf, Precedence: uint8(pbQoSRules[i].Precidence), Segregation: pbQoSRules[i].Seg, QFI: pbQoSRules[i].Qfi}

	}

	return qos

}

func addPFFilter(pfList []*pb.PacketFilters) []nas.PacketFilter {
	pf := make([]nas.PacketFilter, len(pfList))
	for i := range pfList {
		pf[i] = nas.PacketFilter{Identifier: uint8(pfList[i].Identifier), Direction: pfList[i].Direction, Components: pfList[i].Components}
	}

	return pf
}
