package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	//"net"

	"github.com/gorilla/mux"
	"k8s.io/klog"

	//"w5gc.io/wipro5gcore/openapi"
	"w5gc.io/wipro5gcore/openapi/openapiserver"
	"w5gc.io/wipro5gcore/stubs/upfgw/config"
	"w5gc.io/wipro5gcore/stubs/upfgw/sm"
)

const (
	ApiChannelCapacity = 100
)

type SessionMessage struct {
	MsgType                       sm.MessageType
	SessionMsg                    sm.SMContextMessage
	SmContextRefID                string
	BinaryDataN1SmMessage         *os.File
	BinaryDataN2SmInformation     *os.File
	BinaryDataN2SmInformationExt1 *os.File
}

type ApiServer interface {
	Start()
	WatchApiChannel() chan *SessionMessage
}

type ApiServerInfo struct {
	serverStartTime time.Time
	apiChannel      chan *SessionMessage
	//	router          http.Handler
	nodeInfo config.SmfNodeInfo

	// individualController *IndividualSMContextAPIController
	// collectionController *SMContextsCollectionAPIController
}

// Added --> Including code from api files generated
// SMContextsCollectionAPIController binds http requests to an api service and writes the service results to the http response
// type SMContextsCollectionAPIController struct {
// 	service      openapiserver.SMContextsCollectionAPIServicer
// 	errorHandler openapiserver.ErrorHandler
// }

// SMContextsCollectionAPIOption for how the controller is set up.
//type SMContextsCollectionAPIOption func(*SMContextsCollectionAPIController)

// IndividualSMContextAPIController binds http requests to an api service and writes the service results to the http response
// type IndividualSMContextAPIController struct {
// 	service      openapiserver.IndividualSMContextAPIServicer
// 	errorHandler openapiserver.ErrorHandler
// }

// IndividualSMContextAPIOption for how the controller is set up.
//type IndividualSMContextAPIOption func(*IndividualSMContextAPIController)

func NewApiServer(cfg config.SmfNodeInfo) ApiServer {
	//Commented by Mounika --> PDU Session will be handled by HSMF and ISMF
	//      IndividualPDUSessionHSMFAPIService := NewIndividualPDUSessionHSMFAPIService()
	//      IndividualPDUSessionHSMFAPIController := NewIndividualPDUSessionHSMFAPIController(IndividualPDUSessionHSMFAPIService)

	// IndividualSMContextAPIService := openapiserver.NewIndividualSMContextAPIService()
	// IndividualSMContextAPIController := NewIndividualSMContextAPIController(IndividualSMContextAPIService)

	//      PDUSessionsCollectionAPIService := NewPDUSessionsCollectionAPIService()
	//      PDUSessionsCollectionAPIController := NewPDUSessionsCollectionAPIController(PDUSessionsCollectionAPIService)

	// SMContextsCollectionAPIService := openapiserver.NewSMContextsCollectionAPIService()
	// SMContextsCollectionAPIController := NewSMContextsCollectionAPIController(SMContextsCollectionAPIService)

	//router := NewRouter(
	//              IndividualPDUSessionHSMFAPIController,
	//	IndividualSMContextAPIController,
	//              PDUSessionsCollectionAPIController,
	//	SMContextsCollectionAPIController,
	//)

	return &ApiServerInfo{
		//	router:     router,
		nodeInfo:   cfg,
		apiChannel: make(chan *SessionMessage, ApiChannelCapacity),
		//		individualController: IndividualSMContextAPIController,
		//		collectionController: SMContextsCollectionAPIController,
	}
}

func (a *ApiServerInfo) Start() {
	klog.Infof("Started SMF pdusmsp API server")
	router := NewRouter(a.Routes())
	klog.Infof("Started the server on Port: %v", a.nodeInfo.ApiPort)
	klog.Fatal(http.ListenAndServe(a.nodeInfo.ApiPort, router))

}

// Following code added from generated code

// NewIndividualSMContextAPIController creates a default api controller
// func NewIndividualSMContextAPIController(s openapiserver.IndividualSMContextAPIServicer, opts ...IndividualSMContextAPIOption) Router {
// 	controller := &IndividualSMContextAPIController{
// 		service:      s,
// 		errorHandler: openapiserver.DefaultErrorHandler,
// 	}

// 	for _, opt := range opts {
// 		opt(controller)
// 	}

// 	return controller
// }

// NewSMContextsCollectionAPIController creates a default api controller
// func NewSMContextsCollectionAPIController(s openapiserver.SMContextsCollectionAPIServicer, opts ...SMContextsCollectionAPIOption) Router {
// 	controller := &SMContextsCollectionAPIController{
// 		service:      s,
// 		errorHandler: openapiserver.DefaultErrorHandler,
// 	}

// 	for _, opt := range opts {
// 		opt(controller)
// 	}

// 	return controller
// }

// Routes returns all the api routes for the SMContextsCollectionAPIController
func (a *ApiServerInfo) Routes() Routes {
	return Routes{
		"PostSmContexts": Route{
			strings.ToUpper("Post"),
			"/nsmf-pdusession/v1/sm-contexts",
			a.PostSmContexts,
		},
		"ReleaseSmContext": Route{
			strings.ToUpper("Post"),
			"/nsmf-pdusession/v1/sm-contexts/{smContextRef}/release",
			a.ReleaseSmContext,
		},
		"RetrieveSmContext": Route{
			strings.ToUpper("Post"),
			"/nsmf-pdusession/v1/sm-contexts/{smContextRef}/retrieve",
			a.RetrieveSmContext,
		},
		"UpdateSmContext": Route{
			strings.ToUpper("Post"),
			"/nsmf-pdusession/v1/sm-contexts/{smContextRef}/modify",
			a.UpdateSmContext,
		},
	}
}

// PostSmContexts - Create SM Context
func (a *ApiServerInfo) PostSmContexts(w http.ResponseWriter, r *http.Request) {
	klog.Info("Inside PostSmContexts function")
	klog.Info("Creating SMContext")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	jsonDataParam := r.FormValue("jsonData")

	//Next two lines are added by Mounika --> To convert json data into struct
	smContextCreateDataParam := openapiserver.SmContextCreateData{}

	json.Unmarshal([]byte(jsonDataParam), &smContextCreateDataParam)

	klog.Infof("Json data is:%v", jsonDataParam)
	klog.Infof("Input data is : %+v", smContextCreateDataParam)

	binaryDataN1SmMessageParam, err := ReadFormFileToTempFile(r, "binaryDataN1SmMessage")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN1SmMessageParam)

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_CREATE_SM_CONTEXT_REQUEST, SessionMsg: smContextCreateDataParam, SmContextRefID: "", BinaryDataN2SmInformation: nil, BinaryDataN1SmMessage: binaryDataN1SmMessageParam, BinaryDataN2SmInformationExt1: nil}

}

// ReleaseSmContext - Release SM Context
func (a *ApiServerInfo) ReleaseSmContext(w http.ResponseWriter, r *http.Request) {
	klog.Info("Inside ReleaseSmContext function")

	params := mux.Vars(r)
	smContextRefParam := params["smContextRef"]

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	jsonDataParam := r.FormValue("jsonData")

	smContextReleaseDataParam := openapiserver.SmContextReleaseData{}

	json.Unmarshal([]byte(jsonDataParam), &smContextReleaseDataParam)

	klog.Infof("Json data is:%v", jsonDataParam)
	klog.Infof("Input data is : %+v", smContextReleaseDataParam)

	if err := openapiserver.AssertSmContextReleaseDataRequired(smContextReleaseDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	if err := openapiserver.AssertSmContextReleaseDataConstraints(smContextReleaseDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	klog.Infof("Data Checks passed")

	binaryDataN2SmInformationParam, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformation")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN2SmInformationParam)

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_RELEASE_SM_CONTEXT_REQUEST, SessionMsg: smContextReleaseDataParam, SmContextRefID: smContextRefParam, BinaryDataN2SmInformation: binaryDataN2SmInformationParam, BinaryDataN1SmMessage: nil, BinaryDataN2SmInformationExt1: nil}

	//Channel <- SMCchannel{Mode: "ReleaseSmContext", ChannelData: smContextReleaseDataParam, Input: smContextRefParam, BinaryDataN1SmMessage: nil, BinaryDataN2SmInformation: binaryDataN2SmInformationParam, BinaryDataN2SmInformationExt1: nil}
}

// RetrieveSmContext - Retrieve SM Context
func (a *ApiServerInfo) RetrieveSmContext(w http.ResponseWriter, r *http.Request) {
	klog.Info("Inside RetrieveSMContext function")
	params := mux.Vars(r)
	smContextRefParam := params["smContextRef"]
	smContextRetrieveDataParam := openapiserver.SmContextRetrieveData{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&smContextRetrieveDataParam); err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Input data is: %+v", smContextRetrieveDataParam)
	if err := openapiserver.AssertSmContextRetrieveDataRequired(smContextRetrieveDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	if err := openapiserver.AssertSmContextRetrieveDataConstraints(smContextRetrieveDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	klog.Infof("Data Checks passed")

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST, SessionMsg: smContextRetrieveDataParam, SmContextRefID: smContextRefParam, BinaryDataN2SmInformation: nil, BinaryDataN1SmMessage: nil, BinaryDataN2SmInformationExt1: nil}

}

// UpdateSmContext - Update SM Context
func (a *ApiServerInfo) UpdateSmContext(w http.ResponseWriter, r *http.Request) {
	klog.Info("Inside UpdateSmContext")
	klog.Info("Update SM Context is being processed")
	params := mux.Vars(r)
	smContextRefParam := params["smContextRef"]

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	jsonDataParam := r.FormValue("jsonData")

	smContextUpdateDataParam := openapiserver.SmContextUpdateData{}

	json.Unmarshal([]byte(jsonDataParam), &smContextUpdateDataParam)

	klog.Infof("Json data is:%v", jsonDataParam)
	klog.Infof("Input data is : %+v", smContextUpdateDataParam)

	if err := openapiserver.AssertSmContextUpdateDataRequired(smContextUpdateDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	if err := openapiserver.AssertSmContextUpdateDataConstraints(smContextUpdateDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	klog.Infof("Data Checks passed")

	binaryDataN1SmMessageParam, err := ReadFormFileToTempFile(r, "binaryDataN1SmMessage")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN1SmMessageParam)

	binaryDataN2SmInformationParam, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformation")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN2SmInformationParam)

	binaryDataN2SmInformationExt1Param, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformationExt1")
	if err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	klog.Infof("Binary Data file: %v", binaryDataN2SmInformationExt1Param)

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_UPDATE_SM_CONTEXT_REQUEST, SessionMsg: smContextUpdateDataParam, SmContextRefID: smContextRefParam, BinaryDataN2SmInformation: binaryDataN2SmInformationParam, BinaryDataN1SmMessage: binaryDataN1SmMessageParam, BinaryDataN2SmInformationExt1: binaryDataN2SmInformationExt1Param}

	//Channel <- SMCchannel{Mode: "ReleaseSmContext", ChannelData: smContextUpdateDataParam, Input: smContextRefParam, BinaryDataN1SmMessage: binaryDataN1SmMessageParam, BinaryDataN2SmInformation: binaryDataN2SmInformationParam, BinaryDataN2SmInformationExt1: binaryDataN2SmInformationExt1Param}

}

func (a *ApiServerInfo) WatchApiChannel() chan *SessionMessage {
	return a.apiChannel
}
