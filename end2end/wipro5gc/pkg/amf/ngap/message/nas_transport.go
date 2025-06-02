package message

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
	"time"

	"k8s.io/klog"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"w5gc.io/wipro5gcore/pkg/amf/context"
	"w5gc.io/wipro5gcore/pkg/amf/metrics"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/db/redis"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/generator"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc/protos/ngapNas/pb"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/ngapConvert"

	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"

	// "github.com/free5gc/openapi/models"

	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

type RanUe struct {
	RanUeNgapId int64
	AmfUeNgapId int64
}

func HandleUplinkNASTransport(ran *context.AmfRan, message *asn1gen.NGAPPDU, grpc *grpc.Grpc, client *redis.RedisClient) {
	klog.Infof("INSIDE UPLINK NAS TRANSPORT\n")

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
		klog.Infoln("RAN is nil")
		return
	}

	if message == nil {
		klog.Infoln("NGAP Message is nil")
		return
	}

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Infoln("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var uplinkNasTransport *asn1gen.UplinkNASTransport = &asn1gen.UplinkNASTransport{}
	_, err := asn1gen.Unmarshal(x, uplinkNasTransport)
	if err != nil {
		klog.Error("unmarshalling failed for initial ue message type : ", err)
		return
	}
	uplinkNasTransportIEs := uplinkNasTransport.ProtocolIEs
	for i := 0; i < len(uplinkNasTransportIEs); i++ {
		ie := uplinkNasTransportIEs[i]
		// klog.Infof("ie : %+v", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			amfuengapid_val := []byte(ie.Value)
			if amfuengapid_val == nil {
				klog.Infoln("AMF UE NGAP ID is nil")
				return
			}
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				klog.Infoln("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId > math.MaxUint32 {
				klog.Infoln("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdNASPDU:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), nasPdu)
			if err != nil {
				klog.Infoln("unmarshalling of NAS PDU failed : ", err)
				return
			}
			if *nasPdu == nil {
				klog.Infoln("NAS PDU is nil")
				return
			}

		case asn1gen.Asn1vIdUserLocationInformation:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), userLocationInformation)
			if err != nil {
				klog.Infoln("unmarshalling of User Location Information failed : ", err)
				return
			}
			if (*userLocationInformation == asn1gen.UserLocationInformation{}) {
				klog.Infoln("User Location Information is nil")
				return
			}
		}
		klog.Infof("USER LOCATION INFORMATION : %+v", userLocationInformation)
		klog.Infof("USER LOCATION INFORMATION NR : %+v", userLocationInformation.U.UserLocationInformationNR)
		klog.Infof("AMF UE NGAP ID : %+v", amfUeNgapId)
		klog.Infof("RAN UE NGAP ID : %+v", ranUeNgapId)
	}

	timestamp := binary.BigEndian.Uint32(*(userLocationInformation.U.UserLocationInformationNR.TimeStamp))
	ueLocationTimestamp := time.Unix(int64(timestamp), 0)

	amfUe, err := client.Read("amfUeNgapIdToSupi:" + strconv.Itoa(int(amfUeNgapId)))
	if err != nil {
		klog.Error("error in reading supi from db : ", err)
		return
	}
	klog.Info("amfUe struct : ", amfUe.AmfUeNgapId, amfUe.Supi)
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
		Supi:          amfUe.Supi,
	}

	// serializeMsg, _ := proto.Marshal(&createSmContext)
	serializeMsg, err := anypb.New(&createSmContext)
	if err != nil {
		klog.Errorln("Failed to marshal create sm context data : ", err)
		return
	}
	req := &pb.DataRequest{
		Data:    serializeMsg,
		ReqType: "Uplink NAS",
	}
	res := (*grpc).SendData(req)
	if res.Error != "" {
		klog.Errorln(res.Error)
	}
	// else {
	//TO DO
	/*if err := SendPduSessionResourceSetupRequest(ran, res.NasPdu); err != nil {
		klog.Errorln("Error in sending Pdu Session Resource Setup Request : ", err)
	}*/
	// }
	// klog.Infoln("AFTER SENDING REQUEST TYPE IN CHANNEL")
	/*establishmentAccept := nas.PduSessionEstablishmentAccept{
		ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES",
		PDUsessionId:                  "VAL_1",
		PTI:                           1,
		MessageType:                   "PDU_SESSION_ESTABLISHMENT_ACCEPT",
		PduSessionType:                "IPV4",
		SSCmode:                       "SSC_MODE_1",
		SessionAmbr:                   "MULT_1Kbps",
	}
	klog.Infof("Dummy Data for Establishment Accept: %+v\n", establishmentAccept)
	establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
	if err != nil {
		klog.Infoln("error while encoding PDU Session Establishment Accept: ", err)
		return
	}*/
	/*establishmentAcceptByteArray, err := os.ReadFile("/home/vboxuser/asn1c-v772-for-vm/golang/sample_per/ngap/src/nas/testFiles/MinPDUSessionEstablishmentAccept")
	if err != nil {
		klog.Infoln("error in reading pdu session establishment accept file : ", err)
		return
	}
	*/
	/*establishmentAcceptByteArray := []byte{46, 1, 1, 194, 17, 0, 9, 1, 0, 6, 49, 49, 1, 1, 255, 9, 6, 1, 3, 232, 1, 3, 232, 41, 5, 1, 10, 62, 0, 1, 34, 4, 1, 17, 34, 51, 121, 0, 6, 9, 32, 65, 1, 1, 9, 123, 0, 27, 128, 0, 13, 4, 8, 8, 8, 8, 0, 3, 16, 32, 1, 72, 96, 72, 96, 0, 0, 0, 0, 0, 0, 0, 0, 136, 136, 37, 10, 9, 105, 110, 116, 101, 114, 110, 101, 116, 50}
	klog.Infoln("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
	dlNAS := nas.DLNasModel{
		ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES",
		SecurityHeaderType:            "Plain 5GS NAS message",
		MessageType:                   "DL_NAS_TRANSPORT",
		PayLoadContainerType:          "N1 SM information",
	}
	dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
	dlNasByteArray[3] = 1
	if err != nil {
		klog.Infoln("Error while encoding DL NAS Transport:", err)
		return
	}
	klog.Infoln("Encoded DL NAS Transport:", dlNasByteArray)
	establishmentAcceptByteArrayLen := len(establishmentAcceptByteArray)
	intermediateArray := []byte{0, byte(establishmentAcceptByteArrayLen)}
	dlNasByteArray = append(dlNasByteArray, intermediateArray...)
	nasPduByteArray1 := append(dlNasByteArray, establishmentAcceptByteArray...)
	nasPdu1 := asn1gen.NASPDU(nasPduByteArray1)
	klog.Infoln("nas message : dl nas + pdu session est. accept = ", nasPdu1)
	if err := SendPduSessionResourceSetupRequest(ran, nasPdu1); err != nil {
		klog.Infoln("Error in sending Pdu Session Resource Setup Request : ", err)
		return
	}*/
}

// for registration request
func HandleInitialUEMessage(ran *context.AmfRan, message *asn1gen.NGAPPDU, grpc *grpc.Grpc, client *redis.RedisClient) {
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

	// get the values from the message

	metrics.RegistrationAttempts.Inc()

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
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("RAN is nil")
		return
	}

	if message == nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var initialUeMessage *asn1gen.InitialUEMessage = &asn1gen.InitialUEMessage{}
	_, err := asn1gen.Unmarshal(x, initialUeMessage)
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("unmarshalling failed for initial ue message type: ", err)
		return
	}
	initialUeMessageIEs := initialUeMessage.ProtocolIEs
	for i := 0; i < len(initialUeMessageIEs); i++ {
		ie := initialUeMessageIEs[i]
		klog.Infof("IE : %+v\n", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				klog.Error("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId > math.MaxUint32 {
				klog.Error("RAN UE NGAP ID has an invalid value")
				return
			}

		case asn1gen.Asn1vIdNASPDU:
			klog.Infoln("[NGAP] Decode NAS PDU")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), nasPdu)
			if err != nil {
				klog.Error("unmarshalling of NAS PDU failed : ", err)
				return
			}
			if *nasPdu == nil {
				klog.Error("NAS PDU is nil")
				return
			}

		case asn1gen.Asn1vIdUserLocationInformation:
			klog.Infoln("[NGAP] Decode User Location Information")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), userLocationInformation)
			if err != nil {
				klog.Error("unmarshalling of User Location Information failed : ", err)
				return
			}
			if (*userLocationInformation == asn1gen.UserLocationInformation{}) {
				klog.Error("User Location Information is nil")
				return
			}

		case asn1gen.Asn1vIdRRCEstablishmentCause:
			klog.Infoln("[NGAP] Decode RRC Establishment Cause")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &rrcEstablishmentCause)
			if err != nil {
				klog.Error("unmarshalling of RRC Establishment Cause failed : ", err)
				return
			}
			if rrcEstablishmentCause > 11 {
				klog.Error("RRC Establishment Cause is nil")
				return
			}

		case asn1gen.Asn1vIdFiveGSTMSI:
			klog.Infoln("[NGAP] Decode 5G TMSI")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), fiveGSTMSI)
			if err != nil {
				klog.Info("unmarshalling of 5GSTMSI failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(fiveGSTMSI, asn1gen.FiveGSTMSI{})) {
				klog.Info("5GSTMSI is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdAMFSetID:
			klog.Infoln("[NGAP] Decode AMF Set ID")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), amfSetId)
			if err != nil {
				klog.Info("unmarshalling of AMF Set ID failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(amfSetId, asn1gen.AMFSetID{})) {
				klog.Info("AMF Set ID is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdUEContextRequest:
			klog.Infoln("[NGAP] Decode UE Context Request")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ueContextRequest)
			if err != nil {
				klog.Info("unmarshalling of UE Context Request failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			//UeContextRequest is an enum whose value is 0,1
			if ueContextRequest != 0 && ueContextRequest != 1 {
				klog.Info("UE Context Request has an invalid value")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdAllowedNSSAI:
			klog.Infoln("[NGAP] Decode Allowed NSSAI")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), allowedNssai)
			if err != nil {
				klog.Info("unmarshalling of Allowed NSSAI failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if len(allowedNssai) == 0 {
				klog.Info("No Allowed NSSAI present")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdSourceToTargetAMFInformationReroute:
			klog.Infoln("[NGAP] Decode SourceToTargetAMFInformationReroute")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), sourceToTargetAMFInformationReroute)
			if err != nil {
				klog.Info("unmarshalling of Source to Target AMF Information Reroute failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if (reflect.DeepEqual(sourceToTargetAMFInformationReroute, asn1gen.SourceToTargetAMFInformationReroute{})) {
				klog.Info("Source to Target AMF Information Reroute is nil")
				//commented return as this IE presence is optional
				//return
			}

		case asn1gen.Asn1vIdSelectedPLMNIdentity:
			klog.Infoln("[NGAP] Decode Selected PLMN Identity")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), selectedPlmnIdentity)
			if err != nil {
				klog.Info("unmarshalling of Selected PLMN Identity failed : ", err)
				//commented return as this IE presence is optional
				//return
			}
			if *selectedPlmnIdentity == nil {
				klog.Info("Selected PLMN Identity is nil")
				//commented return as this IE presence is optional
				//return
			}
		}
	}

	//RAN-UE-NGAP-ID
	//to do : ask if a ranUe struct should be created?

	// send NAS PDU to NAS-MM module
	nasMsg := pb.NasMessage{
		NasPdu: *nasPdu,
	}
	// serializeNasMsg, _ := proto.Marshal(&nasMsg)
	serializeNasMsg, err := anypb.New(&nasMsg)
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Errorln("Error in marshaling nas message : ", err)
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
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error(res.Error)
		return
	}

	// store SUPI and AMF UE NGAP ID in DB
	amfUeNgapId := generator.GenerateAmfUeNgapId()
	if amfUeNgapId == uint64(math.Pow(2, 40)) {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("Max value of AMF UE NGAP ID reached")
		return
	}
	supi := res.Supi
	amfUeRan := redis.AmfUeRan{
		Supi:        supi,
		AmfUeNgapId: amfUeNgapId,
		RanUeNgapId: uint64(ranUeNgapId),
	}
	amfUeRanJson, err := json.Marshal(amfUeRan)
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("failed to marshal amfUe struct : ", err)
	}
	instance, err := client.Read("supiToAmfUeNgapId:" + supi)
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Info("supi is not yet registered: ", err)
	}
	if instance != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("instance of this supi already exists!", instance)
		return
	}
	_, err = client.Write("amfUeNgapIdToSupi:"+strconv.Itoa(int(amfUeNgapId)), string(amfUeRanJson))
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("failed to write to amfUe DB : ", err)
	}
	_, err = client.Write("supiToAmfUeNgapId:"+supi, string(amfUeRanJson))
	if err != nil {
		metrics.RegistrationFailures.WithLabelValues("ERROR").Inc()
		klog.Error("failed to write to amfUe DB : ", err)
	}

	// send Downlink NAS Transport or UE Context Setup Request to RAN
	klog.Info("UE Context Request value : ", ueContextRequest)
	//ueransim
	if ueContextRequest == 0 {
		SendInitialContextSetupRequest(ran, ranUeNgapId, amfUeNgapId, res.NasPdu)
	} else {
		SendDownlinkNASTransport(ran, ranUeNgapId, amfUeNgapId, res.NasPdu)
	}
	//n3iwf
}

func SendInitialContextSetupRequest(ran *context.AmfRan, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, nasMsg asn1gen.NASPDU) {
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

	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	klog.Info("Initial Context Setup Request : amfUeNgapId : ", amfUeNgapId)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling amf ue ngap id in initial context setup request : ", err)
		return
	}

	klog.Info("Initial Context Setup Request : ranUeNgapId : ", ranUeNgapId)
	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling ran ue ngap id in initial context setup request : ", err)
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
		klog.Infoln("error in marshaling guami in initial context setup request : ", err)
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
		klog.Infoln("error in marshaling allowed nssai in initial context setup request : ", err)
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
		klog.Infoln("error in marshaling ue security capabilities in initial context setup request : ", err)
		return
	}

	securityKey := asn1gen.SecurityKey(ngapConvert.HexToBitString("0000000000000000000000000000000000000000000000000000000000000000", 256))
	encodedSecurityKey, err := asn1gen.Marshal(securityKey)

	encodedNasMsg, err := asn1gen.Marshal(nasMsg)
	if err != nil {
		klog.Infoln("error in marshaling nas msg in initial context setup request : ", err)
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
		klog.Infoln("error in marshaling initial context setup request : ", err)
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
		klog.Infoln("error in marshaling ngap msg : ", err)
		return
	}

	if ran == nil {
		klog.Infoln(errors.New("RAN is nil"))
		return
	}

	if len(encodedNgapMsg) == 0 {
		klog.Infoln(errors.New("packet length is 0"))
		return
	}

	if ran.Conn == nil {
		klog.Infoln(errors.New("RAN address is nil"))
		return
	}

	n, err := ran.Conn.Write(encodedNgapMsg)
	if err != nil {
		err := "Write error : " + err.Error()
		klog.Infoln(errors.New(err))
		return
	} else {
		klog.Infoln("[INITIAL CONTEXT SETUP] Wrote ", n, " bytes")
	}
}

func HandleInitialContextSetupResponse(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	if ran == nil {
		klog.Infoln("RAN is nil")
		return
	}

	if message == nil {
		klog.Infoln("NGAP Message is nil")
		return
	}
	var ranUeNgapId asn1gen.RANUENGAPID
	var amfUeNgapId asn1gen.AMFUENGAPID
	criticalityDiagnostics := asn1gen.CriticalityDiagnostics{}
	pduSessionResourceSetupListSuRes := asn1gen.PDUSessionResourceSetupListSURes{}
	pduSessionResourceFailedToSetupListSuRes := asn1gen.PDUSessionResourceFailedToSetupListSURes{}

	successfulOutcome := message.U.SuccessfulOutcome
	if successfulOutcome == nil {
		klog.Infoln("Successful Outcome message is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	var initialContextSetupResponse *asn1gen.InitialContextSetupResponse = &asn1gen.InitialContextSetupResponse{}
	_, err := asn1gen.Unmarshal(x, initialContextSetupResponse)
	if err != nil {
		klog.Error("unmarshalling failed for initial context setup response type : ", err)
		return
	}
	initialContextSetupResponseIEs := initialContextSetupResponse.ProtocolIEs
	for i := 0; i < len(initialContextSetupResponseIEs); i++ {
		ie := initialContextSetupResponseIEs[i]
		if reflect.ValueOf(ie).IsZero() {
			klog.Infof("ie is null\n")
			return
		} else {
			klog.Infof("ie : %+v", ie)
		}
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			} else {
				klog.Infoln("Initial Context Setup Response : AMF UE NGAP ID : ", amfUeNgapId)
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId < 0 || amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				klog.Infoln("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			} else {
				klog.Infoln("Initial Context Setup Response : RAN UE NGAP ID : ", ranUeNgapId)
			}
			if ranUeNgapId < 0 || ranUeNgapId > math.MaxUint32 {
				klog.Infoln("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceSetupListSuRes)
			if err != nil {
				klog.Error("Error in unmarshaling pduSessionResourceSetupListSuRes : ", err)
			}
			if len(pduSessionResourceSetupListSuRes) == 0 {
				klog.Error("pduSessionResourceSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceFailedToSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceFailedToSetupListSuRes)
			if err != nil {
				klog.Error("Error in unmarshaling pduSessionResourceFailedToSetupListSuRes : ", err)
			}
			if len(pduSessionResourceFailedToSetupListSuRes) == 0 {
				klog.Error("pduSessionResourceFailedToSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &criticalityDiagnostics)
			if err != nil {
				klog.Error("Error in unmarshaling criticalityDiagnostics")
			}

		}
	}
	//todo: process the decoded values

}

func HandleInitialContextSetupFailure(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	if ran == nil {
		klog.Infoln("RAN is nil")
		return
	}

	if message == nil {
		klog.Infoln("NGAP Message is nil")
		return
	}
	var ranUeNgapId asn1gen.RANUENGAPID
	var amfUeNgapId asn1gen.AMFUENGAPID
	cause := asn1gen.Cause{}
	criticalityDiagnostics := asn1gen.CriticalityDiagnostics{}
	pduSessionResourceFailedToSetupListSuRes := asn1gen.PDUSessionResourceFailedToSetupListSURes{}
	unsuccessfulOutcome := message.U.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		klog.Infoln("Unsuccessful Outcome message is nil")
		return
	}
	x := []byte(unsuccessfulOutcome.Value)
	var initialContextSetupFailure *asn1gen.InitialContextSetupFailure = &asn1gen.InitialContextSetupFailure{}
	_, err := asn1gen.Unmarshal(x, initialContextSetupFailure)
	if err != nil {
		klog.Error("unmarshalling failed for initial context setup failure type : ", err)
		return
	}
	initialContextSetupFailureIEs := initialContextSetupFailure.ProtocolIEs
	for i := 0; i < len(initialContextSetupFailureIEs); i++ {
		ie := initialContextSetupFailureIEs[i]
		if reflect.ValueOf(ie).IsZero() {
			klog.Infof("ie is null\n")
			return
		} else {
			klog.Infof("ie : %+v", ie)
		}
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of AMF UE NGAP ID failed : ", err)
				return
			}
			var x uint64 = uint64(math.Pow(2, 40) - 1)
			if amfUeNgapId < 0 || amfUeNgapId > asn1gen.AMFUENGAPID(x) {
				klog.Infoln("value of AMF UE NGAP ID is out of range")
				return
			}
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &ranUeNgapId)
			if err != nil {
				klog.Infoln("unmarshalling of RAN UE NGAP ID failed : ", err)
				return
			}
			if ranUeNgapId < 0 || ranUeNgapId > math.MaxUint32 {
				klog.Infoln("value of RAN UE NGAP ID is out of range")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceFailedToSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceFailedToSetupListSuRes)
			if err != nil {
				klog.Error("Error in unmarshaling pduSessionResourceFailedToSetupListSuRes : ", err)
			}
			if len(pduSessionResourceFailedToSetupListSuRes) == 0 {
				klog.Error("pduSessionResourceFailedToSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &criticalityDiagnostics)
			if err != nil {
				klog.Error("Error in unmarshaling criticalityDiagnostics")
			}

		case asn1gen.Asn1vIdCause:
			klog.Infoln("[NGAP] Decode IE Cause")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), cause)
			if err != nil {
				klog.Error("Error in unmarshaling Cause : ", err)
				return
			}
			if cause == (asn1gen.Cause{}) {
				klog.Error("Cause is nil")
				return
			}
		}
	}
	//todo: process the decoded values
}

func SendDownlinkNASTransport(ran *context.AmfRan, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, nasMsg asn1gen.NASPDU) error {
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
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		return err
	}

	klog.Infoln("ranUeNgapId:", ranUeNgapId)

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
		klog.Infoln("[DOWNLINK NAS TRANSPORT] Wrote ", n, " bytes")
	}
	metrics.RegistrationSuccess.Inc()
	return nil
}

func HandleErrorIndicationFromRan(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	klog.Infoln()
	klog.Infoln("---------- ERROR INDICATION FROM RAN -----------")
	if ran == nil {
		klog.Infoln("RAN is nil")
		return
	}

	if message == nil {
		klog.Infoln("NGAP Message is nil")
		return
	}

	// var ranUeNgapId *asn1gen.RANUENGAPID
	ranUeNgapId := new(asn1gen.RANUENGAPID)
	amfUeNgapId := new(asn1gen.AMFUENGAPID)
	cause := &asn1gen.Cause{}
	criticalityDiagnostics := &asn1gen.CriticalityDiagnostics{}

	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Infoln("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	errorInd := asn1gen.ErrorIndication{}
	_, err := asn1gen.Unmarshal(x, &errorInd)
	if err != nil {
		klog.Infoln("error in unmarshaling error indication received from gNB : ", err)
		return
	}

	for i := 0; i < len(errorInd.ProtocolIEs); i++ {
		ie := errorInd.ProtocolIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdRANUENGAPID:
			_, err := asn1gen.Unmarshal(ie.Value, ranUeNgapId)
			if err != nil {
				klog.Infoln("error in unmarshaling ran ue ngap id in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal(ie.Value, amfUeNgapId)
			if err != nil {
				klog.Infoln("error in unmarshaling amf ue ngap id in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdCause:
			_, err := asn1gen.Unmarshal(ie.Value, cause)
			if err != nil {
				klog.Infoln("error in unmarshaling cause in error indication received from gNB : ", err)
				return
			}
		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal(ie.Value, criticalityDiagnostics)
			if err != nil {
				klog.Infoln("error in unmarshaling criticality diagnostics in error indication received from gNB : ", err)
				return
			}
		}
	}

	//todo : process decoded values

	klog.Infoln("DECODED ERROR INDICATION DATA :")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	klog.Infoln(w, "RAN UE NGAP ID\t = ", *ranUeNgapId)
	fmt.Println("RAN UE NGAP ID = ", *ranUeNgapId)
	klog.Infoln(w, "AMF UE NGAP ID\t = ", *amfUeNgapId)
	fmt.Println("AMF UE NGAP ID = ", *amfUeNgapId)
	klog.Infof("Cause\t = %+v\n", *cause)
	fmt.Println("Cause = ", *cause)
	klog.Infof("Criticality Diagnostics\t = %+v\n", *criticalityDiagnostics)
	fmt.Println("Criticality Diagnostics = ", *criticalityDiagnostics)
	w.Flush()
}

func SendErrorIndicationFromAmf(ran *context.AmfRan, others ...interface{}) {
	klog.Infoln()
	klog.Infoln("---------- ERROR INDICATION FROM AMF -----------")

	var errorInd asn1gen.ErrorIndication = asn1gen.ErrorIndication{}
	if len(others) != 0 {
		klog.Infoln("There are other parameters in error indication! Please encode them as well.")
	}
	encodedErrorInd, err := asn1gen.Marshal(errorInd)
	if err != nil {
		klog.Infoln("error while marshaling error indication :", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	klog.Infoln(w, "Error Indication\t = %+v\n", encodedErrorInd)
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
		klog.Infoln("error while marshaling error indication :", err)
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
	klog.Infoln("ENCODED ERROR INDICATION : ", encodedNgapMsg)
}
