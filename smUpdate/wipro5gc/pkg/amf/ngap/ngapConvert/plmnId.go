package ngapConvert

import (
	"encoding/hex"
	"fmt"
	"strings"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen"

	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

func PlmnIdToModels(ngapPlmnId asn1gen.PLMNIdentity) (modelsPlmnid openapi_commn_client.PlmnId) {
	value := ngapPlmnId
	hexString := strings.Split(hex.EncodeToString(value), "")
	modelsPlmnid.Mcc = hexString[1] + hexString[0] + hexString[3]
	if hexString[2] == "f" {
		modelsPlmnid.Mnc = hexString[5] + hexString[4]
	} else {
		modelsPlmnid.Mnc = hexString[2] + hexString[5] + hexString[4]
	}
	return
}

func PlmnIdToNgap(modelsPlmnid openapi_commn_client.PlmnId) asn1gen.PLMNIdentity {
	var hexString string
	mcc := strings.Split(modelsPlmnid.Mcc, "")
	mnc := strings.Split(modelsPlmnid.Mnc, "")
	if len(modelsPlmnid.Mnc) == 2 {
		hexString = mcc[1] + mcc[0] + "f" + mcc[2] + mnc[1] + mnc[0]
	} else {
		hexString = mcc[1] + mcc[0] + mnc[0] + mcc[2] + mnc[2] + mnc[1]
	}

	var ngapPlmnId asn1gen.PLMNIdentity
	if plmnId, err := hex.DecodeString(hexString); err != nil {
		fmt.Printf("Decode plmn failed: %+v", err)
	} else {
		ngapPlmnId = plmnId
	}
	return ngapPlmnId
}
