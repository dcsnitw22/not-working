package ngapConvert

import (
	"encoding/hex"
	"strconv"

	"w5gc.io/wipro5gcore/asn1gen"

	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

func SNssaiToModels(ngapSnssai asn1gen.SNSSAI) (modelsSnssai openapi_commn_client.Snssai) {
	i, e := strconv.ParseInt(hex.EncodeToString(ngapSnssai.SST), 10, 32)
	if e == nil {
		modelsSnssai.Sst = int32(i)
	}
	if ngapSnssai.SD != nil {
		sd := hex.EncodeToString(*(ngapSnssai.SD))
		modelsSnssai.Sd = &sd
	}
	return
}
