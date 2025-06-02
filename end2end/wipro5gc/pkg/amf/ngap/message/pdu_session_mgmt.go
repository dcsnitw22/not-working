package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/context"
	"w5gc.io/wipro5gcore/pkg/amf/metrics"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/db/redis"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/ngapConvert"

	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
)

// var N1Message = make(chan ([]byte), 1)
// var N2Info = make(chan ([]byte), 1)

// TO DO : should be called after AMF has received n1n2messagetransfer. process n1 and n2 information.
func SendPduSessionResourceSetupRequest(ran *context.AmfRan, n1 []byte, n2 []byte, supi string, client *redis.RedisClient) error {
	// n1 := <-grpcserver.N1Message
	// n2 := <-grpcserver.N2Info

	metrics.CreateAttempts.Inc()

	klog.Info("Inside send pdu session resource setup request function")
	klog.Info("n1 byte data : ", n1, " n2 byte data : ", n2)
	n2d := &asn1gen.PDUSessionResourceSetupRequestTransfer{}
	err := json.Unmarshal(n2, n2d)
	if err != nil {
		return err
	}
	klog.Infof("n2 data: +%v", n2d)

	//AMF UE NGAP ID
	//DONE :fetch this AMFUENGAPID from redis based on the supi received; hardcoded for now
	amfUeRan, err := client.Read("supiToAmfUeNgapId:" + supi)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		klog.Error("could not read from redis : ", err)
	}
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeRan.AmfUeNgapId)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	//RAN UE NGAP ID
	ranUeNgapId := asn1gen.RANUENGAPID(amfUeRan.RanUeNgapId)
	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	//NAS PDU
	// nasPduFile, err = os.Open("../testdata/MinPDUSessionEstablishmentRequest")
	// hard coded MinPDUSessionEstablishmentAccept will not work, because the Request file coming from UERANSIM UE has different data. Accept file should be according to that file.
	//TODO: replace with PDUSessionEstablishmentAccept received from SMF
	/*nasPduBin, err := os.ReadFile("../testdata/MinPDUSessionEstablishmentAccept")
	if err != nil {
		return err
	}
	nasPdu := asn1gen.NASPDU(nasPduBin)
	encodedNasPdu, err := asn1gen.Marshal(nasPdu)*/
	nasMsg := asn1gen.NASPDU(n1)
	encodedNasPdu, err := asn1gen.Marshal(nasMsg)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	//TODO: replace with values received from SMF in N2 SM Information

	// UL NG-U UP TNL Information
	/*ipstr := "127.0.0.1"
	ipaddr, err := net.ResolveIPAddr("ip", ipstr)
	if err != nil {
		return err
	}
	ip := ipaddr.IP
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
				GTPTEID: asn1gen.GTPTEID("1231"),
			},
			ChoiceExtensions: nil,
		},
	}
	encodedUlNguUpTnlInformation, err := asn1gen.Marshal(ulNguUpTnlInformation)
	if err != nil {
		return err
	}*/

	//PDU Session Type
	/*pduSessionType := asn1gen.PDUSessionType(0)
	encodedPduSessionType, err := asn1gen.Marshal(pduSessionType)
	if err != nil {
		return err
	}*/

	//QoS Flow Setup Request List
	/*qosFlowSetupRequestList := asn1gen.QosFlowSetupRequestList{
		asn1gen.QosFlowSetupRequestItem{
			QosFlowIdentifier: asn1gen.QosFlowIdentifier(1),
			QosFlowLevelQosParameters: asn1gen.QosFlowLevelQosParameters{
				QosCharacteristics: asn1gen.QosCharacteristics{
					T: 1,
					U: struct {
						NonDynamic5QI    *asn1gen.NonDynamic5QIDescriptor
						Dynamic5QI       *asn1gen.Dynamic5QIDescriptor
						ChoiceExtensions *asn1gen.ProtocolIESingleContainer
					}{
						NonDynamic5QI: &asn1gen.NonDynamic5QIDescriptor{
							FiveQI: asn1gen.FiveQI(0),
						},
						Dynamic5QI:       nil,
						ChoiceExtensions: nil,
					},
				},
				AllocationAndRetentionPriority: asn1gen.AllocationAndRetentionPriority{
					PriorityLevelARP:        asn1gen.PriorityLevelARP(1),
					PreEmptionCapability:    asn1gen.PreEmptionCapability(0),
					PreEmptionVulnerability: asn1gen.PreEmptionVulnerability(0),
				},
				GBRQosInformation: &asn1gen.GBRQosInformation{},
			},
		},
	}
	encodedQosFlowSetupRequestList, err := asn1gen.Marshal(qosFlowSetupRequestList)
	if err != nil {
		return err
	}*/

	//PDU Session Resource Setup Request List
	sd := asn1gen.SD(asn1rt.OctetString(ngapConvert.HexToBitString("000001", 24).Bytes))
	pduSessionResourceSetupListSuReq := asn1gen.PDUSessionResourceSetupListSUReq{
		asn1gen.PDUSessionResourceSetupItemSUReq{
			PDUSessionID: asn1gen.PDUSessionID(1),
			SNSSAI: asn1gen.SNSSAI{
				SST: asn1gen.SST{1},
				SD:  &sd,
			},
			/*PDUSessionResourceSetupRequestTransfer: asn1gen.PDUSessionResourceSetupRequestTransfer{
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
			},*/
			PDUSessionResourceSetupRequestTransfer: *n2d,
		},
	}
	encodedPduSessionResourceSetupListSuReq, err := asn1gen.Marshal(pduSessionResourceSetupListSuReq)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	pduSessionResourceSetupRequest := asn1gen.PDUSessionResourceSetupRequest{
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
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdPDUSessionResourceSetupListSUReq,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedPduSessionResourceSetupListSuReq,
			},
		},
	}
	encodedPduSessionResourceSetupRequest, err := asn1gen.Marshal(pduSessionResourceSetupRequest)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: 1,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			InitiatingMessage: &asn1gen.InitiatingMessage{
				ProcedureCode: asn1gen.Asn1vIdPDUSessionResourceSetup,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedPduSessionResourceSetupRequest,
			},
			SuccessfulOutcome:   nil,
			UnsuccessfulOutcome: nil,
			ExtElem1:            nil,
		},
	}
	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return err
	}

	if ran == nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return errors.New("RAN is nil")
	}

	if len(encodedNgapMsg) == 0 {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return errors.New("packet length is 0")
	}

	if ran.Conn == nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		return errors.New("RAN address is nil")
	}

	n, err := ran.Conn.Write(encodedNgapMsg)
	klog.Info("Sent encoded Ngap message to RAN")
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		err := "Write error : " + err.Error()
		return errors.New(err)
	} else {
		println("Wrote ", n, " bytes")
	}
	return nil
}

func HandlePduSessionResourceSetupResponse(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	if ran == nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		fmt.Println("RAN is nil")
		return
	}

	if message == nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		fmt.Println("NGAP Message is nil")
		return
	}

	var amfUeNgapId asn1gen.AMFUENGAPID
	var ranUeNgapId asn1gen.RANUENGAPID
	pduSessionResourceSetupListSuRes := asn1gen.PDUSessionResourceSetupListSURes{}
	pduSessionResourceFailedToSetupListSuRes := asn1gen.PDUSessionResourceFailedToSetupListSURes{}
	successfulOutcome := message.U.SuccessfulOutcome
	if successfulOutcome == nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		fmt.Println("Successful Outcome message is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	var pduSessionResourceSetupResponse *asn1gen.PDUSessionResourceSetupResponse = &asn1gen.PDUSessionResourceSetupResponse{}
	_, err := asn1gen.Unmarshal(x, pduSessionResourceSetupResponse)
	if err != nil {
		metrics.CreateFailures.WithLabelValues("ERROR").Inc()
		println("unmarshalling failed for pdu session resource setup response type")
		return
	}
	pduSessionResourceSetupResponseIEs := pduSessionResourceSetupResponse.ProtocolIEs
	for i := 0; i < len(pduSessionResourceSetupResponseIEs); i++ {
		ie := pduSessionResourceSetupResponseIEs[i]
		if reflect.ValueOf(ie).IsZero() {
			metrics.CreateFailures.WithLabelValues("ERROR").Inc()
			fmt.Printf("ie is null\n")
			return
		} else {
			fmt.Printf("ie : %+v", ie)
		}
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &amfUeNgapId)
			if err != nil {
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
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
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
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
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
				klog.Error("Error in unmarshaling pduSessionResourceSetupListSuRes : ", err)
			}
			if len(pduSessionResourceSetupListSuRes) == 0 {
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
				klog.Error("pduSessionResourceSetupListSuRes is nil")
				return
			}

		case asn1gen.Asn1vIdPDUSessionResourceFailedToSetupListSURes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceFailedToSetupListSuRes)
			if err != nil {
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
				klog.Error("Error in unmarshaling pduSessionResourceFailedToSetupListSuRes : ", err)
			}
			if len(pduSessionResourceFailedToSetupListSuRes) == 0 {
				metrics.CreateFailures.WithLabelValues("ERROR").Inc()
				klog.Error("pduSessionResourceFailedToSetupListSuRes is nil")
				return
			}
		}
	}
	//todo: process the decoded values
	metrics.CreateSuccess.Inc()
}

func SendPduSessionResourceReleaseCommand(ran *context.AmfRan) error {
	//AMF UE NGAP ID
	amfUeNgapId := asn1gen.AMFUENGAPID(10)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		return err
	}

	//RAN UE NGAP ID
	ranUeNgapId := asn1gen.RANUENGAPID(1)
	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		return err
	}

	//RAN Paging Priority
	ranPagingPriority := asn1gen.RANPagingPriority(11)
	encodedRanPagingPriority, err := asn1gen.Marshal(ranPagingPriority)
	if err != nil {
		return err
	}

	//NAS PDU
	// nasPduFile, err = os.Open("../testdata/MinPDUSessionReleaseRequest")
	// hard coded MinPDUSessionReleaseAccept will not work, because the Request file coming from UERANSIM UE has different data. Accept file should be according to that file.
	// TODO: replace with sth
	nasPduBin, err := os.ReadFile("../testdata/MinPDUSessionReleaseAccept")
	if err != nil {
		return err
	}
	nasPdu := asn1gen.NASPDU(nasPduBin)
	encodedNasPdu, err := asn1gen.Marshal(nasPdu)
	if err != nil {
		return err
	}

	//TODO: replace with values received from SMF in N2 SM Information
	causeMisc := asn1gen.CauseMisc(6)
	cause := asn1gen.Cause{
		T: 5,
		U: struct {
			RadioNetwork     *asn1gen.CauseRadioNetwork
			Transport        *asn1gen.CauseTransport
			Nas              *asn1gen.CauseNas
			Protocol         *asn1gen.CauseProtocol
			Misc             *asn1gen.CauseMisc
			ChoiceExtensions *asn1gen.ProtocolIESingleContainer
		}{
			Misc: &causeMisc,
		},
	}
	//PDU Session Resource to Release List
	pduSessionResourceToReleaseListRelCmd := asn1gen.PDUSessionResourceToReleaseListRelCmd{
		asn1gen.PDUSessionResourceToReleaseItemRelCmd{
			PDUSessionID: asn1gen.PDUSessionID(1),
			PDUSessionResourceReleaseCommandTransfer: asn1gen.PDUSessionResourceReleaseCommandTransfer{
				Cause: cause,
			},
		},
	}
	encodedPduSessionResourceToReleaseListRelCmd, err := asn1gen.Marshal(pduSessionResourceToReleaseListRelCmd)
	if err != nil {
		return err
	}

	pduSessionResourceReleaseCommand := asn1gen.PDUSessionResourceReleaseCommand{
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
				Id:          asn1gen.Asn1vIdRANPagingPriority,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedRanPagingPriority,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdNASPDU,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedNasPdu,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdPDUSessionResourceToReleaseListRelCmd,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedPduSessionResourceToReleaseListRelCmd,
			},
		},
	}
	encodedPduSessionResourceReleaseCommand, err := asn1gen.Marshal(pduSessionResourceReleaseCommand)
	if err != nil {
		return err
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: 1,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			InitiatingMessage: &asn1gen.InitiatingMessage{
				ProcedureCode: asn1gen.Asn1vIdPDUSessionResourceRelease,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedPduSessionResourceReleaseCommand,
			},
			SuccessfulOutcome:   nil,
			UnsuccessfulOutcome: nil,
			ExtElem1:            nil,
		},
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
		println("Wrote ", n, " bytes")
	}
	return nil
}

func HandlePduSessionResourceReleaseResponse(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
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
	pduSessionResourceReleasedList := asn1gen.PDUSessionResourceReleasedListRelRes{}
	userLocationInformation := &asn1gen.UserLocationInformation{}

	successfulOutcome := message.U.SuccessfulOutcome
	if successfulOutcome == nil {
		fmt.Println("Successful Outcome message is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	var pduSessionResourceReleaseResponse *asn1gen.PDUSessionResourceReleaseResponse = &asn1gen.PDUSessionResourceReleaseResponse{}
	_, err := asn1gen.Unmarshal(x, pduSessionResourceReleaseResponse)
	if err != nil {
		println("unmarshalling failed for pdu session resource setup response type")
		return
	}
	pduSessionResourceReleaseResponseIEs := pduSessionResourceReleaseResponse.ProtocolIEs
	for i := 0; i < len(pduSessionResourceReleaseResponseIEs); i++ {
		ie := pduSessionResourceReleaseResponseIEs[i]
		fmt.Printf("ie : %+v", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
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

		case asn1gen.Asn1vIdPDUSessionResourceReleasedListRelRes:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pduSessionResourceReleasedList)
			if err != nil {
				klog.Error("unmarshaling of PDU Session Resource Released List failed")
				return
			}
			if len(pduSessionResourceReleasedList) == 0 {
				klog.Error("pduSessionResourceReleasedList is empty")
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
		case asn1gen.Asn1vIdCriticalityDiagnostics:
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &criticalityDiagnostics)
			if err != nil {
				fmt.Println("unmarshalling of Criticality Diagnostics failed : ", err)
				return
			}
			if len(criticalityDiagnostics.IEsCriticalityDiagnostics) == 0 {
				klog.Error("criticalityDiagnostics list is nil")
				return
			}
			if criticalityDiagnostics.ProcedureCode == nil {
				klog.Error("criticalityDiagnostics : procedure code is nil")
				return
			}
			if criticalityDiagnostics.TriggeringMessage == nil {
				klog.Error("criticalityDiagnostics : triggering message is nil")
				return
			}
			if criticalityDiagnostics.ProcedureCriticality == nil {
				klog.Error("criticalityDiagnostics : procedure criticality is nil")
				return
			}
			if *(criticalityDiagnostics.ProcedureCode) > 255 {
				klog.Error("criticalityDiagnostics : procedure code value is invalid")
				return
			}
			if *(criticalityDiagnostics.TriggeringMessage) > 2 {
				klog.Error("criticalityDiagnostics : triggering message value is invalid")
				return
			}
			if *(criticalityDiagnostics.ProcedureCriticality) > 2 {
				klog.Error("criticalityDiagnostics : procedure criticality value is invalid")
				return
			}
		}
	}
	//todo : process the decoded values
}
