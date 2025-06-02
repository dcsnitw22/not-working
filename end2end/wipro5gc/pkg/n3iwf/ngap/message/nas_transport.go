package message

import (
	"errors"
	"net"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
)

func buildInitialUEMessage() {
	// Implementation for building Initial UE Message
	// Fill in the necessary IEs for the Initial UE Message
}

func buildUplinkNASTransport() {
	// Implementation for building Uplink NAS Transport
	// Fill in the necessary IEs for the Uplink NAS Transport
}

func buildInitialContextSetupResponse() {
	// Implementation for building Initial Context Setup Response
	// Fill in the necessary IEs for the Initial Context Setup Response
}

func buildInitialContextSetupFailure() {
	// Implementation for building Initial Context Setup Failure
	// Fill in the necessary IEs for the Initial Context Setup Failure
}

func HandleInitialContextSetupRequest(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
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
	var initialContextSetupRequest *asn1gen.InitialContextSetupRequest = &asn1gen.InitialContextSetupRequest{}
	_, err := asn1gen.Unmarshal(x, initialContextSetupRequest)
	if err != nil {
		klog.Error("unmarshalling failed for initial context setup request type : ", err)
		return
	}
	initialContextSetupRequestIEs := initialContextSetupRequest.ProtocolIEs
	for i := 0; i < len(initialContextSetupRequestIEs); i++ {
		ie := initialContextSetupRequestIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			klog.Infof("AMF UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdRANUENGAPID:
			klog.Infof("RAN UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdGUAMI:
			klog.Infof("GUAMI: %v", ie.Value)
		case asn1gen.Asn1vIdAllowedNSSAI:
			klog.Infof("Allowed NSSAI: %v", ie.Value)
		case asn1gen.Asn1vIdUESecurityCapabilities:
			klog.Infof("UE Security Capabilities: %v", ie.Value)
		case asn1gen.Asn1vIdSecurityKey:
			klog.Infof("Security Key: %v", ie.Value)
		case asn1gen.Asn1vIdNASPDU:
			klog.Infof("NAS PDU: %v", ie.Value)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}

	// Process the decoded values and send a response if necessary
}

func HandleDownlinkNASTransport(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
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
	var downlinkNasTransport *asn1gen.DownlinkNASTransport = &asn1gen.DownlinkNASTransport{}
	_, err := asn1gen.Unmarshal(x, downlinkNasTransport)
	if err != nil {
		klog.Error("unmarshalling failed for downlink NAS transport type : ", err)
		return
	}
	downlinkNasTransportIEs := downlinkNasTransport.ProtocolIEs
	for i := 0; i < len(downlinkNasTransportIEs); i++ {
		ie := downlinkNasTransportIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdAMFUENGAPID:
			klog.Infof("AMF UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdRANUENGAPID:
			klog.Infof("RAN UE NGAP ID: %v", ie.Value)
		case asn1gen.Asn1vIdNASPDU:
			klog.Infof("NAS PDU: %v", ie.Value)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}

	// Process the decoded values and send a response if necessary
}

func SendInitialContextSetupResponse(conn net.Conn, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, nasMsg asn1gen.NASPDU) {
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling amf ue ngap id in initial context setup response : ", err)
		return
	}

	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling ran ue ngap id in initial context setup response : ", err)
		return
	}

	encodedNasMsg, err := asn1gen.Marshal(nasMsg)
	if err != nil {
		klog.Infoln("error in marshaling nas msg in initial context setup response : ", err)
		return
	}

	initialContextSetupResponse := asn1gen.InitialContextSetupResponse{
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
	encodedInitialContextSetupResponse, err := asn1gen.Marshal(initialContextSetupResponse)
	if err != nil {
		klog.Infoln("error in marshaling initial context setup response : ", err)
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
				ProcedureCode: asn1gen.Asn1vIdInitialContextSetup,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedInitialContextSetupResponse,
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
		klog.Infoln("[INITIAL CONTEXT SETUP RESPONSE] Wrote ", n, " bytes")
	}
}

func SendInitialContextSetupFailure(conn net.Conn, ranUeNgapId asn1gen.RANUENGAPID, amfUeNgapIdVal uint64, cause asn1gen.Cause) {
	amfUeNgapId := asn1gen.AMFUENGAPID(amfUeNgapIdVal)
	encodedAmfUeNgapId, err := asn1gen.Marshal(amfUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling amf ue ngap id in initial context setup failure : ", err)
		return
	}

	encodedRanUeNgapId, err := asn1gen.Marshal(ranUeNgapId)
	if err != nil {
		klog.Infoln("error in marshaling ran ue ngap id in initial context setup failure : ", err)
		return
	}

	encodedCause, err := asn1gen.Marshal(cause)
	if err != nil {
		klog.Infoln("error in marshaling cause in initial context setup failure : ", err)
		return
	}

	initialContextSetupFailure := asn1gen.InitialContextSetupFailure{
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
	encodedInitialContextSetupFailure, err := asn1gen.Marshal(initialContextSetupFailure)
	if err != nil {
		klog.Infoln("error in marshaling initial context setup failure : ", err)
		return
	}

	ngapMsg := asn1gen.NGAPPDU{
		T: asn1gen.NGAPPDUUnsuccessfulOutcomeTAG,
		U: struct {
			InitiatingMessage   *asn1gen.InitiatingMessage
			SuccessfulOutcome   *asn1gen.SuccessfulOutcome
			UnsuccessfulOutcome *asn1gen.UnsuccessfulOutcome
			ExtElem1            *asn1rt.Asn1ChoiceExt
		}{
			UnsuccessfulOutcome: &asn1gen.UnsuccessfulOutcome{
				ProcedureCode: asn1gen.Asn1vIdInitialContextSetup,
				Criticality:   asn1gen.CriticalityReject,
				Value:         encodedInitialContextSetupFailure,
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
		klog.Infoln("[INITIAL CONTEXT SETUP FAILURE] Wrote ", n, " bytes")
	}
}
