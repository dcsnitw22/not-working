package context

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/ngapConvert"
)

type SupportedTAI struct {
	Tai        openapi_commn_client.Tai
	SNssaiList []openapi_commn_client.Snssai
}

func InTaiList(servedTai openapi_commn_client.Tai, taiList []openapi_commn_client.Tai) bool {
	for _, tai := range taiList {
		/*if reflect.DeepEqual(tai, servedTai) {
			return true
		}*/
		// fmt.Println(servedTai.PlmnId.Mcc, "==", tai.PlmnId.Mcc, "&&", servedTai.PlmnId.Mnc, "==", tai.PlmnId.Mnc, "&&", servedTai.Tac, "==", tai.Tac)
		if servedTai.PlmnId.Mcc == tai.PlmnId.Mcc && servedTai.PlmnId.Mnc == tai.PlmnId.Mnc && servedTai.Tac == tai.Tac {

			return true
		}
	}
	return false
}

type AmfRan struct {
	// no such field in asn1gen code
	//RanPresent int
	RanId  *openapi_commn_client.GlobalRanNodeId
	Name   string
	AnType openapi_commn_client.AccessType
	/* socket Connect*/
	Conn net.Conn
	/* Supported TA List */
	SupportedTAList []SupportedTAI

	/* RAN UE List */
	RanUeList sync.Map // RanUeNgapId as key

	/* logger */
	Log *logrus.Entry
}

func (ran *AmfRan) SetRanId(ranNodeId *asn1gen.GlobalRANNodeID) {
	ranId := ngapConvert.RanIdToModels(*ranNodeId)
	// ran.RanPresent = ranNodeId.Present
	ran.RanId = &ranId
	ran.AnType = openapi_commn_client.AccessType("3GPP_ACCESS")
	// if ranNodeId.Present == ngapType.GlobalRANNodeIDPresentGlobalN3IWFID {
	// 	ran.AnType = models.AccessType_NON_3_GPP_ACCESS
	// } else {
	// 	ran.AnType = models.AccessType__3_GPP_ACCESS
	// }
}
