package apiclient

import (
	"context"
	"encoding/json"
	"time"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
)

type ApiClient interface {
	Start()
	N1N2MessageTransfer(sourceIP string, ueCtxId string, n1N2MessageTransferReqData openapi_commn_client.N1N2MessageTransferReqData, binaryDataN1MessageContentFile []byte, binaryDataN2InfoContentFile []byte)
}

type ApiClientInfo struct {
	clientStartTime time.Time
	nodeInfo        config.SmfNodeInfo
	openApiClient   map[string]*openapi_commn_client.APIClient
}

// changed this paramter to type config.CscConfig instead of config.SmfNodeInfo
func NewApiClient(cfg *config.PdusmspConfig) ApiClient {
	c := &ApiClientInfo{
		nodeInfo:      cfg.NodeInfo,
		openApiClient: make(map[string]*openapi_commn_client.APIClient),
	}
	N11AmfNodes := cfg.N11AmfNodes
	for i := 0; i < len(N11AmfNodes); i++ {
		AmfServerAddress := N11AmfNodes[i].NodeId
		AmfServerPort := N11AmfNodes[i].Port
		klog.Info("Server address and port is ", AmfServerAddress, " ", AmfServerPort)
		var OpenApiCfg *openapi_commn_client.Configuration = openapi_commn_client.NewConfiguration(AmfServerAddress, AmfServerPort)
		c.openApiClient[AmfServerAddress] = openapi_commn_client.NewAPIClient(OpenApiCfg)
	}
	return c
}

// called by sm
func (a *ApiClientInfo) N1N2MessageTransfer(sourceIP string, ueCtxId string, n1N2MessageTransferReqData openapi_commn_client.N1N2MessageTransferReqData, binaryDataN1MessageContentFile []byte, binaryDataN2InfoContentFile []byte) {

	klog.Info("In N1N2MessageTransfer func in apiclient.go")

	ueContextId := ueCtxId
	// n1n2MessageTransferReqDataJsonFile, err := os.Open("/home/ubuntu/go/src/w5gc.io/wipro5gcore/testdata/n1n2MessageTransferReqData.json")
	// if err != nil {
	// 	fmt.Println("Error in opening sm context create file...", err)
	// } else {
	// 	fmt.Println("Successfully opened sm context json file n1n2MessageTransferReqData.json...")
	// }
	// defer n1n2MessageTransferReqDataJsonFile.Close()

	// read opened JSON file as a byte array.
	// byteValue, err := io.ReadAll(n1n2MessageTransferReqDataJsonFile)
	// fmt.Println("json data as string:", string(byteValue))
	// if err != nil {
	// 	fmt.Println("Error in reading sm context data...", err)
	// } else {
	// 	fmt.Println("Successfully read sm context data...")
	// }

	// var n1N2MessageTransferReqData openapi_commn_client.N1N2MessageTransferReqData

	// klog.Infof("%s", byteValue)
	// e := json.Unmarshal(byteValue, &n1N2MessageTransferReqData)
	// klog.Infof("n1n2Message:%+v", n1N2MessageTransferReqData)
	klog.Infof("n1container:%+v", n1N2MessageTransferReqData.N1MessageContainer)
	klog.Infof("n2container:%+v", n1N2MessageTransferReqData.N2InfoContainer)
	klog.Infof("old guami: %+v", n1N2MessageTransferReqData.OldGuami)
	klog.Infof("pdu session id : %+v", n1N2MessageTransferReqData.PduSessionId)
	klog.Infof("ue context id : %+v", ueContextId)
	// x, err := json.Marshal(n1N2MessageTransferReqData)
	// if err != nil {
	// 	klog.Info("Error in marshaling request data into JSON")
	// }
	// klog.Infof("json data before sending req : %s", x)
	apiN1N2MessageTransferRequest := a.openApiClient[sourceIP].N1N2MessageCollectionDocumentAPI.N1N2MessageTransfer(context.Background(), ueContextId).N1N2MessageTransferReqData(n1N2MessageTransferReqData)

	// binaryDataN1MessageContentFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n1msgtest")
	// if e == nil {
	// 	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN1MessageContent(binaryDataN1MessageContentFile)
	// } else {
	// 	klog.Info("Error in reading N1 message file:", e)
	// }
	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN1MessageContent(binaryDataN1MessageContentFile)
	// defer binaryDataN1MessageContentFile.Close()

	// binaryDataN2InfoContentFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n2infotest")
	// binaryDataN2InfoContentFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n2infotest")
	// if e == nil {
	// 	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN2InfoContent(binaryDataN2InfoContentFile)
	// } else {
	// 	klog.Info("Error in reading N2 information file:", e)
	// }
	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN2InfoContent(binaryDataN2InfoContentFile)
	// defer binaryDataN2InfoContentFile.Close()

	resp, httpRes, err := apiN1N2MessageTransferRequest.Execute()
	x, err := json.Marshal(n1N2MessageTransferReqData)
	klog.Infof("JSON data after sending req : %s", x)
	klog.Info(resp)
	klog.Info(httpRes)
	klog.Info(err)
}

func (a *ApiClientInfo) Start() {
	klog.Infof("Started AMF csc API client")
}
