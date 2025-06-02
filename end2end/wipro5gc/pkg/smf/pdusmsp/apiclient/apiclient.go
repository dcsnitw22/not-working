package apiclient

import (
	"context"
	"time"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
)

type ApiClient interface {
	Start()
	N1N2MessageTransfer(sourceIP string, supi string, n1N2MessageTransferReqData openapi_commn_client.N1N2MessageTransferReqData, binaryDataN1MessageContentFile []byte, binaryDataN2InfoContentFile []byte)
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
func (a *ApiClientInfo) N1N2MessageTransfer(sourceIP string, supi string, n1N2MessageTransferReqData openapi_commn_client.N1N2MessageTransferReqData, binaryDataN1MessageContentFile []byte, binaryDataN2InfoContentFile []byte) {
	//uecontextid takes value of supi according to spec
	ueContextId := supi
	klog.Infof("n1n2Message:%+v", n1N2MessageTransferReqData)
	apiN1N2MessageTransferRequest := a.openApiClient[sourceIP].N1N2MessageCollectionDocumentAPI.N1N2MessageTransfer(context.Background(), ueContextId).N1N2MessageTransferReqData(n1N2MessageTransferReqData)
	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN1MessageContent(binaryDataN1MessageContentFile)
	apiN1N2MessageTransferRequest = apiN1N2MessageTransferRequest.BinaryDataN2InfoContent(binaryDataN2InfoContentFile)
	resp, httpRes, err := apiN1N2MessageTransferRequest.Execute()
	// x, err := json.Marshal(n1N2MessageTransferReqData)
	// klog.Infof("JSON data after sending req : %s", x)
	klog.Info(resp)
	klog.Info(httpRes)
	klog.Info(err)
}

func (a *ApiClientInfo) Start() {
	klog.Infof("Started AMF csc API client")
}
