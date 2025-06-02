package message

import (
	"encoding/hex"
	"errors"

	"k8s.io/klog"

	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
	"w5gc.io/wipro5gcore/pkg/amf/context"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/ngapConvert"

	// change to openapi package path when VMs are up
	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

const (
	MaxNumOfSlice int = 1024
)

// taken as constant for now
var SupportTaiLists = []openapi_commn_client.Tai{
	{PlmnId: openapi_commn_client.PlmnId{Mcc: "286", Mnc: "01"}, Tac: "000001"},
	{PlmnId: openapi_commn_client.PlmnId{Mcc: "904", Mnc: "217"}, Tac: "92718E"},
}

func NewSupportedTAI() (tai context.SupportedTAI) {
	tai.SNssaiList = make([]openapi_commn_client.Snssai, 0, MaxNumOfSlice)
	return
}

func SendNGSetupResponse(ran *context.AmfRan) error {
	//AMF Name
	var amfName asn1gen.AMFName = asn1gen.AMFName(ran.Name)
	encodedAmfName, err := asn1gen.Marshal(amfName)
	if err != nil {
		return err
	}

	//Served GUAMI List
	backupAmf := asn1gen.AMFName("testBackupAmf")
	plmnid := openapi_commn_client.PlmnId{
		Mcc: "286",
		Mnc: "01",
	}
	guami := openapi_commn_client.Guami{
		PlmnId: plmnid,
		AmfId:  "218A9E",
	}
	regionId, setId, prtId := ngapConvert.AmfIdToNgap(guami.AmfId)
	var servedGuamiList asn1gen.ServedGUAMIList = asn1gen.ServedGUAMIList{
		asn1gen.ServedGUAMIItem{
			GUAMI: asn1gen.GUAMI{
				PLMNIdentity: ngapConvert.PlmnIdToNgap(plmnid),
				AMFRegionID:  asn1gen.AMFRegionID(regionId),
				AMFSetID:     asn1gen.AMFSetID(setId),
				AMFPointer:   asn1gen.AMFPointer(prtId),
			},
			BackupAMFName: &backupAmf,
		},
	}

	encodedServedGuamiList, err := asn1gen.Marshal(servedGuamiList)
	if err != nil {
		return err
	}

	//Relative AMF Capacity
	//hardcoded random value
	var relAmfCap asn1gen.RelativeAMFCapacity = 10
	encodedRelAmfCap, err := asn1gen.Marshal(relAmfCap)
	if err != nil {
		return err
	}

	//PLMN Support List
	sd := asn1gen.SD(asn1rt.OctetString(ngapConvert.HexToBitString("000001", 24).Bytes))
	// sd := asn1gen.SD("001")
	convertedPlmnId := ngapConvert.PlmnIdToNgap(plmnid)
	// fmt.Info("Converted PLMN ID to NGAP : ", convertedPlmnId)
	var plmnSupportList asn1gen.PLMNSupportList = asn1gen.PLMNSupportList{
		asn1gen.PLMNSupportItem{
			PLMNIdentity: convertedPlmnId,
			SliceSupportList: asn1gen.SliceSupportList{
				asn1gen.SliceSupportItem{
					SNSSAI: asn1gen.SNSSAI{
						SST: asn1gen.SST("1"),
						SD:  &sd,
					},
				},
			},
		},
	}

	encodedPlmnSupportList, err := asn1gen.Marshal(plmnSupportList)
	if err != nil {
		e := "plmnsupportlist error : " + err.Error()
		return errors.New(e)
	}

	var ngSetupResponse asn1gen.NGSetupResponse = asn1gen.NGSetupResponse{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdAMFName,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedAmfName,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdServedGUAMIList,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedServedGuamiList,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdRelativeAMFCapacity,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedRelAmfCap,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdPLMNSupportList,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedPlmnSupportList,
			},
		},
	}
	encodedNgSetupResponse, err := asn1gen.Marshal(ngSetupResponse)
	if err != nil {
		return err
	}

	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = asn1gen.NGAPPDUSuccessfulOutcomeTAG
	ngapMsg.U.SuccessfulOutcome = &asn1gen.SuccessfulOutcome{
		ProcedureCode: asn1gen.Asn1vIdNGSetup,
		Criticality:   asn1gen.CriticalityReject,
		Value:         encodedNgSetupResponse,
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
		klog.Info("Wrote ", n, " bytes")
	}
	return nil
}

func SendNGSetupFailure(ran *context.AmfRan, cause asn1gen.Cause) error {
	encodedCause, err := asn1gen.Marshal(cause)
	if err != nil {
		return err
	}
	var ngSetupFailure asn1gen.NGSetupFailure = asn1gen.NGSetupFailure{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdCause,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedCause,
			},
		},
	}

	encodedNgSetupFailure, err := asn1gen.Marshal(ngSetupFailure)
	if err != nil {
		return err
	}

	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = asn1gen.NGAPPDUUnsuccessfulOutcomeTAG
	ngapMsg.U.UnsuccessfulOutcome = &asn1gen.UnsuccessfulOutcome{
		ProcedureCode: asn1gen.Asn1vIdNGSetup,
		Criticality:   asn1gen.CriticalityReject,
		Value:         encodedNgSetupFailure,
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
		klog.Info("Wrote ", n, " bytes")
	}
	return nil

}

func HandleNGSetupRequest(ran *context.AmfRan, message *asn1gen.NGAPPDU) {

	var globalRANNodeID *asn1gen.GlobalRANNodeID = &asn1gen.GlobalRANNodeID{}
	var rANNodeName asn1gen.RANNodeName = ""
	var supportedTAList *asn1gen.SupportedTAList = &asn1gen.SupportedTAList{}
	var pagingDRX asn1gen.PagingDRX = 5
	var cause asn1gen.Cause

	if ran == nil {
		klog.Error("ran is nil")
		return
	}
	if message == nil {
		klog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Error("Initiating Message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var nGSetupRequest *asn1gen.NGSetupRequest = &asn1gen.NGSetupRequest{}
	_, err := asn1gen.Unmarshal(x, nGSetupRequest)
	if err != nil {
		klog.Info("unmarshalling failed for ng setup request type")
		return
	}
	// if nGSetupRequest == asn1gen.NGSetupRequest{} {
	// 	klog.Info("NGSetupRequest is nil")
	// 	return
	// }

	for i := 0; i < len(nGSetupRequest.ProtocolIEs); i++ {
		ie := nGSetupRequest.ProtocolIEs[i]
		klog.Infof("IE : %+v\n", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdGlobalRANNodeID:
			klog.Info("[NGAP] Decode IE GlobalRANNodeID")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), globalRANNodeID)
			// fmt.Info("globalRANNodeID:", globalRANNodeID)
			if err != nil {
				klog.Error("Error in unmarshaling global RAN Node ID : ", err)
				return
			}
			if *globalRANNodeID == (asn1gen.GlobalRANNodeID{}) {
				klog.Error("GlobalRANNodeID is nil")
				return
			}
		case asn1gen.Asn1vIdSupportedTAList:
			klog.Info("[NGAP] Decode Supported TA List")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), supportedTAList)
			if err != nil {
				klog.Error("Error in unmarshaling Supported TA List : ", err)
			}
			if len(*supportedTAList) == 0 {
				klog.Error("SupportedTAList is nil")
				return
			}
		case asn1gen.Asn1vIdRANNodeName:
			klog.Info("[NGAP] Decode RAN Node Name")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &rANNodeName)
			if err != nil {
				klog.Error("Error in unmarshaling RAN Node Name : ", err)
			}
			if rANNodeName == "" {
				klog.Error("RANNodeName is nil")
				return
			}
		case asn1gen.Asn1vIdDefaultPagingDRX:
			klog.Info("[NGAP] Decode IE DefaultPagingDRX")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), &pagingDRX)
			if err != nil {
				klog.Error("Error in unmarshaling RAN Node Name : ", err)
			}
			if pagingDRX == 5 {
				klog.Error("DefaultPagingDRX is nil")
				return
			}
		}
	}

	ran.SetRanId(globalRANNodeID)
	if rANNodeName != "" {
		ran.Name = string(rANNodeName)
	}
	if pagingDRX != 5 {
		klog.Infof("PagingDRX[%d]", pagingDRX)
	}

	for i := 0; i < len(*supportedTAList); i++ {
		supportedTAItem := (*supportedTAList)[i]
		tac := hex.EncodeToString(supportedTAItem.TAC)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList); j++ {
			supportedTAI := NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			//fmt.Infof("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)

			} else {
				break
			}
		}
	}

	//fmt.Infof("RAN Struct : %+v", ran)

	if len(ran.SupportedTAList) == 0 {
		klog.Error("NG-Setup failure: No supported TA exist in NG-Setup request")
		cause.T = asn1gen.CauseMiscTAG
		var cmu asn1gen.CauseMisc = asn1gen.CauseMiscUnspecified
		cause.U.Misc = &cmu
	} else {
		var found bool
		for _, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, SupportTaiLists) {
				// Info("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			klog.Error("NG-Setup failure: Cannot find Served TAI in AMF")
			cause.T = asn1gen.CauseMiscTAG
			var cmuplmn asn1gen.CauseMisc = asn1gen.CauseMiscUnknownPLMN
			cause.U.Misc = &cmuplmn
		}
	}

	//fmt.Info("Cause : ", cause)
	// Info("-- End of NG Setup Request. Send NG Setup Response/Failure --")

	// if cause.Present == ngapType.CausePresentNothing {
	// 	ngap_message.SendNGSetupResponse(ran)
	// } else {
	// 	ngap_message.SendNGSetupFailure(ran, cause)
	// }

	if cause == (asn1gen.Cause{}) {
		err := SendNGSetupResponse(ran)
		if err != nil {
			klog.Error("Error in sending ng setup response : ", err)
		}
	} else {
		klog.Error("NG setup failed. failure cause : ", cause)
		err := SendNGSetupFailure(ran, cause)
		if err != nil {
			klog.Error("error in sending NG setup failure message : ", err)
		}
	}
}

// NG Reset initiated by AMF
func SendNGResetMessage(ran *context.AmfRan) error {
	//Cause
	var unknownCause asn1gen.CauseMisc = 5
	var cause asn1gen.Cause
	cause.T = 5
	cause.U.Misc = &unknownCause
	encodedCause, err := asn1gen.Marshal(cause)
	if err != nil {
		return err
	}

	//Choice
	//reset all
	/*var resetAll asn1gen.ResetAll = 0
	var resetType asn1gen.ResetType
	resetType.T = 1
	resetType.U.NGInterface = &resetAll
	encodedResetType, err := asn1gen.Marshal(resetType)
	if err != nil {
		return err
	}*/

	//reset part
	var amfUeNgapId asn1gen.AMFUENGAPID = 10
	var ranUeNgapId asn1gen.RANUENGAPID = 1
	resetPart := asn1gen.UEAssociatedLogicalNGConnectionList{
		asn1gen.UEAssociatedLogicalNGConnectionItem{
			AMFUENGAPID: &amfUeNgapId,
			RANUENGAPID: &ranUeNgapId,
		},
	}
	var resetType asn1gen.ResetType
	resetType.T = 2
	resetType.U.PartOfNGInterface = &resetPart
	/*resetType := asn1gen.ResetType{
		T: 2,
		U: struct {
			NGInterface       *asn1gen.ResetAll
			PartOfNGInterface *asn1gen.UEAssociatedLogicalNGConnectionList
			ChoiceExtensions  *asn1gen.ProtocolIESingleContainer
		}{
			NGInterface:       nil,
			PartOfNGInterface: &resetPart,
			ChoiceExtensions:  nil,
		},
	}*/
	encodedResetType, err := asn1gen.Marshal(resetType)
	if err != nil {
		return err
	}

	//NGReset
	var ngReset asn1gen.NGReset = asn1gen.NGReset{
		ProtocolIEs: asn1gen.ProtocolIEContainer{
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdCause,
				Criticality: asn1gen.CriticalityIgnore,
				Value:       encodedCause,
			},
			asn1gen.ProtocolIEField{
				Id:          asn1gen.Asn1vIdResetType,
				Criticality: asn1gen.CriticalityReject,
				Value:       encodedResetType,
			},
		},
	}
	encodedNgReset, err := asn1gen.Marshal(ngReset)
	if err != nil {
		return err
	}

	//NGAPPDU
	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = 0
	ngapMsg.U.InitiatingMessage = &asn1gen.InitiatingMessage{
		ProcedureCode: asn1gen.Asn1vIdNGReset,
		Criticality:   asn1gen.CriticalityReject,
		Value:         encodedNgReset,
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
		klog.Info("Wrote ", n, " bytes")
	}
	return nil
}

func HandleNGResetAcknowledge(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	var ueAssociatedLogicalNGConnectionList asn1gen.UEAssociatedLogicalNGConnectionList = asn1gen.UEAssociatedLogicalNGConnectionList{}
	var criticalityDiagnostics asn1gen.CriticalityDiagnostics = asn1gen.CriticalityDiagnostics{}
	if ran == nil {
		klog.Error("ran is nil")
		return
	}
	if message == nil {
		klog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.U.SuccessfulOutcome
	if successfulOutcome == nil {
		klog.Error("Successful Outcome response is nil")
		return
	}
	x := []byte(successfulOutcome.Value)
	var nGResetAck *asn1gen.NGResetAcknowledge = &asn1gen.NGResetAcknowledge{}
	_, err := asn1gen.Unmarshal(x, nGResetAck)
	if err != nil {
		klog.Info("unmarshalling failed for ng reset acknowledge type")
		return
	}

	for i := 0; i < len(nGResetAck.ProtocolIEs); i++ {
		ie := nGResetAck.ProtocolIEs[i]
		klog.Infof("IE : %+v\n", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdUEAssociatedLogicalNGConnectionList:
			klog.Info("[NGAP] Decode IE UEAssociatedLogicalNGConnectionList")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), ueAssociatedLogicalNGConnectionList)
			if err != nil {
				klog.Error("Error in unmarshaling UE Associated Logical NG Connection List : ", err)
				return
			}
			if len(ueAssociatedLogicalNGConnectionList) == 0 {
				klog.Error("UEAssociatedLogicalNGConnectionList is empty")
				return
			}

		case asn1gen.Asn1vIdCriticalityDiagnostics:
			klog.Info("[NGAP] Decode IE Criticality Diagnostics")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), criticalityDiagnostics)
			if err != nil {
				klog.Error("Error in unmarshaling Criticality Diagnostics : ", err)
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
	//handle ue associated logical ng connection list and criticality diagnostics

}

// NG Reset initiated by NG RAN
func HandleNGResetMessage(ran *context.AmfRan, message *asn1gen.NGAPPDU) {
	var cause asn1gen.Cause = asn1gen.Cause{}
	var resetType asn1gen.ResetType = asn1gen.ResetType{}
	if ran == nil {
		klog.Error("ran is nil")
		return
	}
	if message == nil {
		klog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.U.InitiatingMessage
	if initiatingMessage == nil {
		klog.Error("Initiating Message is nil")
		return
	}
	x := []byte(initiatingMessage.Value)
	var nGReset *asn1gen.NGReset = &asn1gen.NGReset{}
	_, err := asn1gen.Unmarshal(x, nGReset)
	if err != nil {
		klog.Info("unmarshalling failed for ng reset type")
		return
	}

	for i := 0; i < len(nGReset.ProtocolIEs); i++ {
		ie := nGReset.ProtocolIEs[i]
		klog.Infof("IE : %+v\n", ie)
		switch ie.Id {
		case asn1gen.Asn1vIdCause:
			klog.Info("[NGAP] Decode IE Cause")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), cause)
			if err != nil {
				klog.Error("Error in unmarshaling Cause : ", err)
				return
			}
			if cause == (asn1gen.Cause{}) {
				klog.Error("Cause is nil")
				return
			}

		case asn1gen.Asn1vIdResetType:
			klog.Info("[NGAP] Decode IE Reset Type")
			_, err := asn1gen.Unmarshal([]byte(ie.Value), resetType)
			if err != nil {
				klog.Error("Error in unmarshaling reset type : ", err)
				return
			}
			if resetType.T != 1 && resetType.T != 0 {
				klog.Error("Reset Type is invalid")
				return
			}
		}
	}
	//handle cause and reset type
	err = SendNGResetAcknowledge(ran)
	if err != nil {
		klog.Error("Error in NG Reset Ack : ", err)
		return
	}
}

func SendNGResetAcknowledge(ran *context.AmfRan) error {
	var ngResetAck asn1gen.NGResetAcknowledge = asn1gen.NGResetAcknowledge{}
	encodedNgResetAck, err := asn1gen.Marshal(ngResetAck)
	if err != nil {
		return err
	}

	var ngapMsg asn1gen.NGAPPDU = asn1gen.NGAPPDU{}
	ngapMsg.T = 2
	ngapMsg.U.SuccessfulOutcome = &asn1gen.SuccessfulOutcome{
		ProcedureCode: asn1gen.Asn1vIdNGReset,
		Criticality:   asn1gen.CriticalityReject,
		Value:         encodedNgResetAck,
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
		klog.Info("Wrote ", n, " bytes")
	}
	return nil
}
