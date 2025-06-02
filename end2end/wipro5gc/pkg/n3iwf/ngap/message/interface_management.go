package message

import (
	"errors"
	"net"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
)

// TO DO : Fill in the necessary IEs for the NGSetupRequest message
func buildNGSetupRequest() asn1gen.NGAPPDU {
	plmnId := []byte{0x82, 0xf6, 0x10}
	globalRANNodeID := asn1gen.GlobalRANNodeID{
		T: 3,
		U: struct {
			GlobalGNBID      *asn1gen.GlobalGNBID
			GlobalNgENBID    *asn1gen.GlobalNgENBID
			GlobalN3IWFID    *asn1gen.GlobalN3IWFID
			ChoiceExtensions *asn1gen.ProtocolIESingleContainer
		}{
			GlobalN3IWFID: &asn1gen.GlobalN3IWFID{
				PLMNIdentity: asn1gen.PLMNIdentity(plmnId),
				N3IWFID: asn1gen.N3IWFID{
					T: 1,
					U: struct {
						N3IWFID          *asn1rt.BitString
						ChoiceExtensions *asn1gen.ProtocolIESingleContainer
					}{
						N3IWFID: &asn1rt.BitString{
							BitLength: 16,
							Bytes:     []byte{0x01, 0x01, 0x01, 0x01},
						},
					},
				},
			},
		},
	}
	encodedGlobalRANNodeID, err := asn1gen.Marshal(globalRANNodeID)
	if err != nil {
		klog.Errorf("failed to encode GlobalRANNodeID: %v", err)
		return asn1gen.NGAPPDU{}
	}

	rANNodeName := asn1gen.RANNodeName("n3iwf")
	encodedRANNodeName, err := asn1gen.Marshal(rANNodeName)
	if err != nil {
		klog.Errorf("failed to encode RAN node name: %v", err)
		return asn1gen.NGAPPDU{}
	}

	sd := asn1gen.SD("001")
	tac := []byte{0x00, 0x00, 0x01}
	supportedTAList := asn1gen.SupportedTAList{
		asn1gen.SupportedTAItem{
			TAC: asn1gen.TAC(tac),
			BroadcastPLMNList: asn1gen.BroadcastPLMNList{
				asn1gen.BroadcastPLMNItem{
					PLMNIdentity: asn1gen.PLMNIdentity(plmnId),
					TAISliceSupportList: asn1gen.SliceSupportList{
						asn1gen.SliceSupportItem{
							SNSSAI: asn1gen.SNSSAI{
								SST: asn1gen.SST("1"),
								SD:  &sd,
							},
						},
					},
				},
			},
		},
	}
	encodedSupportedTAList, err := asn1gen.Marshal(supportedTAList)
	if err != nil {
		klog.Errorf("failed to encode supported TA list: %v", err)
		return asn1gen.NGAPPDU{}
	}

	defualtPagingDRX := asn1gen.PagingDRX(0)
	encodedDefaultPagingDRX, err := asn1gen.Marshal(defualtPagingDRX)
	if err != nil {
		klog.Errorf("failed to encode default paging DRX: %v", err)
		return asn1gen.NGAPPDU{}
	}

	// Build the NGSetupRequest message
	ngSetupRequest := asn1gen.NGSetupRequest{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdAMFName,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedRANNodeName,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdGlobalRANNodeID,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedGlobalRANNodeID,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdSupportedTAList,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedSupportedTAList,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdDefaultPagingDRX,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedDefaultPagingDRX,
			},
		},
	}

	encoodedNgSetupRequest, err := asn1gen.Marshal(ngSetupRequest)
	if err != nil {
		klog.Errorf("failed to encode NGSetupRequest: %v", err)
		return asn1gen.NGAPPDU{}
	}
	// Wrap the NGSetupRequest in an NGAPPDU
	ngapPDU := asn1gen.NGAPPDU{}
	ngapPDU.T = asn1gen.NGAPPDUInitiatingMessageTAG
	ngapPDU.U.InitiatingMessage = &asn1gen.InitiatingMessage{
		ProcedureCode: asn1gen.Asn1vIdNGSetup,
		Criticality:   asn1gen.CriticalityReject,
		Value:         encoodedNgSetupRequest,
	}

	return ngapPDU
}

func SendN2SetupRequest(conn net.Conn) error {
	n2SetupRequest := buildNGSetupRequest()
	if n2SetupRequest == (asn1gen.NGAPPDU{}) {
		return errors.New("Failed to build N2 Setup Request : NGAPPDU is empty")
	}
	encodedRequest, err := asn1gen.Marshal(n2SetupRequest)
	if err != nil {
		klog.Errorf("failed to encode N2 Setup Request: %v", err)
		return err
	}
	_, err = conn.Write(encodedRequest)
	if err != nil {
		klog.Errorf("failed to write N2 Setup Request to connection: %v", err)
		return err
	}
	klog.Infoln("N2 Setup Request sent successfully")
	return nil
}

func HandleNGSetupResponse(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
	successfulOutcome := ngapPDU.U.SuccessfulOutcome
	if successfulOutcome == nil {
		klog.Error("Initiating Message is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	ngSetupResponse := asn1gen.NGSetupResponse{}
	_, err := asn1gen.Unmarshal(x, &ngSetupResponse)
	if err != nil {
		klog.Errorf("failed to decode NGSetupResponse: %v", err)
		return
	}
	klog.Infof("Received NG Setup Response: %v", ngSetupResponse)
	for i := 0; i < len(ngSetupResponse.ProtocolIEs); i++ {
		ie := ngSetupResponse.ProtocolIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdAMFName:
			amfName := asn1gen.AMFName("")
			_, err := asn1gen.Unmarshal(ie.Value, &amfName)
			if err != nil {
				klog.Errorf("failed to decode AMFName: %v", err)
				return
			}
			klog.Infof("AMF Name: %v", amfName)
		case asn1gen.Asn1vIdServedGUAMIList:
			servedGUAMIList := asn1gen.ServedGUAMIList{}
			_, err := asn1gen.Unmarshal(ie.Value, &servedGUAMIList)
			if err != nil {
				klog.Errorf("failed to decode ServedGUAMIList: %v", err)
				return
			}
			klog.Infof("Served GUAMI List: %v", servedGUAMIList)
		case asn1gen.Asn1vIdRelativeAMFCapacity:
			relativeAMFCapacity := asn1gen.RelativeAMFCapacity(0)
			_, err := asn1gen.Unmarshal(ie.Value, &relativeAMFCapacity)
			if err != nil {
				klog.Errorf("failed to decode RelativeAMFCapacity: %v", err)
				return
			}
			klog.Infof("Relative AMF Capacity: %v", relativeAMFCapacity)
		case asn1gen.Asn1vIdPLMNSupportList:
			plmnSupportList := asn1gen.PLMNSupportList{}
			_, err := asn1gen.Unmarshal(ie.Value, &plmnSupportList)
			if err != nil {
				klog.Errorf("failed to decode PLMNSupportList: %v", err)
				return
			}
			klog.Infof("PLMN Support List: %v", plmnSupportList)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}

}

func HandleNGSetupFailure(conn net.Conn, ngapPDU *asn1gen.NGAPPDU) {
	unsuccessfulOutcome := ngapPDU.U.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		klog.Error("Unsuccessful Outcome is nil")
		return
	}
	x := []byte(unsuccessfulOutcome.Value)
	ngSetupFailure := asn1gen.NGSetupFailure{}
	_, err := asn1gen.Unmarshal(x, &ngSetupFailure)
	if err != nil {
		klog.Errorf("failed to decode NGSetupFailure: %v", err)
		return
	}
	klog.Infof("Received NG Setup Failure: %v", ngSetupFailure)
	for i := 0; i < len(ngSetupFailure.ProtocolIEs); i++ {
		ie := ngSetupFailure.ProtocolIEs[i]
		switch ie.Id {
		case asn1gen.Asn1vIdCause:
			cause := asn1gen.Cause{}
			_, err := asn1gen.Unmarshal(ie.Value, &cause)
			if err != nil {
				klog.Errorf("failed to decode Cause: %v", err)
				return
			}
			klog.Infof("Cause: %v", cause)
		default:
			klog.Infof("Unknown IE: %v", ie.Id)
		}
	}
}
