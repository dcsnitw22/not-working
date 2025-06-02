package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_server"
	"w5gc.io/wipro5gcore/openapi/openapiserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/cs"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcclient"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/ngapNas"
)

const (
	ApiChannelCapacity = 100
)

var N1Msg, N2Info chan ([]byte)

type ApiServer interface {
	Start(*grpc.Grpc)
	WatchApiChannel() chan *N1N2Message
}

type N1N2Message struct {
	N1N2Msg                    cs.N1N2Message
	UeContextID                string
	BinaryDataN1MessageContent *os.File
	BinaryDataN2InfoContent    *os.File
}

type ApiServerInfo struct {
	serverStartTime time.Time
	apiChannel      chan *N1N2Message
	nodeInfo        config.AmfNodeInfo
}

func NewApiServer(cfg config.AmfNodeInfo) ApiServer {
	return &ApiServerInfo{
		nodeInfo:   cfg,
		apiChannel: make(chan *N1N2Message, ApiChannelCapacity),
	}
}

func (a *ApiServerInfo) Start(g *grpc.Grpc) {
	klog.Infof("Starting AMF Communication API server")
	router := NewRouter(a.Routes())
	klog.Infof("Started the server on Port: %v", a.nodeInfo.ApiPort)
	klog.Fatal(http.ListenAndServe(a.nodeInfo.ApiPort, router))
}

func (a *ApiServerInfo) WatchApiChannel() chan *N1N2Message {
	return a.apiChannel
}

func (a *ApiServerInfo) Routes() Routes {
	return Routes{
		"N1N2MessageTransfer": Route{
			strings.ToUpper("Post"),
			"/namf-comm/v1/ue-contexts/{ueContextId}/n1-n2-messages",
			a.N1N2MessageTransfer,
		},
	}
}

// N1N2MessageTransfer - Namf_Communication N1N2 Message Transfer (UE Specific) service Operation
func (a *ApiServerInfo) N1N2MessageTransfer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ueContextIdParam := params["ueContextId"]

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapi_commn_server.ParsingError{Err: err}, nil)
		return
	}

	jsonDataParam := r.FormValue("jsonData")

	n1N2MessageTransferReqDataParam := openapi_commn_server.N1N2MessageTransferReqData{}

	json.Unmarshal([]byte(jsonDataParam), &n1N2MessageTransferReqDataParam)

	klog.Infof("Json data is:%v", jsonDataParam)
	klog.Infof("Input data is : %+v", n1N2MessageTransferReqDataParam)

	// n1N2MessageTransferReqDataParam := openapi_commn_server.N1N2MessageTransferReqData{}
	// d := json.NewDecoder(r.Body)
	// d.DisallowUnknownFields()
	// if err := d.Decode(&n1N2MessageTransferReqDataParam); err != nil {
	// 	openapi_commn_server.DefaultErrorHandler(w, r, &openapi_commn_server.ParsingError{Err: err}, nil)
	// 	klog.Infof("decoder gave error")
	// 	return
	// }
	// klog.Info(n1N2MessageTransferReqDataParam)
	if err := openapi_commn_server.AssertN1N2MessageTransferReqDataRequired(n1N2MessageTransferReqDataParam); err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, err, nil)
		return
	}
	if err := openapi_commn_server.AssertN1N2MessageTransferReqDataConstraints(n1N2MessageTransferReqDataParam); err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, err, nil)
		return
	}

	klog.Infof("Data Checks passed")

	_, fileHeader, err := r.FormFile("binaryDataN1MessageContent")
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	formFile, err := fileHeader.Open()
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	defer formFile.Close()
	binaryDataN1MessageContent, err := ioutil.ReadAll(formFile)
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data: %v", binaryDataN1MessageContent)

	_, fileHeader2, err := r.FormFile("binaryDataN2InfoContent")
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	formFile2, err := fileHeader2.Open()
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	defer formFile2.Close()
	binaryDataN2InfoContent, err := ioutil.ReadAll(formFile2)
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	// binaryDataN1SmMessageParam, err := ReadFormFileToTempFile(r, "binaryDataN1SmMessage")
	// if err != nil {
	//      //a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	//      openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	//      return
	// }
	klog.Infof("Binary Data: %v", binaryDataN2InfoContent)

	/*binaryDataN1MessageContentParam, err := ReadFormFileToTempFile(r, "binaryDataN1MessageContent")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapi_commn_server.DefaultErrorHandler(w, r, &openapi_commn_server.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN1MessageContentParam)

	binaryDataN2InfoContentParam, err := ReadFormFileToTempFile(r, "binaryDataN2InfoContent")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapi_commn_server.DefaultErrorHandler(w, r, &openapi_commn_server.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN2InfoContentParam)*/

	//added variable s here
	//TODO integration with backend services will be done later
	var s openapi_commn_server.N1N2MessageCollectionDocumentAPIService
	result, err := s.N1N2MessageTransfer(r.Context(), ueContextIdParam, n1N2MessageTransferReqDataParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		openapi_commn_server.DefaultErrorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	/*N1Msg = make(chan ([]byte))
	N2Info = make(chan ([]byte))
	N1Msg <- binaryDataN1MessageContent
	N2Info <- binaryDataN2InfoContent*/
	klog.Info("Sending n1n2 data in channel")
	n1n2 := &ngapNas.DLRequest{
		N1DataBytes:  binaryDataN1MessageContent,
		N2DataBytes:  binaryDataN2InfoContent,
		UeContextId:  ueContextIdParam,
		PduSessionId: n1N2MessageTransferReqDataParam.PduSessionId,
	}
	grpcclient.N1N2Data <- n1n2
}
