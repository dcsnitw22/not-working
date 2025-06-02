package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	//"net"

	"github.com/gorilla/mux"
	"k8s.io/klog"

	//"w5gc.io/wipro5gcore/openapi"
	"w5gc.io/wipro5gcore/openapi/openapiserver"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/sm"
)

const (
	ApiChannelCapacity = 100
)

type SessionMessage struct {
	MsgType                       sm.MessageType
	SessionMsg                    sm.SMContextMessage
	SmContextRefID                string
	BinaryDataN1SmMessage         []byte
	BinaryDataN2SmInformation     []byte
	BinaryDataN2SmInformationExt1 []byte
	Writer                        http.ResponseWriter
	Request                       *http.Request
}

type Receiver struct {
	RecievedResponse openapiserver.ImplResponse
	RecievedErr      error
}

type ApiServer interface {
	Start()
	WatchApiChannel() chan *SessionMessage
	WatchRecChannel() chan *Receiver
}

type ApiServerInfo struct {
	serverStartTime time.Time
	apiChannel      chan *SessionMessage
	apiReceiver     chan *Receiver
	//	router          http.Handler
	nodeInfo        config.SmfNodeInfo
	RequestResponse openapiserver.ImplResponse
	ErrorResponse   error

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
		nodeInfo:    cfg,
		apiChannel:  make(chan *SessionMessage, ApiChannelCapacity),
		apiReceiver: make(chan *Receiver),
		//		individualController: IndividualSMContextAPIController,
		//		collectionController: SMContextsCollectionAPIController,
	}
}

// func GetResponse(resp openapiserver.ImplResponse,err error)

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

	//validate whether amf node ip is present or not
	//remove ip check we will di amfi check
	// AmfNodeIPAddress := r.Header.Get("X-Real-Ip")
	// if AmfNodeIPAddress == "" {
	// 	AmfNodeIPAddress = r.Header.Get("X-Forwarded-For")
	// }
	// if AmfNodeIPAddress == "" {
	// 	AmfNodeIPAddress = r.RemoteAddr
	// }
	// AmfNodeIPAddress = strings.Split(AmfNodeIPAddress, ":")[0]
	// N11AmfNodes := (config.PdusmspCfg).N11AmfNodes
	// found := false
	// klog.Info(N11AmfNodes)
	// for i := 0; i < len(N11AmfNodes); i++ {
	// 	klog.Info(N11AmfNodes[i].NodeId)
	// 	if AmfNodeIPAddress == N11AmfNodes[i].NodeId {
	// 		found = true
	// 		break
	// 	}
	// }
	// klog.Info(found)
	// if !found {
	// 	klog.Info("Request has not been sent from a peer AMF node")
	// 	err := errors.New("request has not been sent from a peer amf node")
	// 	// klog.Infof("writer: %v,\n request: %v,\n error: %v\n", w, r, err)
	// 	openapiserver.DefaultErrorHandler(w, r, err, &openapiserver.ImplResponse{
	// 		Code: http.StatusBadRequest,
	// 		Body: err.Error(),
	// 	})
	// 	return
	// }

	klog.Info("Creating SMContext")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	jsonDataParam := r.FormValue("jsonData")

	//Next two lines are added by Mounika --> To convert json data into struct
	smContextCreateDataParam := openapiserver.SmContextCreateData{}

	klog.Info(json.Unmarshal([]byte(jsonDataParam), &smContextCreateDataParam))
	klog.Infof("Json data is:%v", jsonDataParam)
	klog.Infof("Input data is : %+v", smContextCreateDataParam)

	if err := openapiserver.AssertSmContextCreateDataRequired(smContextCreateDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}
	if err := openapiserver.AssertSmContextCreateDataConstraints(smContextCreateDataParam); err != nil {
		//a.individualController.errorHandler(w, r, err, nil)
		openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
		return
	}

	klog.Infof("Data Checks passed")

	nasMsgParam := r.FormValue("binaryDataN1SmMessage")
	var binaryDataN1SmMessageParam []byte
	json.Unmarshal([]byte(nasMsgParam), &binaryDataN1SmMessageParam)

	// binaryDataN1SmMessageParam, err := ReadFormFileToTempFile(r, "binaryDataN1SmMessage")
	// if err != nil {
	// 	//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	return
	// }
	klog.Infof("Binary Data: %v", binaryDataN1SmMessageParam)
	//send writer in channel
	a.apiChannel <- &SessionMessage{
		MsgType:                       sm.NSMF_CREATE_SM_CONTEXT_REQUEST,
		SessionMsg:                    smContextCreateDataParam,
		SmContextRefID:                "",
		BinaryDataN2SmInformation:     nil,
		BinaryDataN1SmMessage:         binaryDataN1SmMessageParam,
		BinaryDataN2SmInformationExt1: nil,
		Writer:                        w,
		Request:                       r,
	}

	//TODO make a reciever function using code written below
	rec := <-a.apiReceiver
	// klog.Info(rec)

	if rec.RecievedErr != nil {
		openapiserver.DefaultErrorHandler(
			w, r, &openapiserver.ParsingError{
				Err: rec.RecievedErr,
			}, &rec.RecievedResponse,
		)
		return
	}
	EncodeJSONResponse(rec.RecievedResponse.Body, &rec.RecievedResponse.Code, w)
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

	// binaryDataN2SmInformationParam, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformation")
	// if err != nil {
	// 	//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	return
	// }

	nasMsgParam := r.FormValue("binaryDataN2SmInformation")
	var binaryDataN2SmInformationParam []byte
	json.Unmarshal([]byte(nasMsgParam), &binaryDataN2SmInformationParam)
	klog.Infof("Binary Data: %v", binaryDataN2SmInformationParam)

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_RELEASE_SM_CONTEXT_REQUEST,
		SessionMsg:                    smContextReleaseDataParam,
		SmContextRefID:                smContextRefParam,
		BinaryDataN2SmInformation:     binaryDataN2SmInformationParam,
		BinaryDataN1SmMessage:         nil,
		BinaryDataN2SmInformationExt1: nil}
	//TODO  create a function to handle situation below
	rec := <-a.apiReceiver
	// klog.Info(rec)

	if rec.RecievedErr != nil {
		openapiserver.DefaultErrorHandler(
			w, r, &openapiserver.ParsingError{
				Err: rec.RecievedErr,
			}, &rec.RecievedResponse,
		)
		return
	}
	EncodeJSONResponse(rec.RecievedResponse.Body, &rec.RecievedResponse.Code, w)

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

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST,
		SessionMsg:                smContextRetrieveDataParam,
		SmContextRefID:            smContextRefParam,
		BinaryDataN2SmInformation: nil,
		BinaryDataN1SmMessage:     nil, BinaryDataN2SmInformationExt1: nil}
	rec := <-a.apiReceiver
	// klog.Info(rec)

	if rec.RecievedErr != nil {
		openapiserver.DefaultErrorHandler(
			w, r, &openapiserver.ParsingError{
				Err: rec.RecievedErr,
			}, &rec.RecievedResponse,
		)
		return
	}
	EncodeJSONResponse(rec.RecievedResponse.Body, &rec.RecievedResponse.Code, w)

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

	// binaryDataN1SmMessageParam, err := ReadFormFileToTempFile(r, "binaryDataN1SmMessage")
	// if err != nil {
	// 	//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	return
	// }

	nasMsgParam := r.FormValue("binaryDataN1SmMessage")
	var binaryDataN1SmMessageParam []byte
	json.Unmarshal([]byte(nasMsgParam), &binaryDataN1SmMessageParam)
	klog.Infof("Binary Data: %v", binaryDataN1SmMessageParam)

	// binaryDataN2SmInformationParam, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformation")
	// if err != nil {
	// 	//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	return
	// }
	nasMsgParam2 := r.FormValue("binaryDataN2SmInformation")
	var binaryDataN2SmInformationParam []byte
	json.Unmarshal([]byte(nasMsgParam2), &binaryDataN2SmInformationParam)
	klog.Infof("Binary Data: %v", binaryDataN2SmInformationParam)

	// binaryDataN2SmInformationExt1Param, err := ReadFormFileToTempFile(r, "binaryDataN2SmInformationExt1")
	// if err != nil {
	// 	//a.individualController.errorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	openapiserver.DefaultErrorHandler(w, r, &openapiserver.ParsingError{Err: err}, nil)
	// 	return
	// }

	nasMsgParam3 := r.FormValue("binaryDataN2SmInformationExt1")
	var binaryDataN2SmInformationExt1Param []byte
	json.Unmarshal([]byte(nasMsgParam3), &binaryDataN2SmInformationExt1Param)
	klog.Infof("Binary Data: %v", binaryDataN2SmInformationExt1Param)

	a.apiChannel <- &SessionMessage{MsgType: sm.NSMF_UPDATE_SM_CONTEXT_REQUEST,
		SessionMsg:                    smContextUpdateDataParam,
		SmContextRefID:                smContextRefParam,
		BinaryDataN2SmInformation:     binaryDataN2SmInformationParam,
		BinaryDataN1SmMessage:         binaryDataN1SmMessageParam,
		BinaryDataN2SmInformationExt1: binaryDataN2SmInformationExt1Param}

	rec := <-a.apiReceiver
	// klog.Info(rec)

	if rec.RecievedErr != nil {
		openapiserver.DefaultErrorHandler(
			w, r, &openapiserver.ParsingError{
				Err: rec.RecievedErr,
			}, &rec.RecievedResponse,
		)
		return
	}
	EncodeJSONResponse(rec.RecievedResponse.Body, &rec.RecievedResponse.Code, w)

}

func (a *ApiServerInfo) WatchApiChannel() chan *SessionMessage {
	return a.apiChannel
}

func (a *ApiServerInfo) WatchRecChannel() chan *Receiver {
	return a.apiReceiver
}
