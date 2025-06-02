package message

import (
	"errors"
	"net"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
)

func HandlePduSessionResourceSetupRequest(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
	if conn == nil {
		klog.Infoln("Connection is nil")
		return
	}

	if ngapPDU == nil {
		klog.Infoln("NGAP PDU is nil")
		return
	}

	initiatingMessage := ngapPDU.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Infoln("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var pduSessionResourceSetupRequest *asn1gen.PDUSessionResourceSetupRequest = &asn1gen.PDUSessionResourceSetupRequest{}
	_, err := asn1gen.Unmarshal(x, pduSessionResourceSetupRequest)
	if err != nil {
		klog.Error("unmarshalling failed for PDU session resource setup request type : ", err)
		return
	}
	pduSessionResourceSetupRequestIEs := pduSessionResourceSetupRequest.ProtocolIEs
	for i := 0; i < len(pduSessionResourceSetupRequestIEs); i++ {
		ie := pduSessionResourceSetupRequestIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			klog.Infof("AMF UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdRANUENGAPID:
			klog.Infof("RAN UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdNASPDU:
			klog.Infof("NAS PDU: %v", ie.Value)
		case asn1gen.Asn1vIdPDUSessionResourceSetupListSUReq:
			klog.Infof("PDU Session Resource Setup List SU Req: %v", ie.Value)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}
	// Process the decoded values and send a response if necessary
}

func HandlePduSessionResourceReleaseCommand(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
	if conn == nil {
		klog.Infoln("Connection is nil")
		return
	}

	if ngapPDU == nil {
		klog.Infoln("NGAP PDU is nil")
		return
	}

	initiatingMessage := ngapPDU.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Infoln("Initiating message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var pduSessionResourceReleaseCommand *asn1gen.PDUSessionResourceReleaseCommand = &asn1gen.PDUSessionResourceReleaseCommand{}
	_, err := asn1gen.Unmarshal(x, pduSessionResourceReleaseCommand)
	if err != nil {
		klog.Error("unmarshalling failed for PDU session resource release command type : ", err)
		return
	}
	pduSessionResourceReleaseCommandIEs := pduSessionResourceReleaseCommand.ProtocolIEs
	for i := 0; i < len(pduSessionResourceReleaseCommandIEs); i++ {
		ie := pduSessionResourceReleaseCommandIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			klog.Infof("AMF UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdRANUENGAPID:
			klog.Infof("RAN UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdNASPDU:
			klog.Infof("NAS PDU: %v", ie.Value)
		case asn1gen.Asn1vIdPDUSessionResourceToReleaseListRelCmd:
			klog.Infof("PDU Session Resource To Release List Rel Cmd: %v", ie.Value)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}

	// Process the decoded values and send a response if necessary
}

func BuildPduSessionResourceSetupResponse(conn net.Conn, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, nasMsg asn1gen.NASPDU) {
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling amf ue ngap id in PDU session resource setup response : ", err)
		return
	}

	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling ran ue ngap id in PDU session resource setup response : ", err)
		return
	}

	encodedNasMsg, err := asn1gen.Marshal(nasMsg)
	if err != nil {
		klog.Infoln("error in marshaling nas msg in PDU session resource setup response : ", err)
		return
	}

	pduSessionResourceSetupResponse := asn1gen.PDUSessionResourceSetupResponse{
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
				Value:       encodedNasMsg,
			},
		},
	}
	encodedPduSessionResourceSetupResponse, err := asn1gen.Marshal(pduSessionResourceSetupResponse)
	if err != nil {
		klog.Infoln("error in marshaling PDU session resource setup response : ", err)
		return
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: asn1gen.NGAPPDUSuccessfulOutcomeTAG,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			SuccessfulOutcome: &asn1gen.SuccessfulOutcome{
				ProcedureCode: asn1gen.Asn1vIdPDUSessionResourceSetup,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedPduSessionResourceSetupResponse,
			},
		},
	}

	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		klog.Infoln("error in marshaling ngap msg : ", err)
		return
	}

	if conn == nil {
		klog.Infoln(errors.New("Connection is nil"))
		return
	}

	if len(encodedNgapMsg) == 0 {
		klog.Infoln(errors.New("packet length is 0"))
		return
	}

	n, err := conn.Write(encodedNgapMsg)
	if err != nil {
		err := "Write error : " + err.Error()
		klog.Infoln(errors.New(err))
		return
	} else {
		klog.Infoln("[PDU SESSION RESOURCE SETUP RESPONSE] Wrote ", n, " bytes")
	}
}

func BuildPduSessionResourceReleaseResponse(conn net.Conn, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, cause asn1gen.Cause) {
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling amf ue ngap id in PDU session resource release response : ", err)
		return
	}

	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling ran ue ngap id in PDU session resource release response : ", err)
		return
	}

	encodedCause, err := asn1gen.Marshal(cause)
	if err != nil {
		klog.Infoln("error in marshaling cause in PDU session resource release response : ", err)
		return
	}

	pduSessionResourceReleaseResponse := asn1gen.PDUSessionResourceReleaseResponse{
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
				Id:          asn1gen.Asn1vIdCause,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedCause,
			},
		},
	}
	encodedPduSessionResourceReleaseResponse, err := asn1gen.Marshal(pduSessionResourceReleaseResponse)
	if err != nil {
		klog.Infoln("error in marshaling PDU session resource release response : ", err)
		return
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: asn1gen.NGAPPDUSuccessfulOutcomeTAG,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			SuccessfulOutcome: &asn1gen.SuccessfulOutcome{
				ProcedureCode: asn1gen.Asn1vIdPDUSessionResourceRelease,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedPduSessionResourceReleaseResponse,
			},
		},
	}

	encodedNgapMsg, err := asn1gen.Marshal(ngapMsg)
	if err != nil {
		klog.Infoln("error in marshaling ngap msg : ", err)
		return
	}

	if conn == nil {
		klog.Infoln(errors.New("Connection is nil"))
		return
	}

	if len(encodedNgapMsg) == 0 {
		klog.Infoln(errors.New("packet length is 0"))
		return
	}

	n, err := conn.Write(encodedNgapMsg)
	if err != nil {
		err := "Write error : " + err.Error()
		klog.Infoln(errors.New(err))
		return
	} else {
		klog.Infoln("[PDU SESSION RESOURCE RELEASE RESPONSE] Wrote ", n, " bytes")
	}
}
