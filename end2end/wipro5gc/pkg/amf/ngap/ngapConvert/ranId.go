package ngapConvert

import (
	"w5gc.io/wipro5gcore/asn1gen"

	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

func RanIdToModels(ranNodeId asn1gen.GlobalRANNodeID) (ranId openapi_commn_client.GlobalRanNodeId) {
	choice := ranNodeId.T
	switch choice {
	case asn1gen.GlobalRANNodeIDGlobalGNBIDTAG:
		//output
		ranId.GNbId = new(openapi_commn_client.GNbId)
		gnbId := ranId.GNbId

		//input
		ngapGnbId := ranNodeId.U.GlobalGNBID
		plmnid := PlmnIdToModels(ngapGnbId.PLMNIdentity)

		//output
		ranId.PlmnId = &plmnid

		// if ngapGnbId.GNBID.Present == ngapType.GNBIDPresentGNBID {
		if ngapGnbId.GNBID.T == asn1gen.GNBIDGNBIDTAG {
			choiceGnbId := ngapGnbId.GNBID.U.GNBID
			gnbId.BitLength = int32(choiceGnbId.BitLength)
			gnbId.GNBValue = BitStringToHex(choiceGnbId)
		}
		// case ngapType.GlobalRANNodeIDPresentGlobalNgENBID:
		// 	ngapNgENBID := ranNodeId.GlobalNgENBID
		// 	plmnid := PlmnIdToModels(ngapNgENBID.PLMNIdentity)
		// 	ranId.PlmnId = &plmnid
		// 	if ngapNgENBID.NgENBID.Present == ngapType.NgENBIDPresentMacroNgENBID {
		// 		macroNgENBID := ngapNgENBID.NgENBID.MacroNgENBID
		// 		ranId.NgeNbId = "MacroNGeNB-" + BitStringToHex(macroNgENBID)
		// 	} else if ngapNgENBID.NgENBID.Present == ngapType.NgENBIDPresentShortMacroNgENBID {
		// 		shortMacroNgENBID := ngapNgENBID.NgENBID.ShortMacroNgENBID
		// 		ranId.NgeNbId = "SMacroNGeNB-" + BitStringToHex(shortMacroNgENBID)
		// 	} else if ngapNgENBID.NgENBID.Present == ngapType.NgENBIDPresentLongMacroNgENBID {
		// 		longMacroNgENBID := ngapNgENBID.NgENBID.LongMacroNgENBID
		// 		ranId.NgeNbId = "LMacroNGeNB-" + BitStringToHex(longMacroNgENBID)
		// 	}
		// case ngapType.GlobalRANNodeIDPresentGlobalN3IWFID:
		// 	ngapN3IWFID := ranNodeId.GlobalN3IWFID
		// 	plmnid := PlmnIdToModels(ngapN3IWFID.PLMNIdentity)
		// 	ranId.PlmnId = &plmnid
		// 	if ngapN3IWFID.N3IWFID.Present == ngapType.N3IWFIDPresentN3IWFID {
		// 		choiceN3IWFID := ngapN3IWFID.N3IWFID.N3IWFID
		// 		ranId.N3IwfId = BitStringToHex(choiceN3IWFID)
		// 	}
	}

	return ranId
}
