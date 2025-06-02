package main

import (
	"context"
	"fmt"
	"log"
	"nasMain/pkg/grpcSmfNas/pb"
	"nasMain/pkg/nas"

	// "nas"
	// "nasMain/grpcSmfNas/pb"
	// grpcSMserver "nasMain/grpcSmfNas/server"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
)

func Trail() {
	// go grpcSMserver.StartSmfNasGrpc()

	// <-time.After(time.Second)
	// Server address and port
	serverAddr := "grpcnassmf-service.nassm.svc.cluster.local:50052"

	// Create a new gRPC client
	clientConn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewSmfNasClient(clientConn)

	// // Input data
	// modificationRequestFile, _ := os.OpenFile("/home/ubuntu/NAS/testFiles/MinPDUSessionModificationRequest", os.O_RDONLY, 0)
	// // Get the file size
	// fileInfo, _ := modificationRequestFile.Stat()
	// fileSize := fileInfo.Size()
	// // Create a byte array with the size of the file
	// modificationRequestByteArray := make([]byte, fileSize)
	// // Read the bytes from the file into the byte array
	// _, _ = modificationRequestFile.Read(modificationRequestByteArray)
	// fmt.Println("Byte Array:", modificationRequestByteArray)
	// // byteArray := []byte{0x01, 0x02, 0x03, 0x04} // Example byte array
	// anyMessage, err := anypb.New(&pb.ByteDataWrapper{ByteArray: modificationRequestByteArray})
	// if err != nil {
	// 	log.Fatalf("Failed to create Any message: %v", err)
	// }

	// reqType := "" // Empty string for the request type

	releaseReject := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_RELEASE_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
	// fmt.Printf("Dummy Data for Release Reject: %+v\n", releaseReject)

	// protoMsg := pb.PDUSRelRejModel{Epd: releaseReject.ExtendedProtocolDiscriminator, PdusessionID: int32(releaseReject.PDUsessionId), Pti: int32(releaseReject.PTI), MsgType: releaseReject.MessageType, SmCause: releaseReject.SMcause}

	// anyMessage, err := anypb.New(&protoMsg)
	// if err != nil {
	// 	log.Fatalf("Failed to create Any message: %v", err)
	// }

	// reqType := "PDU_SESSION_RELEASE_REJECT"

	fmt.Println("### Encode DL NAS Message ###")
	dlNAS := nas.DLNasModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "DL_NAS_TRANSPORT", PayLoadContainerType: "N1 SM information", PayLoadContainer: releaseReject, PduSessionIdIEI: 1, PduSessionId: 2}
	fmt.Printf("Dummy Data for DL NAS Transport: %+v\n", dlNAS)
	dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
	fmt.Println("Encoded DL NAS Transport:", dlNasByteArray)
	fmt.Println("Error while encoding DL NAS Transport:", err)

	//UpLink NAS
	fmt.Println("### Decode Uplink NAS message ###")
	ulNasFile, _ := os.OpenFile("/home/wipro/mounika_nas/nas-sm/testFiles/MinULNASTransport", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ := ulNasFile.Stat()
	fileSize := fileInfo.Size()
	// Create a byte array with the size of the file
	ulNasByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = ulNasFile.Read(ulNasByteArray)
	fmt.Println("Byte Array:", ulNasByteArray)
	ul := []byte{126, 0, 103, 1, 0, 21, 46, 1, 1, 193, 255, 255, 145, 161, 40, 1, 0, 123, 0, 7, 128, 0, 10, 0, 0, 13, 0, 18, 1, 129, 34, 4, 1, 0, 0, 1, 37, 9, 8, 105, 110, 116, 101, 114, 110, 101, 116}
	ulNasPDU, err := nas.DecodeULNas(ul)
	fmt.Printf("Decoded UL Nas Message:  %+v\n", ulNasPDU)
	fmt.Println("Error while decoding UL Nas Message:", err)

	fmt.Println("### Establishment Accept ###")
	comp := []string{"Match-all type"}
	pf := nas.PacketFilter{Identifier: 15, Direction: "BIDIRECTIONAL", Components: comp}
	pflist := []nas.PacketFilter{pf}
	qos := nas.QoSRule{QoSIdentifier: "QRI 1", Operation: "Create new QoS rule", DQR: "DEFAULT_QoS_RULE", Precedence: 255, Segregation: "Segregation not requested", QFI: "QFI 1", PacketFilterList: pflist}
	ses := nas.SessionAMBR{IEI: 1, UnitUL: "MULT_1Kbps", RateUL: 30, UnitDL: "MULT_1Kbps", RateDL: 40}
	qosList := []nas.QoSRule{qos, qos}
	establishmentAccept := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_ACCEPT", PduSessionType: "IPV4", SSCmode: "SSC_MODE_1", AuthorizedQoSRules: qosList, SessionAmbr: ses}
	// pbPF1 := pb.PF{Identifier: uint32(pf.Identifier), Direction: pf.Direction, Components: comp}
	// anyPF, err := anypb.New(&pbPF1)
	// anyPF1 := pb.PacketFilters{anyPF}
	// anyPFList := []pb.PacketFilters{anyPF}

	pf1 := pb.PacketFilters{Identifier: uint32(pf.Identifier), Direction: pf.Direction, Components: pf.Components}

	// packetFilter := &pb.PacketFilters{
	// 	Pf: &anypb.Any{ // Nested PacketFilter data
	// 		TypeUrl: "type.googleapis.com/PacketFilter", // Example type URL
	// 		Value: []byte(`{
	// 			"identifier": 15,
	// 			"direction": "BIDIRECTIONAL",
	// 			"components": ["Match-all type"]
	// 		}`), // Example JSON
	// 	},
	// }
	pbQos := pb.QosRules{Qosidentifier: qos.QoSIdentifier, Operation: qos.Operation, Dqr: qos.DQR, Precidence: uint32(qos.Precedence), Seg: qos.Segregation, Qfi: qos.QFI, Pf: []*pb.PacketFilters{&pf1}}
	//pbPFList := pb.PacketFilters{anyPF}
	// Create context with timeout
	pbQosList := []*pb.QosRules{&pbQos}

	pbSes := pb.Sessionambr{Iei: int32(ses.IEI), UnitUL: ses.UnitUL, RateUL: int32(ses.RateUL), UnitDL: ses.UnitDL, RateDL: int32(ses.RateDL)}

	esPb := pb.PDUSEstAccModel{Epd: establishmentAccept.ExtendedProtocolDiscriminator, PdusessionID: int32(establishmentAccept.PDUsessionId), Pti: int32(establishmentAccept.PTI), MsgType: establishmentAccept.MessageType, PdusType: establishmentAccept.PduSessionType, SscMode: establishmentAccept.SSCmode, QosIEI: int32(establishmentAccept.QosRuleIEI), Qosrules: pbQosList, Sessionambr: &pbSes}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	anyMessage, err := anypb.New(&esPb)
	if err != nil {
		log.Fatalf("Failed to create Any message: %v", err)
	}

	reqType := "PDU_SESSION_ESTABLISHMENT_ACCEPT"

	// Send the request
	resp, err := client.SendSMData(ctx, &pb.SMDataRequest{
		NasMessage: anyMessage,
		TypeReq:    reqType,
	})
	if err != nil {
		log.Fatalf("Error sending SMData: %v", err)
	}

	// Process the response
	log.Printf("Received response: %v", resp.NasResponse.Value)

	establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
	fmt.Println("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
	fmt.Println("Error while encoding PDU Session Establishment Accept:", err)

	// //PDU Session Release Reject
	// fmt.Println("### Encode PDU Session Release Reject ###")
	// releaseReject := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_RELEASE_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
	// fmt.Printf("Dummy Data for Release Reject: %+v\n", releaseReject)
	// releaseRejectByteArray, err := nas.ReRouteEncode(releaseReject.ExtendedProtocolDiscriminator, releaseReject.MessageType, releaseReject)
	// // releaseRejectByteArray, err := nas.EncodePduSessionReleaseReject(releaseReject)
	// fmt.Println("Encoded Session Release Reject:", releaseRejectByteArray)
	// fmt.Println("Error while encoding PDU Session Release Reject:", err)

	// fmt.Println()

}

func main() {
	Trail()
}

// 	fmt.Println("###Check Reroute###")

// 	val := []byte{126, 0, 65, 121, 0, 13, 1, 130, 246, 16, 0, 0, 0, 0, 0, 0, 0, 0, 16, 46, 4, 240, 240, 240, 240}
// 	l := 25
// 	reg := make([]byte, l)
// 	copy(reg, val)

// 	fmt.Println(reg)
// 	epd, mt, err := nas.Classify(reg)
// 	fmt.Println("err:", err)
// 	res, err := nas.ReRouteDecode(epd, mt, reg)
// 	fmt.Println("err:", err)
// 	fmt.Println(res)

// 	fmt.Println("### Decode Mobility Management Procedures ###")

// 	fmt.Println()

// 	//Registration Request
// 	fmt.Println("### Decode Registration Request ###")
// 	registrationRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinRegistrationRequest", os.O_RDONLY, 0)
// 	// Get the file size
// 	fileInfo, _ := registrationRequestFile.Stat()
// 	fileSize := fileInfo.Size()
// 	// Create a byte array with the size of the file
// 	registrationRequestByteArray := make([]byte, fileSize)
// 	// Read the bytes from the file into the byte array
// 	_, _ = registrationRequestFile.Read(registrationRequestByteArray)

// 	fmt.Println("Byte Array:", registrationRequestByteArray)
// 	registrationPDU, err := nas.DecodeRegistrationRequest(registrationRequestByteArray)
// 	fmt.Printf("Decoded PDU Session Registration Request: %+v\n", registrationPDU)
// 	fmt.Println("Error while decoding PDU Session Registration Request:", err)

// 	fmt.Println()

// 	//UpLink NAS
// 	fmt.Println("### Decode Uplink NAS message ###")
// 	ulNasFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinULNASTransport", os.O_RDONLY, 0)
// 	// Get the file size
// 	fileInfo, _ = ulNasFile.Stat()
// 	fileSize = fileInfo.Size()
// 	// Create a byte array with the size of the file
// 	ulNasByteArray := make([]byte, fileSize)
// 	// Read the bytes from the file into the byte array
// 	_, _ = ulNasFile.Read(ulNasByteArray)
// 	fmt.Println("Byte Array:", ulNasByteArray)
// 	ulNasPDU, err := nas.DecodeULNas(ulNasByteArray)
// 	fmt.Printf("Decoded UL Nas Message:  %+v\n", ulNasPDU)
// 	fmt.Println("Error while decoding UL Nas Message:", err)

// 	fmt.Println()

// 	fmt.Println("### Encoding Mobility Management Procedures ###")

// 	fmt.Println()

// 	// fmt.Println("### Encode Registration Accept ###")
// 	// registrationAccept := nas.RegistrationAcceptModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "REGISTRATION_ACCEPT", RegResult: "3GPP access", Sms: "SMS over NAS not allowed", NssaPerformed: "Network slice-specific authentication and authorization is not to be performed", EmergencyReg: "Not registered for emergency services", RoamingReg: "No additional information"}
// 	// fmt.Printf("Dummy Data for Registartion Accept: %+v\n", registrationAccept)
// 	// registrationAcceptByteArray, err := nas.ReRouteEncode(registrationAccept.ExtendedProtocolDiscriminator, registrationAccept.MessageType, registrationAccept)
// 	// // registrationAcceptByteArray, err := nas.EncodeRegistrationAccept(registrationAccept)
// 	// fmt.Println("Encoded Registartion Accept: ", registrationAcceptByteArray)
// 	// fmt.Println("Error while encoding Registration accept:", err)

// 	fmt.Println()

// 	// fmt.Println("### Encode DL NAS Message ###")
// 	// dlNAS := nas.DLNasModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "DL_NAS_TRANSPORT", PayLoadContainerType: "N1 SM information"}
// 	// fmt.Printf("Dummy Data for DL NAS Transport: %+v\n", dlNAS)
// 	// dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
// 	// fmt.Println("Encoded DL NAS Transport:", dlNasByteArray)
// 	// fmt.Println("Error while encoding DL NAS Transport:", err)

// 	fmt.Println()

// 	fmt.Println("### Decode Session Management Procedures ###")

// 	fmt.Println()

// 	//PDU session establishment request
// 	fmt.Println("### Decode PDU Session Establishment Request ###")
// 	establishmentRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionEstablishmentRequest", os.O_RDONLY, 0)
// 	// Get the file size
// 	fileInfo, _ = establishmentRequestFile.Stat()
// 	fileSize = fileInfo.Size()
// 	// Create a byte array with the size of the file
// 	establishmentRequestByteArray := make([]byte, fileSize)
// 	// Read the bytes from the file into the byte array
// 	_, _ = establishmentRequestFile.Read(establishmentRequestByteArray)
// 	fmt.Println("Byte Array:", establishmentRequestByteArray)
// 	establishmentPDU, err := nas.DecodePduSessionEstablishmentRequest(establishmentRequestByteArray)
// 	fmt.Printf("Decoded PDU Session Establishment Request: %+v\n", establishmentPDU)
// 	fmt.Println("Error while decoding PDU Session Establishment Request:", err)

// 	fmt.Println()

// 	//PDU Session Modification Request
// 	fmt.Println("### Decode PDU Session Modification Request ###")
// 	modificationRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionModificationRequest", os.O_RDONLY, 0)
// 	// Get the file size
// 	fileInfo, _ = modificationRequestFile.Stat()
// 	fileSize = fileInfo.Size()
// 	// Create a byte array with the size of the file
// 	modificationRequestByteArray := make([]byte, fileSize)
// 	// Read the bytes from the file into the byte array
// 	_, _ = modificationRequestFile.Read(modificationRequestByteArray)
// 	fmt.Println("Byte Array:", modificationRequestByteArray)
// 	modificationPDU, err := nas.DecodePduSessionModificationRequest(modificationRequestByteArray)
// 	fmt.Printf("Decoded PDU Session Modification Request: %+v\n", modificationPDU)
// 	fmt.Println("Error while decoding PDU Session Modification Request:", err)

// 	fmt.Println()

// 	//PDU Session Release Request
// 	fmt.Println("### Decode PDU Session Release Request ###")
// 	releaseRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionReleaseRequest", os.O_RDONLY, 0)
// 	// Get the file size
// 	fileInfo, _ = releaseRequestFile.Stat()
// 	fileSize = fileInfo.Size()
// 	// Create a byte array with the size of the file
// 	releaseRequestByteArray := make([]byte, fileSize)
// 	// Read the bytes from the file into the byte array
// 	_, _ = releaseRequestFile.Read(releaseRequestByteArray)
// 	fmt.Println("Byte Array:", releaseRequestByteArray)
// 	releasePDU, err := nas.DecodePduSessionReleaseRequest(releaseRequestByteArray)
// 	fmt.Printf("Decoded PDU Session Release Request: %+v\n", releasePDU)
// 	fmt.Println("Error while decoding PDU Session Release Request:", err)

// 	fmt.Println()

// 	fmt.Println("### Encoding Session Management Procedures")
// 	fmt.Println()

// 	//PDU Session Establishment Accept
// 	// fmt.Println("### Encode PDU Session Establishment Accept ###")
// 	// establishmentAccept := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_ACCEPT", PduSessionType: "IPV4", SSCmode: "SSC_MODE_1"}
// 	// fmt.Printf("Dummy Data for Establishment Accept: %+v\n", establishmentAccept)
// 	// establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
// 	// fmt.Println("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
// 	// fmt.Println("Error while encoding PDU Session Establishment Accept:", err)

// 	// dummyestablishmentAcceptFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionEstablishmentAccept", os.O_RDONLY, 0)
// 	// // Get the file size
// 	// fileInfo, _ = dummyestablishmentAcceptFile.Stat()
// 	// fileSize = fileInfo.Size()
// 	// // Create a byte array with the size of the file
// 	// dummyestablishmentAcceptByteArray := make([]byte, fileSize)
// 	// // Read the bytes from the file into the byte array
// 	// _, _ = dummyestablishmentAcceptFile.Read(dummyestablishmentAcceptByteArray)
// 	// fmt.Println("Byte Array:", dummyestablishmentAcceptByteArray)

// 	fmt.Println()

// 	//PDU Session Establishment Reject
// 	fmt.Println("### Encode PDU Session Establishment Reject ###")
// 	establishmentReject := nas.PduSessionEstablishmentReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
// 	fmt.Printf("Dummy Data for Establishment Reject: %+v\n", establishmentReject)
// 	establishmentRejectByteArray, err := nas.EncodePduSessionEstablishmentReject(establishmentReject)
// 	fmt.Println("Encoded Session Establishment Reject:", establishmentRejectByteArray)
// 	fmt.Println("Error while encoding PDU Session Establishment Reject:", err)

// 	fmt.Println()

// 	//PDU Session Modification Reject
// 	fmt.Println("### Encode PDU Session Modification Reject ###")
// 	modificationReject := nas.PduSessionModificationReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_MODIFICATION_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
// 	fmt.Printf("Dummy Data for Modification Reject: %+v\n", modificationReject)
// 	modificationRejectByteArray, err := nas.EncodePduSessionModificationReject(modificationReject)
// 	fmt.Println("Encoded Session Modification Reject:", modificationRejectByteArray)
// 	fmt.Println("Error while encoding PDU Session Modification Reject:", err)

// 	fmt.Println()

// 	//PDU Session Release Reject
// 	fmt.Println("### Encode PDU Session Release Reject ###")
// 	releaseReject := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_RELEASE_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
// 	fmt.Printf("Dummy Data for Release Reject: %+v\n", releaseReject)
// 	releaseRejectByteArray, err := nas.ReRouteEncode(releaseReject.ExtendedProtocolDiscriminator, releaseReject.MessageType, releaseReject)
// 	// releaseRejectByteArray, err := nas.EncodePduSessionReleaseReject(releaseReject)
// 	fmt.Println("Encoded Session Release Reject:", releaseRejectByteArray)
// 	fmt.Println("Error while encoding PDU Session Release Reject:", err)

// 	fmt.Println()

// 	// dlNasFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MaxPDUSessionEstablishmentAccept", os.O_RDONLY, 0)
// 	// // Get the file size
// 	// fileInfo, _ = dlNasFile.Stat()
// 	// fileSize = fileInfo.Size()
// 	// // Create a byte array with the size of the file
// 	// dlNasByteArray = make([]byte, fileSize)
// 	// // Read the bytes from the file into the byte array
// 	// _, _ = dlNasFile.Read(dlNasByteArray)
// 	// fmt.Println("Byte Array:", dlNasByteArray[41])

// 	fmt.Println("### Encode DL NAS Message ###")
// 	dlNAS := nas.DLNasModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "DL_NAS_TRANSPORT", PayLoadContainerType: "N1 SM information", PayLoadContainer: releaseReject}
// 	fmt.Printf("Dummy Data for DL NAS Transport: %+v\n", dlNAS)
// 	dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
// 	fmt.Println("Encoded DL NAS Transport:", dlNasByteArray)
// 	fmt.Println("Error while encoding DL NAS Transport:", err)

// 	fmt.Println("### Establishment Accept ###")
// 	comp := []string{"Match-all type"}
// 	pf := nas.PacketFilter{Identifier: 15, Direction: "BIDIRECTIONAL", Components: comp}
// 	pflist := []nas.PacketFilterList{pf}
// 	qos := nas.QoSRule{QoSIdentifier: "QRI 1", Operation: "Create new QoS rule", DQR: "DEFAULT_QoS_RULE", Precedence: 255, Segregation: "Segregation not requested", QFI: "QFI 1", PacketFilterList: pflist}
// 	ses := nas.SessionAMBR{IEI: 1, UnitUL: "MULT_1Kbps", RateUL: 30, UnitDL: "MULT_1Kbps", RateDL: 40}
// 	qosList := []nas.QoSRule{qos, qos}
// 	establishmentAccept := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: 1, PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_ACCEPT", PduSessionType: "IPV4", SSCmode: "SSC_MODE_1", AuthorizedQoSRules: qosList, SessionAmbr: ses}
// 	fmt.Printf("Dummy Data for Establishment Accept: %+v\n", establishmentAccept)
// 	establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
// 	fmt.Println("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
// 	fmt.Println("Error while encoding PDU Session Establishment Accept:", err)

// }
