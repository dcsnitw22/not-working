package grpcSMserver

import (
	"fmt"
	"log"

	// "nas"
	// "nasMain/grpcSmfNas/pb"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/grpcSmfNas/pb"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/nas"
)

type SmfNasGrpcServer struct {
	pb.UnimplementedSmfNasServer
}

func StartSmfNasGrpc() {
	fmt.Println("Started SMF NAS GRPC SERVER")
	klog.Infof("Started SMF-NAS gRPC server")
	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	pb.RegisterSmfNasServer(server, &SmfNasGrpcServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *SmfNasGrpcServer) SendSMData(ctx context.Context, req *pb.SMDataRequest) (*pb.SMDataResponse, error) {
	var res *pb.SMDataResponse
	if req.TypeReq == "" {
		nm, err := decodeRequest(req)
		if err != nil {
			return nil, err
		}

		nasMsg, err := anypb.New(nm)
		if err != nil {
			return nil, err
		}

		res = &pb.SMDataResponse{NasResponse: nasMsg}
	} else {
		nm, err := encodeRequest(req)
		if err != nil {
			return nil, err
		}

		nasMsg, err := anypb.New(nm)
		if err != nil {
			return nil, err
		}

		res = &pb.SMDataResponse{NasResponse: nasMsg}

	}
	return res, nil
}

func decodeRequest(req *pb.SMDataRequest) (proto.Message, error) {
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
		nasMsgs := nasMsg.(nas.PduSessionEstablishmentRequest)
		res = &pb.PDUSEstReqModel{Epd: nasMsgs.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsgs.PDUsessionId)), Pti: proto.Int32(int32(nasMsgs.PTI)), MsgType: nasMsgs.MessageType, DatarateUL: nasMsgs.MaxIntegrityProtectedDataRateUl, DatarateDL: nasMsgs.MaxIntegrityProtectedDataRateDl}
	case "PDU_SESSION_MODIFICATION_REQUEST":
		nasMsgs := nasMsg.(nas.PduSessionModificationRequest)
		res = &pb.PDUSModReqModel{Epd: nasMsgs.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsgs.PDUsessionId)), Pti: proto.Int32(int32(nasMsgs.PTI)), MsgType: nasMsgs.MessageType}
	case "PDU_SESSION_RELEASE_REQUEST":
		nasMsgs := nasMsg.(nas.PduSessionReleaseRequest)
		res = &pb.PDUSRelReqModel{Epd: nasMsgs.ExtendedProtocolDiscriminator, PdusessionID: proto.Int32(int32(nasMsgs.PDUsessionId)), Pti: proto.Int32(int32(nasMsgs.PTI)), MsgType: nasMsgs.MessageType}

	}

	return res, nil

}

func encodeRequest(req *pb.SMDataRequest) (proto.Message, error) {
	var res proto.Message

	switch req.TypeReq {
	case "PDU_SESSION_ESTABLISHMENT_ACCEPT":
		var msg pb.PDUSEstAccModel
		err := req.NasMessage.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		qosRules := addQoSRules(msg.Qosrules)
		// fmt.Println(msg.Sessionambr)
		sessionAMBR := nas.SessionAMBR{IEI: int(msg.Sessionambr.Iei), UnitUL: msg.Sessionambr.UnitUL, RateUL: int(msg.Sessionambr.RateUL), UnitDL: msg.Sessionambr.UnitDL, RateDL: int(msg.Sessionambr.RateDL)}
		// fmt.Println("After session AMBR")
		nasMsg := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, PduSessionType: msg.PdusType, SSCmode: msg.SscMode, QosRuleIEI: int(msg.QosIEI), AuthorizedQoSRules: qosRules, SessionAmbr: sessionAMBR}
		// fmt.Println(nasMsg)
		byteArray, err := nas.ReRouteEncode(nasMsg.ExtendedProtocolDiscriminator, nasMsg.MessageType, nasMsg)
		if err != nil {

			fmt.Println("Hello")
			return nil, err

		}

		res = &pb.ByteDataWrapper{ByteArray: byteArray}
	case "PDU_SESSION_ESTABLISHMENT_REJECT":
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

	case "PDU_SESSION_MODIFICATION_REJECT":
		var msg pb.PDUSModRejModel
		err := req.NasMessage.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		nasMsg := nas.PduSessionModificationReject{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, SMcause: msg.SmCause}
		byteArray, err := nas.ReRouteEncode(nasMsg.ExtendedProtocolDiscriminator, nasMsg.MessageType, nasMsg)

		if err != nil {
			return nil, err

		}
		res = &pb.ByteDataWrapper{ByteArray: byteArray}

	case "PDU_SESSION_RELEASE_REJECT":
		var msg pb.PDUSRelRejModel
		err := req.NasMessage.UnmarshalTo(&msg)
		if err != nil {
			return nil, err
		}

		nasMsg := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: msg.Epd, PDUsessionId: int(msg.PdusessionID), PTI: int(msg.Pti), MessageType: msg.MsgType, SMcause: msg.SmCause}
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
	// fmt.Println(qos)
	return qos

}

func addPFFilter(pfList []*pb.PacketFilters) []nas.PacketFilter {
	pf := make([]nas.PacketFilter, len(pfList))
	for i := range pfList {
		pf[i] = nas.PacketFilter{Identifier: uint8(pfList[i].Identifier), Direction: pfList[i].Direction, Components: pfList[i].Components}
	}
	// fmt.Println()
	// fmt.Println(pf)

	return pf
}
