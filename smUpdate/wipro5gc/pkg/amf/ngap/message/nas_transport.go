package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"text/tabwriter"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"w5gc.io/wipro5gcore/pkg/amf/context"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/ngapConvert"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen/asn1rt"
	logger "w5gc.io/wipro5gcore/pkg/amf/ngap/log"

	// "github.com/free5gc/openapi/models"

	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

type RanUe struct {
	RanUeNgapId int64
	AmfUeNgapId int64
}

// var AnType openapi_commn_client.AccessType
// var Snssai openapi_commn_client.Snssai
// var NrLocation *openapiclient.NrLocation
// var PduSessionId int32
// var N1SmContainer []byte

// var RequestType chan (string)

/*
type AmfUe struct {
	RanUe map[models.AccessType]*RanUe
}
*/

func HandleUplinkNASTransport(ran *context.AmfRan, message *asn1gen.NGAPPDU, grpc *grpc.Grpc) {
	fmt.Printf("INSIDE UPLINK NAS TRANSPORT\n")

	/*
		Message Type M YES ignore
		AMF UE NGAP ID M YES reject
		RAN UE NGAP ID M YES reject
		NAS-PDU M YES reject
		User Location Information M YES ignore
	*/

	var amfUeNgapId asn1gen.AMFUENGAPID
	var ranUeNgapId asn1gen.RANUENGAPID
	var nasPdu *asn1gen.NASPDU = &asn1gen.NASPDU{}
	var userLocationInformation *asn1gen.UserLocationInformation = &asn1gen.UserLocationInformation{}

	if ran == nil {
		fmt.Println("RAN is nil")
		return
	}

	if message == nil {
		fmt.Println("NGAP Message is nil")
		return
	}

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		fmt.Println("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var initialUeMessage *asn1gen.InitialUEMessage = &asn1gen.InitialUEMessage{}
	_, err := asn1gen.Unmarshal(x, initialUeMessage)
	if err != nil {
		println("unmarshalling failed for initial ue message type")
		return
	}
	initialUeMessageIEs := initialUeMessage.ProtocolIEs
	for i := 0; i < len(initialUeMessageIEs); i++ {
		ie := initialUeMessageIEs[i]
		// fmt.Printf("ie : %+v", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			amfuengapid_val := []byte(ie.Value)
			if amfuengapid_val == nil {
				fmt.Println("AMF UE NGAP ID is nil")
				return
			}
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				fmt.Println("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId > math.MaxUint32 {
				fmt.Println("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdNASPDU:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), nasPdu)
			if err != nil {
				fmt.Println("unmarshalling of NAS PDU failed : ", err)
				return
			}
			if *nasPdu == nil {
				fmt.Println("NAS PDU is nil")
				return
			}

		case asn1gen.Asn1vIdUserLocationInformation:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), userLocationInformation)
			if err != nil {
				fmt.Println("unmarshalling of User Location Information failed : ", err)
				return
			}
			if (*userLocationInformation == asn1gen.UserLocationInformation{}) {
				fmt.Println("User Location Information is nil")
				return
			}
			fmt.Printf("USER LOCATION INFORMATION : %+v", userLocationInformation)
			fmt.Printf("USER LOCATION INFORMATION NR : %+v", userLocationInformation.U.UserLocationInformationNR)
		}
	}

	/*nasPduFile, err := os.Create("nasPdu")
	if err != nil {
		fmt.Println("error in creating NAS PDU file : ", err)
		return
	}
	defer nasPduFile.Close()
	_, err = nasPduFile.Write(*nasPdu)
	if err != nil {
		fmt.Println("error in writing to NAS PDU file : ", err)
		return
	}*/
	//TO DO : send nasPduFile to another function to decide type of nas message
	//if type is UL NAS Transport then call the below code. then forward the enclosed NAS-SM message to SMF
	/*ulNasModel, err := nas.DecodeULNas(nasPduFile)
	if err != nil {
		fmt.Println("error in decoding UL NAS PDU : ", err)
		return
	}*/
	/*fmt.Println("before nas classify")
	nasPduByteArray := []byte(*nasPdu)
	epd, messageType, err := nas.Classify(nasPduByteArray)
	if err != nil {
		fmt.Println("error while classifying NAS Message : ", err)
		return
	}
	fmt.Println("epd= ", epd, " , messageType= ", messageType)
	fmt.Println("NAS PDU: ", *nasPdu)
	decodedUlNas, err := nas.DecodeULNas(*nasPdu)
	if err != nil {
		fmt.Println("error while decoding UL NAS Transport message : ", err)
		return
	}
	fmt.Printf("%+v\n", decodedUlNas)
	epd, messageType, err = nas.Classify(nasPduByteArray[6:12])
	if err != nil {
		fmt.Println("error while classifying NAS message : ", err)
		return
	}
	nasDecodedMsg, err := nas.ReRoute(epd, messageType, nasPduByteArray[6:12])
	if err != nil {
		fmt.Println("error while decoding NAS SM message : ", err)
		return
	}
	fmt.Println("nas decoded message : ", nasDecodedMsg)
	if err != nil {
		fmt.Println("error while decoding UL NAS Transport message : ", err)
		return
	}*/
	/*INSERT NEW CODE HERE*/
	// apiclient.AnType = ran.AnType

	// fmt.Println("AN TYPE FROM AMF/UERANSIM: ", apiclient.AnType)

	// apiclient.Snssai = ran.SupportedTAList[0].SNssaiList[0]

	// fmt.Println("SNSSAI FROM AMF/UERANSIM: ", apiclient.Snssai)

	// timestamp := binary.BigEndian.Uint32(*(userLocationInformation.U.UserLocationInformationNR.TimeStamp))
	// ueLocationTimestamp := time.Unix(int64(timestamp), 0)

	/*plmnId := openapiclient.PlmnId(ngapConvert.PlmnIdToModels(userLocationInformation.U.UserLocationInformationNR.TAI.PLMNIdentity))
	tac := hex.EncodeToString(userLocationInformation.U.UserLocationInformationNR.TAI.TAC)
	nrCellId := hex.EncodeToString(userLocationInformation.U.UserLocationInformationNR.NRCGI.NRCellIdentity.Bytes)
	timestamp := binary.BigEndian.Uint32(*(userLocationInformation.U.UserLocationInformationNR.TimeStamp))
	ueLocationTimestamp := time.Unix(int64(timestamp), 0)
	apiclient.NrLocation = &openapiclient.NrLocation{
		Tai:                 *openapiclient.NewTai(plmnId, tac),
		Ncgi:                *openapiclient.NewNcgi(plmnId, nrCellId),
		UeLocationTimestamp: &ueLocationTimestamp,
	}

	fmt.Println("NR LOCATION FROM AMF/UERANSIM: ", apiclient.NrLocation)*/

	/*var i int
	if pduSessionReq, ok := nasDecodedMsg.(nas.PduSessionEstablishmentRequest); ok {
		logger.AppLog.Println("Type assertion of NAS Decoded Message to PDU Session Establishment Request passed.")
		PduSessionIdStr := (pduSessionReq.PDUsessionId)[4]
		fmt.Println("PduSessionIdStr : ", PduSessionIdStr)
		i, err = strconv.Atoi(string(PduSessionIdStr))
		if err != nil {
			logger.AppLog.Errorln("Cannot convert PDU Session ID string to int: ", err)
		} else {
			fmt.Println("PduSessionId : ", i)
		}
	}*/

	// fmt.Println("PDU SESSION ID FROM AMF/UERANSIM: ", apiclient.PduSessionId)

	// apiclient.N1SmContainer = nasPduByteArray[6:]

	// fmt.Println("N1 SM CONTAINER FROM AMF/UERANSIM: ", apiclient.N1SmContainer)
	//needs to be sent from NAS instead of NGAP
	/*y := &create_sm_context_grpc.CreateSmContextDataFromNasMod{
		AnType: string(ran.AnType),
		Snssai: &create_sm_context_grpc.Snssai{
			Sst: ran.SupportedTAList[0].SNssaiList[0].Sst,
			Sd:  *ran.SupportedTAList[0].SNssaiList[0].Sd,
		},
		NrLocation: &create_sm_context_grpc.NrLocation{
			Tai: &create_sm_context_grpc.Tai{
				PlmnId: &create_sm_context_grpc.PlmnId{
					//harcoded for now
					Mcc: "286",
					Mnc: "01",
				},
				Tac: string(userLocationInformation.U.UserLocationInformationNR.TAI.TAC),
			},
			Ncgi: &create_sm_context_grpc.Ncgi{
				PlmnId: &create_sm_context_grpc.PlmnId{
					//harcoded for now
					Mcc: "286",
					Mnc: "01",
				},
				NrCellId: string(userLocationInformation.U.UserLocationInformationNR.NRCGI.NRCellIdentity.Bytes),
			},
			UeLocationTimestamp: timestamppb.New(ueLocationTimestamp),
		},
		PduSessionId: int32(i),
		//ASK MOUNIKA
		N1SmContainer: nasPduByteArray[4:],
	}
	(*grpc).SendDataForCreateSmContext(y)*/
	timestamp := binary.BigEndian.Uint32(*(userLocationInformation.U.UserLocationInformationNR.TimeStamp))
	ueLocationTimestamp := time.Unix(int64(timestamp), 0)
	createSmContext := pb.CreateSmContext{
		AnType: string(ran.AnType),
		Snssai: &pb.Snssai{
			Sst: ran.SupportedTAList[0].SNssaiList[0].Sst,
			Sd:  *ran.SupportedTAList[0].SNssaiList[0].Sd,
		},
		NrLocation: &pb.NrLocation{
			Tai: &pb.Tai{
				PlmnId: &pb.PlmnId{
					//harcoded for now
					Mcc: "286",
					Mnc: "01",
				},
				Tac: string(userLocationInformation.U.UserLocationInformationNR.TAI.TAC),
			},
			Ncgi: &pb.Ncgi{
				PlmnId: &pb.PlmnId{
					//harcoded for now
					Mcc: "286",
					Mnc: "01",
				},
				NrCellId: string(userLocationInformation.U.UserLocationInformationNR.NRCGI.NRCellIdentity.Bytes),
			},
			UeLocationTimestamp: timestamppb.New(ueLocationTimestamp),
		},
		N1SmContainer: *nasPdu,
		//harcoded for now
		PduSessionId: 1,
	}

	// serializeMsg, _ := proto.Marshal(&createSmContext)
	serializeMsg, err := anypb.New(&createSmContext)
	if err != nil {
		logger.AppLog.Errorln("Failed to marshal create sm context data : ", err)
		return
	}
	req := &pb.DataRequest{
		// Data: &anypb.Any{
		// 	TypeUrl: "w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb/CreateSmContext",
		// 	Value:   serializeMsg,
		// },
		Data:    serializeMsg,
		ReqType: "Uplink NAS",
	}
	res := (*grpc).SendData(req)
	if res.Error != "" {
		logger.AppLog.Errorln(res.Error)
	} else {
		//TO DO
		/*if err := SendPduSessionResourceSetupRequest(ran, res.NasPdu); err != nil {
			logger.AppLog.Errorln("Error in sending Pdu Session Resource Setup Request : ", err)
		}*/
	}
	// fmt.Println("AFTER SENDING REQUEST TYPE IN CHANNEL")
	/*establishmentAccept := nas.PduSessionEstablishmentAccept{
		ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES",
		PDUsessionId:                  "VAL_1",
		PTI:                           1,
		MessageType:                   "PDU_SESSION_ESTABLISHMENT_ACCEPT",
		PduSessionType:                "IPV4",
		SSCmode:                       "SSC_MODE_1",
		SessionAmbr:                   "MULT_1Kbps",
	}
	fmt.Printf("Dummy Data for Establishment Accept: %+v\n", establishmentAccept)
	establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
	if err != nil {
		fmt.Println("error while encoding PDU Session Establishment Accept: ", err)
		return
	}*/
	/*establishmentAcceptByteArray, err := os.ReadFile("/home/vboxuser/asn1c-v772-for-vm/golang/sample_per/ngap/src/nas/testFiles/MinPDUSessionEstablishmentAccept")
	if err != nil {
		fmt.Println("error in reading pdu session establishment accept file : ", err)
		return
	}
	*/
	/*establishmentAcceptByteArray := []byte{46, 1, 1, 194, 17, 0, 9, 1, 0, 6, 49, 49, 1, 1, 255, 9, 6, 1, 3, 232, 1, 3, 232, 41, 5, 1, 10, 62, 0, 1, 34, 4, 1, 17, 34, 51, 121, 0, 6, 9, 32, 65, 1, 1, 9, 123, 0, 27, 128, 0, 13, 4, 8, 8, 8, 8, 0, 3, 16, 32, 1, 72, 96, 72, 96, 0, 0, 0, 0, 0, 0, 0, 0, 136, 136, 37, 10, 9, 105, 110, 116, 101, 114, 110, 101, 116, 50}
	fmt.Println("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
	dlNAS := nas.DLNasModel{
		ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES",
		SecurityHeaderType:            "Plain 5GS NAS message",
		MessageType:                   "DL_NAS_TRANSPORT",
		PayLoadContainerType:          "N1 SM information",
	}
	dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
	dlNasByteArray[3] = 1
	if err != nil {
		fmt.Println("Error while encoding DL NAS Transport:", err)
		return
	}
	fmt.Println("Encoded DL NAS Transport:", dlNasByteArray)
	establishmentAcceptByteArrayLen := len(establishmentAcceptByteArray)
	intermediateArray := []byte{0, byte(establishmentAcceptByteArrayLen)}
	dlNasByteArray = append(dlNasByteArray, intermediateArray...)
	nasPduByteArray1 := append(dlNasByteArray, establishmentAcceptByteArray...)
	nasPdu1 := asn1gen.NASPDU(nasPduByteArray1)
	fmt.Println("nas message : dl nas + pdu session est. accept = ", nasPdu1)
	if err := SendPduSessionResourceSetupRequest(ran, nasPdu1); err != nil {
		fmt.Println("Error in sending Pdu Session Resource Setup Request : ", err)
		return
	}*/
}

func HandleInitialUEMessage(ran *context.AmfRan, message *asn1gen.NGAPPDU, grpc *grpc.Grpc) {
	/* Message-Type M YES ignore
	RAN-UE-NGAP-ID M YES reject
	NAS-PDU M YES reject
	User Location Information M YES reject
	RRC Establishment Cause M YES ignore
	5G-S-TMSI O YES reject
	AMF Set ID O YES ignore
	UE Context Request O ENUMERATED(requested, ...) YES ignore
	Allowed NSSAI O YES reject
	Source to Target AMF Information Reroute O YES ignore
	Selected PLMN Identity O PLMN Identity YES ignore */
	var ranUeNgapId asn1gen.RANUENGAPID
	var nasPdu *asn1gen.NASPDU = &asn1gen.NASPDU{}
	var userLocationInformation *asn1gen.UserLocationInformation = &asn1gen.UserLocationInformation{}
	var rrcEstablishmentCause asn1gen.RRCEstablishmentCause
	var fiveGSTMSI asn1gen.FiveGSTMSI = asn1gen.FiveGSTMSI{}
	var amfSetId asn1gen.AMFSetID = asn1gen.AMFSetID{}
	var ueContextRequest asn1gen.UEContextRequest
	var allowedNssai asn1gen.AllowedNSSAI = asn1gen.AllowedNSSAI{}
	var sourceToTargetAMFInformationReroute asn1gen.SourceToTargetAMFInformationReroute = asn1gen.SourceToTargetAMFInformationReroute{}
	var selectedPlmnIdentity *asn1gen.PLMNIdentity = &asn1gen.PLMNIdentity{}

	if ran == nil {
		logger.AppLog.Error("RAN is nil")
		return
	}

	if message == nil {
		logger.AppLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		logger.AppLog.Error("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var initialUeMessage *asn1gen.InitialUEMessage = &asn1gen.InitialUEMessage{}
	_, err := asn1gen.Unmarshal(x, initialUeMessage)
	if err != nil {
		logger.AppLog.Error("unmarshalling failed for initial ue message type: ", err)
		return
	}
	initialUeMessageIEs := initialUeMessage.ProtocolIEs
	for i := 0; i < len(initialUeMessageIEs); i++ {
		ie := initialUeMessageIEs[i]
		logger.AppLog.Printf("IE : %+v\n", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				logger.AppLog.Error("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId > math.MaxUint32 {
				logger.AppLog.Error("RAN UE NGAP ID has an invalid value")
				return
			}

		case asn1gen.Asn1vIdNASPDU:
			logger.AppLog.Println("[NGAP] Decode NAS PDU")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), nasPdu)
			if err != nil {
				logger.AppLog.Error("unmarshalling of NAS PDU failed : ", err)
				return
			}
			if *nasPdu == nil {
				logger.AppLog.Error("NAS PDU is nil")
				return
			}

		case asn1gen.Asn1vIdUserLocationInformation:
			logger.AppLog.Println("[NGAP] Decode User Location Information")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), userLocationInformation)
			if err != nil {
				logger.AppLog.Error("unmarshalling of User Location Information failed : ", err)
				return
			}
			if (*userLocationInformation == asn1gen.UserLocationInformation{}) {
				logger.AppLog.Error("User Location Information is nil")
				return
			}

		case asn1gen.Asn1vIdRRCEstablishmentCause:
			logger.AppLog.Println("[NGAP] Decode RRC Establishment Cause")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &rrcEstablishmentCause)
			if err != nil {
				logger.AppLog.Error("unmarshalling of RRC Establishment Cause failed : ", err)
				return
			}
			if rrcEstablishmentCause > 11 {
				logger.AppLog.Error("RRC Establishment Cause is nil")
				return
			}

		case asn1gen.Asn1vIdFiveGSTMSI:
			logger.AppLog.Println("[NGAP] Decode 5G TMSI")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), fiveGSTMSI)
			if err != nil {
				logger.AppLog.Info("unmarshalling of 5GSTMSI failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(fiveGSTMSI, asn1gen.FiveGSTMSI{})) {
				logger.AppLog.Info("5GSTMSI is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdAMFSetID:
			logger.AppLog.Println("[NGAP] Decode AMF Set ID")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), amfSetId)
			if err != nil {
				logger.AppLog.Info("unmarshalling of AMF Set ID failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(amfSetId, asn1gen.AMFSetID{})) {
				logger.AppLog.Info("AMF Set ID is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdUEContextRequest:
			logger.AppLog.Println("[NGAP] Decode UE Context Request")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ueContextRequest)
			if err != nil {
				logger.AppLog.Info("unmarshalling of UE Context Request failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			//UeContextRequest is an enum whose value is 0,1
			if ueContextRequest != 0 && ueContextRequest != 1 {
				logger.AppLog.Info("UE Context Request has an invalid value")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdAllowedNSSAI:
			logger.AppLog.Println("[NGAP] Decode Allowed NSSAI")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), allowedNssai)
			if err != nil {
				logger.AppLog.Info("unmarshalling of Allowed NSSAI failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if len(allowedNssai) == 0 {
				logger.AppLog.Info("No Allowed NSSAI present")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdSourceToTargetAMFInformationReroute:
			logger.AppLog.Println("[NGAP] Decode SourceToTargetAMFInformationReroute")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), sourceToTargetAMFInformationReroute)
			if err != nil {
				logger.AppLog.Info("unmarshalling of Source to Target AMF Information Reroute failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(sourceToTargetAMFInformationReroute, asn1gen.SourceToTargetAMFInformationReroute{})) {
				logger.AppLog.Info("Source to Target AMF Information Reroute is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdSelectedPLMNIdentity:
			logger.AppLog.Println("[NGAP] Decode Selected PLMN Identity")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), selectedPlmnIdentity)
			if err != nil {
				logger.AppLog.Info("unmarshalling of Selected PLMN Identity failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if *selectedPlmnIdentity == nil {
				logger.AppLog.Info("Selected PLMN Identity is nil")
				//commented return as this IE presence is optional
				//return
			}
		}
	}

	//RAN-UE-NGAP-ID
	//to do : ask if a ranUe struct should be created?

	/*nasPduFile, err := os.Create("nasPdu")
	if err != nil {
		logger.AppLog.Error("error in creating NAS PDU file : ", err)
	}
	defer nasPduFile.Close()
	logger.AppLog.Println("NAS PDU : ", *nasPdu)
	_, err = nasPduFile.Write(*nasPdu)
	if err != nil {
		logger.AppLog.Error("error in writing to NAS PDU file : ", err)
	} else {
		logger.AppLog.Info("written nas pdu data to file successfully")
	}
	nasPduFile.Seek(0, 0)*/

	//checking if written to file successfully
	/*
		var b = make([]byte, len(*nasPdu))
		n, err := nasPduFile.Read(b)
		fmt.Println(n, err, b)
		nasPduFile.Seek(0, 0)
	*/
	nasMsg := pb.NasMessage{
		NasPdu: *nasPdu,
	}
	// serializeNasMsg, _ := proto.Marshal(&nasMsg)
	serializeNasMsg, err := anypb.New(&nasMsg)
	if err != nil {
		logger.AppLog.Errorln("Error in marshaling nas message : ", err)
	}
	regReq := &pb.DataRequest{
		// Data: &anypb.Any{
		// 	TypeUrl: "w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb/NasMessage",
		// 	Value:   serializeNasMsg,
		// },
		Data:    serializeNasMsg,
		ReqType: "Registration Request",
	}
	res := (*grpc).SendData(regReq)
	if res.Error != "" {
		logger.AppLog.Error(res.Error)
	} else {
		fmt.Println("UE Context Request value : ", ueContextRequest)
		if ueContextRequest == 0 {
			SendInitialContextSetupRequest(ran, ranUeNgapId, res.NasPdu)
		} else {
			SendDownlinkNASTransport(ran, ranUeNgapId, res.NasPdu)
		}
	}
	/*err = nas.DecodeInitialUE(nasPduFile)
	if err != nil {
		logger.AppLog.Error(err)
	}
	regAcceptFile := nas.EncodeRegistrationAcceptMain()
	regAcceptFile2, err := os.Open(regAcceptFile.Name())
	if err != nil {
		logger.AppLog.Error(err)
	}
	nasMsg, err := os.ReadFile(regAcceptFile2.Name())
	logger.AppLog.Println("nasmsg received from nas module : ", nasMsg)
	if err != nil {
		logger.AppLog.Error(err)
	} else {
		fmt.Println("UE Context Request value : ", ueContextRequest)
		if ueContextRequest == 0 {
			SendInitialContextSetupRequest(ran, ranUeNgapId, nasMsg)
		} else {
			SendDownlinkNASTransport(ran, ranUeNgapId, nasMsg)
		}
	}*/
	//stub for NG Reset
	//logger.AppLog.Println("SENDING NG RESET MESSAGE FROM AMF TO NG RAN")
	//SendNGResetMessage(ran)
}

// to do: which all parameters to pass? sending only ranuengapid and nasmsg right now
func SendInitialContextSetupRequest(ran *context.AmfRan, ranUeNgapId asn1gen.RANUENGAPID, nasMsg asn1gen.NASPDU) {
	/*
		IE/Group Name											Presence	Range	IE type and reference	Semantics description	Criticality	Assigned Criticality
		Message Type											M					9.3.1.1		YES	reject
		AMF UE NGAP ID											M					9.3.3.1		YES	reject
		RAN UE NGAP ID											M					9.3.3.2		YES	reject
		Old AMF													O					AMF Name 9.3.3.21		YES	reject
		UE Aggregate Maximum Bit Rate							C-ifPDUsessionResourceSetup		9.3.1.58		YES	reject
		Core Network Assistance Information for RRC INACTIVE	O		9.3.1.15		YES	ignore
		GUAMI													M		9.3.3.3		YES	reject
		PDU Session Resource Setup Request List								0..1			YES	reject
		>PDU Session Resource Setup Request Item							1..<maxnoofPDUSessions>			-
		>>PDU Session ID										M		9.3.1.50		-
		>>PDU Session NAS-PDU	O		NAS-PDU 9.3.3.4		-
		>>S-NSSAI 	M		9.3.1.24		-
		>>PDU Session Resource Setup Request Transfer M		OCTET STRING	Containing the PDU Session Resource Setup Request Transfer IE specified in subclause 9.3.4.1.	-
		Allowed NSSAI	M		9.3.1.31	Indicates the S-NSSAIs permitted by the network	YES	reject
		UE Security Capabilities	M		9.3.1.86		YES	reject
		Security Key	M		9.3.1.87		YES	reject
		Trace Activation	O		9.3.1.14		YES	ignore
		Mobility Restriction List	O		9.3.1.85		YES	ignore
		UE Radio Capability	O		9.3.1.74		YES	ignore
		Index to RAT/Frequency Selection Priority	O		9.3.1.61		YES	ignore
		Masked IMEISV	O		9.3.1.54		YES	ignore
		NAS-PDU	O		9.3.3.4		YES	ignore
		Emergency Fallback Indicator	O		9.3.1.26		YES	reject
		RRC Inactive Transition Report Request	O		9.3.1.91		YES	ignore
		UE Radio Capability for Paging	O		9.3.1.68		YES	ignore
		Redirection for Voice EPS Fallback 	O		9.3.1.116		YES	ignore
		Location Reporting Request Type	O		9.3.1.65		YES	ignore
		CN Assisted RAN Parameters Tuning	O		9.3.1.119		YES	ignore
	*/

	amfUeNgapId := asn1gen.AMFUENGAPID(10)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		fmt.Println("error in marshaling amf ue ngap id in initial context setup request : ", err)
		return
	}

	// ranUeNgapId = asn1gen.RANUENGAPID(1)
	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		fmt.Println("error in marshaling ran ue ngap id in initial context setup request : ", err)
		return
	}

	plmnid := openapi_commn_client.PlmnId{
		Mcc: "286",
		Mnc: "01",
	}
	guamiModel := openapi_commn_client.Guami{
		PlmnId: plmnid,
		AmfId:  "218A9E",
	}
	regionId, setId, prtId := ngapConvert.AmfIdToNgap(guamiModel.AmfId)
	guami := asn1gen.GUAMI{
		PLMNIdentity: ngapConvert.PlmnIdToNgap(plmnid),
		AMFRegionID:  asn1gen.AMFRegionID(regionId),
		AMFSetID:     asn1gen.AMFSetID(setId),
		AMFPointer:   asn1gen.AMFPointer(prtId),
	}
	encodedGuami, err := asn1gen.Marshal(guami)
	if err != nil {
		fmt.Println("error in marshaling guami in initial context setup request : ", err)
		return
	}

	sd := asn1gen.SD(asn1rt.OctetString(ngapConvert.HexToBitString("000001", 24).Bytes))
	allowedNssai := asn1gen.AllowedNSSAI{
		asn1gen.AllowedNSSAIItem{
			SNSSAI: asn1gen.SNSSAI{
				SST: asn1gen.SST("1"),
				SD:  &sd,
			},
		},
	}
	encodedAllowedNssai, err := asn1gen.Marshal(allowedNssai)
	if err != nil {
		fmt.Println("error in marshaling allowed nssai in initial context setup request : ", err)
		return
	}

	q := asn1rt.BitString{
		Bytes:     []byte{0, 0, 0, 0},
		BitLength: 16,
	}
	ueSecurityCapabilities := asn1gen.UESecurityCapabilities{
		NRencryptionAlgorithms:             asn1gen.NRencryptionAlgorithms(q),
		NRintegrityProtectionAlgorithms:    asn1gen.NRintegrityProtectionAlgorithms(q),
		EUTRAencryptionAlgorithms:          asn1gen.EUTRAencryptionAlgorithms(q),
		EUTRAintegrityProtectionAlgorithms: asn1gen.EUTRAintegrityProtectionAlgorithms(q),
	}
	encodedUeSecurityCapabilities, err := asn1gen.Marshal(ueSecurityCapabilities)
	if err != nil {
		fmt.Println("error in marshaling ue security capabilities in initial context setup request : ", err)
		return
	}

	securityKey := asn1gen.SecurityKey(ngapConvert.HexToBitString("0000000000000000000000000000000000000000000000000000000000000000", 256))
	encodedSecurityKey, err := asn1gen.Marshal(securityKey)

	encodedNasMsg, err := asn1gen.Marshal(nasMsg)
	if err != nil {
		fmt.Println("error in marshaling nas msg in initial context setup request : ", err)
		return
	}

	initialContextSetupRequest := asn1gen.InitialContextSetupRequest{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdAMFUENGAPID,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedAmfUeNgapId,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdRANUENGAPID,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedRanUeNgapId,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdGUAMI,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedGuami,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdAllowedNSSAI,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedAllowedNssai,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdUESecurityCapabilities,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedUeSecurityCapabilities,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdSecurityKey,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedSecurityKey,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdNASPDU,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedNasMsg,
			},
		},
	}
	encodedInitialContextSetupRequest, err := asn1gen.Marshal(initialContextSetupRequest)
	if err != nil {
		fmt.Println("error in marshaling initial context setup request : ", err)
		return
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: asn1gen.NGAPPDUInitiatingMessageTAG,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			InitiatingMessage: &asn1gen.InitiatingMessage{
				ProcedureCode: asn1gen.Asn1vIdInitialContextSetup,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedInitialContextSetupRequest,
			},
		},
	}

	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		fmt.Println("error in marshaling ngap msg : ", err)
		return
	}

	if ran == nil {
		fmt.Println(errors.New("RAN is nil"))
		return
	}

	if len(encodedNgapMsg) == 0 {
		fmt.Println(errors.New("packet length is 0"))
		return
	}

	if ran.Conn == nil {
		fmt.Println(errors.New("RAN address is nil"))
		return
	}

	n, err := ran.Conn.Write(encodedNgapMsg)
	if err != nil {
		err := "Write error : " + err.Error()
		fmt.Println(errors.New(err))
		return
	} else {
		logger.AppLog.Println("[INITIAL CONTEXT SETUP] Wrote ", n, " bytes")
	}
	return
}

func HandleInitialContextSetupResponse(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	if ran == nil {
		fmt.Println("RAN is nil")
		return
	}

	if message == nil {
		fmt.Println("NGAP Message is nil")
		return
	}
	var ranUeNgapId asn1gen.RANUENGAPID
	var amfUeNgapId asn1gen.AMFUENGAPID
	criticalityDiagnostics := asn1gen.CriticalityDiagnostics{}
	pduSessionResourceSetupListSuRes := asn1gen.PDUSessionResourceSetupListSURes{}
	pduSessionResourceFailedToSetupListSuRes := asn1gen.PDUSessionResourceFailedToSetupListSURes{}

	successfulOutcome := message.U.SuccessfulOutcome
	if successfulOutcome == nil {
		fmt.Println("Successful Outcome message is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	var initialContextSetupResponse *asn1gen.InitialContextSetupResponse = &asn1gen.InitialContextSetupResponse{}
	_, err := asn1gen.Unmarshal(x, initialContextSetupResponse)
	if err != nil {
		println("unmarshalling failed for initial context setup response type")
		return
	}
	initialContextSetupResponseIEs := initialContextSetupResponse.ProtocolIEs
	for i := 0; i < len(initialContextSetupResponseIEs); i++ {
		ie := initialContextSetupResponseIEs[i]
		if reflect.ValueOf(ie).IsZero() {
			fmt.Printf("ie is null\n")
			return
		} else {
			fmt.Printf("ie : %+v", ie)
		}
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId < 0 || amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				fmt.Println("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId < 0 || ranUeNgapId > math.MaxUint32 {
				fmt.Println("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceSetupListSuRes)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling pduSessionResourceSetupListSuRes : ", err)
			}
			if len(pduSessionResourceSetupListSuRes) == 0 {
				logger.AppLog.Error("pduSessionResourceSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceFailedToSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceFailedToSetupListSuRes)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling pduSessionResourceFailedToSetupListSuRes : ", err)
			}
			if len(pduSessionResourceFailedToSetupListSuRes) == 0 {
				logger.AppLog.Error("pduSessionResourceFailedToSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &criticalityDiagnostics)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling criticalityDiagnostics")
			}

		}
	}
	//todo: process the decoded values

}

func HandleInitialContextSetupFailure(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	if ran == nil {
		fmt.Println("RAN is nil")
		return
	}

	if message == nil {
		fmt.Println("NGAP Message is nil")
		return
	}
	var ranUeNgapId asn1gen.RANUENGAPID
	var amfUeNgapId asn1gen.AMFUENGAPID
	cause := asn1gen.Cause{}
	criticalityDiagnostics := asn1gen.CriticalityDiagnostics{}
	pduSessionResourceFailedToSetupListSuRes := asn1gen.PDUSessionResourceFailedToSetupListSURes{}
	unsuccessfulOutcome := message.U.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		fmt.Println("Unsuccessful Outcome message is nil")
		return
	}
	x := []byte(unsuccessfulOutcome.Value)
	var initialContextSetupFailure *asn1gen.InitialContextSetupFailure = &asn1gen.InitialContextSetupFailure{}
	_, err := asn1gen.Unmarshal(x, initialContextSetupFailure)
	if err != nil {
		println("unmarshalling failed for initial context setup failure type")
		return
	}
	initialContextSetupFailureIEs := initialContextSetupFailure.ProtocolIEs
	for i := 0; i < len(initialContextSetupFailureIEs); i++ {
		ie := initialContextSetupFailureIEs[i]
		if reflect.ValueOf(ie).IsZero() {
			fmt.Printf("ie is null\n")
			return
		} else {
			fmt.Printf("ie : %+v", ie)
		}
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId < 0 || amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				fmt.Println("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				fmt.Println("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId < 0 || ranUeNgapId > math.MaxUint32 {
				fmt.Println("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceFailedToSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceFailedToSetupListSuRes)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling pduSessionResourceFailedToSetupListSuRes : ", err)
			}
			if len(pduSessionResourceFailedToSetupListSuRes) == 0 {
				logger.AppLog.Error("pduSessionResourceFailedToSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &criticalityDiagnostics)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling criticalityDiagnostics")
			}

		case asn1gen.Asn1vIdCause:
			logger.AppLog.Println("[NGAP] Decode IE Cause")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), cause)
			if err != nil {
				logger.AppLog.Error("Error in unmarshaling Cause : ", err)
				return
			}
			if cause == (asn1gen.Cause{}) {
				logger.AppLog.Error("Cause is nil")
				return
			}

		}
	}
	//todo: process the decoded values

}

func SendDownlinkNASTransport(ran *context.AmfRan, ranUeNgapId asn1gen.RANUENGAPID, nasMsg asn1gen.NASPDU) error {
	/*
		Message-Type M YES ignore
		AMF-UE-NGAP-ID M YES reject
		RAN-UE-NGAP-ID M YES reject
		Old-AMF O YES reject
		Ran-PAging-Priority O YES ignore
		NAS-PDU M YES reject
		Mobility-Restriction-List O YES ignore
		Index-to-RAT/Frequency-Selection-Priority O YES ignore
		UE-Aggregate-Maximum-Bit-Rate O YES ignore
		Allowed NSSAI O YES reject
	*/

	/*
		var amfUe AmfUe
		ranUe := make(map[models.AccessType]*RanUe)
		ranUe[models.AccessType__3_GPP_ACCESS] = &RanUe{
			RanUeNgapId: 2,
			AmfUeNgapId: 1,
		}
		amfUe.RanUe = ranUe

		amfUeNgapId := asn1gen.AMFUENGAPID(amfUe.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId)
	*/
	amfUeNgapId := asn1gen.AMFUENGAPID(10)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		return err
	}

	logger.AppLog.Println("ranUeNgapId:", ranUeNgapId)

	// ranUeNgapId2 := asn1gen.RANUENGAPID(20)
	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		return err
	}

	var nasPdu asn1gen.NASPDU = nasMsg
	encodedNasPdu, err := asn1gen.Marshal(nasPdu)
	if err != nil {
		return err
	}

	var downlinkNasTransport asn1gen.DownlinkNASTransport = asn1gen.DownlinkNASTransport{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdAMFUENGAPID,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedAmfUeNgapId,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdRANUENGAPID,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedRanUeNgapId,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdNASPDU,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedNasPdu,
			},
		},
	}

	encodedDownlinkNasTransport, err := asn1gen.Marshal(downlinkNasTransport)
	if err != nil {
		return err
	}

	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = asn1gen.NGAPPDUInitiatingMessageTAG
	ngapMsg.U.InitiatingMessage = &asn1gen.InitiatingMessage{
		ProcedureCode: asn1gen.Asn1vIdDownlinkNASTransport,
		Criticality:   asn1gen.CriticalityIgnore,
		Value:         encodedDownlinkNasTransport,
	}

	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		return err
	}

	if ran == nil {
		return errors.New("RAN is nil")
	}

	if len(encodedNgapMsg) == 0 {
		return errors.New("packet length is 0")
	}

	if ran.Conn == nil {
		return errors.New("RAN address is nil")
	}

	n, err := ran.Conn.Write(encodedNgapMsg)
	if err != nil {
		err := "Write error : " + err.Error()
		return errors.New(err)
	} else {
		logger.AppLog.Println("[DOWNLINK NAS TRANSPORT] Wrote ", n, " bytes")
	}
	return nil
}

func HandleErrorIndicationFromRan(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	fmt.Println()
	fmt.Println("---------- ERROR INDICATION FROM RAN -----------")
	if ran == nil {
		fmt.Println("RAN is nil")
		return
	}

	if message == nil {
		fmt.Println("NGAP Message is nil")
		return
	}

	var ranUeNgapId *asn1gen.RANUENGAPID
	var amfUeNgapId *asn1gen.AMFUENGAPID
	var cause *asn1gen.Cause
	var criticalityDiagnostics *asn1gen.CriticalityDiagnostics

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		fmt.Println("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	errorInd := asn1gen.ErrorIndication{}
	_, err := asn1gen.Unmarshal(x, &errorInd)
	if err != nil {
		fmt.Println("error in unmarshaling error indication received from gNB : ", err)
		return
	}

	for i := 0; i < len(errorInd.ProtocolIEs); i++ {
		ie := errorInd.ProtocolIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal(ie.Value, &ranUeNgapId)
			if err != nil {
				fmt.Println("error in unmarshaling ran ue ngap id in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal(ie.Value, &amfUeNgapId)
			if err != nil {
				fmt.Println("error in unmarshaling amf ue ngap id in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdCause:
			_, err := asn1gen.Unmarshal(ie.Value, &cause)
			if err != nil {
				fmt.Println("error in unmarshaling cause in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal(ie.Value, &criticalityDiagnostics)
			if err != nil {
				fmt.Println("error in unmarshaling criticality diagnostics in error indication received from gNB : ", err)
				return
			}
		}
	}

	//todo : process decoded values

	fmt.Println("DECODED ERROR INDICATION DATA :")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "RAN UE NGAP ID\t = ", ranUeNgapId)
	fmt.Fprintln(w, "AMF UE NGAP ID\t = ", amfUeNgapId)
	fmt.Fprintf(w, "Cause\t = %+v\n", cause)
	fmt.Fprintf(w, "Criticality Diagnostics\t = %+v\n", criticalityDiagnostics)
	w.Flush()

}

func SendErrorIndicationFromAmf(ran *context.AmfRan, others ...interface{}) {
	fmt.Println()
	fmt.Println("---------- ERROR INDICATION FROM AMF -----------")

	var errorInd asn1gen.ErrorIndication = asn1gen.ErrorIndication{}
	if len(others) != 0 {
		fmt.Println("There are other parameters in error indication! Please encode them as well.")
	}
	encodedErrorInd, err := asn1gen.Marshal(errorInd)
	if err != nil {
		fmt.Println("error while marshaling error indication :", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "Error Indication\t = %+v\n", encodedErrorInd)
	w.Flush()

	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = 1
	ngapMsg.U.InitiatingMessage = &asn1gen.InitiatingMessage{
		ProcedureCode: asn1gen.Asn1vIdErrorIndication,
		Criticality:   asn1gen.CriticalityIgnore,
		Value:         encodedErrorInd,
	}
	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		fmt.Println("error while marshaling error indication :", err)
		return
	}

	/*
		if ran == nil {
			return errors.New("RAN is nil")
		}

		if len(encodedNgapMsg) == 0 {
			return errors.New("packet length is 0")
		}

		if ran.Conn == nil {
			return errors.New("RAN address is nil")
		}

		n, err := ran.Conn.Write(encodedNgapMsg)
		if err != nil {
			err := "Write error : " + err.Error()
			return errors.New(err)
		} else {
			println("Wrote ", n, " bytes")
		}
	*/
	fmt.Println()
	fmt.Println("ENCODED ERROR INDICATION : ", encodedNgapMsg)
}
