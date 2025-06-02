package apiclient

import (
	"io"

	//"strconv"
	"time"

	//"net"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/openapi/openapiclient"
	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"
	"w5gc.io/wipro5gcore/pkg/amf/metrics"

	// "w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"

	// "w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/create_sm_context_grpc"
	"w5gc.io/wipro5gcore/pkg/amf/csp/sm"
)

/*var AnType openapi_commn_client.AccessType
var Snssai openapi_commn_client.Snssai
var NrLocation *openapiclient.NrLocation
var PduSessionId int32
var N1SmContainer []byte
*/

var Data grpcserver.GrpcMessageInfo

const (
	ApiChannelCapacity = 100
)

type SessionMessage struct {
	MsgType    sm.MessageType
	SessionMsg sm.SMContextMessage
}

type ApiClient interface {
	Start(grpcMsg *grpcserver.GrpcMessage)
	WatchApiChannel() chan *SessionMessage
}

type ApiClientInfo struct {
	clientStartTime time.Time
	apiChannel      chan *SessionMessage
	nodeInfo        config.AmfNodeInfo
	openApiClient   *openapiclient.APIClient
	sessionDb       *redis.Client
}

// changed this paramter to type config.CspConfig instead of config.AmfNodeInfo
func NewApiClient(cfg *config.CspConfig) ApiClient {
	c := &ApiClientInfo{
		apiChannel: make(chan *SessionMessage, ApiChannelCapacity),
		nodeInfo:   cfg.NodeInfo,
	}
	//need config.N11SmfNodesInfo to create OpenApiCfg with server address and server URL
	//need to change NewApiClient parameters
	N11SmfNodes := cfg.N11SmfNodes
	defaultSmfServerAddress := N11SmfNodes[0].NodeId
	defaultSmfServerPort := N11SmfNodes[0].Port
	klog.Info("Server address and port is ", defaultSmfServerAddress, " ", defaultSmfServerPort)
	var OpenApiCfg *openapiclient.Configuration = openapiclient.NewConfiguration(defaultSmfServerAddress, defaultSmfServerPort)
	c.openApiClient = openapiclient.NewAPIClient(OpenApiCfg)

	ctx := context.Background()
	sessionClient, err := sm.NewRedisClient(ctx, 0) //changed from 2 to 0
	if err != nil {
		// Handle error
		klog.Error("unable to connect to session Database")
		klog.Info(err)
	}
	c.sessionDb = sessionClient
	return c
}

func (a *ApiClientInfo) CreateSMContext(grpcMsg *grpcserver.GrpcMessage) {

	metrics.AmfcreateSessionAttempts.Inc()

	var ctx context.Context
	createRequest := a.openApiClient.SMContextsCollectionAPI.PostSmContexts(ctx)
	hd, _ := os.UserHomeDir()
	smContextCreateJsonFile, err := os.Open(hd + "/wipro5gc/testdata/smContextCreate.json")
	if err != nil {
		klog.Info("Error in opening sm context create file...", err)
	} else {
		klog.Info("Successfully opened sm context json file smcontextcreate.json...")
	}
	defer smContextCreateJsonFile.Close()

	// read opened JSON file as a byte array.
	byteValue, err := io.ReadAll(smContextCreateJsonFile)
	// klog.Info("json data as string:", string(byteValue))
	if err != nil {
		klog.Info("Error in reading sm context data...", err)
	} else {
		klog.Info("Successfully read sm context data...")
	}
	var smContextCreateData openapiclient.SmContextCreateData
	e := json.Unmarshal(byteValue, &smContextCreateData)
	if e != nil {
		klog.Error(e)
	}

	if grpcMsg != nil {
		y := *grpcMsg.GrpcMsg
		// x := y.(create_sm_context_grpc.CreateSmContextDataFromNasMod)
		x := y.(*create_sm_context_grpc.CreateSmContextDataFromNasMod)
		// fmt.Printf("create sm context data from nas to csp : +%v", x)
		// x := create_sm_context_grpc.CreateSmContextDataFromNasMod{}
		// e := json.Unmarshal(z, &x)
		if e != nil {
			klog.Error(e)
		} else {
			klog.Infof("create sm context data from nas to csp : +%v", x)
		}
		b := openapiclient.AccessType(x.AnType)
		c := openapiclient.Snssai{
			Sst: x.Snssai.Sst,
			Sd:  &x.Snssai.Sd,
		}
		t := x.NrLocation.UeLocationTimestamp.AsTime()
		d := &openapiclient.NrLocation{
			Tai: openapiclient.Tai{
				PlmnId: openapiclient.PlmnId{
					Mcc: x.NrLocation.Tai.PlmnId.Mcc,
					Mnc: x.NrLocation.Tai.PlmnId.Mnc,
				},
				Tac: x.NrLocation.Tai.Tac,
			},
			Ncgi: openapiclient.Ncgi{
				PlmnId: openapiclient.PlmnId{
					Mcc: x.NrLocation.Ncgi.PlmnId.Mcc,
					Mnc: x.NrLocation.Ncgi.PlmnId.Mnc,
				},
				NrCellId: x.NrLocation.Ncgi.NrCellId,
			},
			// AgeOfLocationInformation: &x.NrLocation.AgeOfLocationInformation,
			UeLocationTimestamp: &t,
			/*
				GeographicalInformation:  &x.NrLocation.GeographicalInformation,
				GeodeticInformation:      &x.NrLocation.GeodeticInformation,
				GlobalGnbId:              x.NrLocation.GlobalGnbId,
			*/
		}
		f := x.PduSessionId
		g := x.N1SmContainer
		h := x.Supi
		klog.Info("n1 sm container data : ", g)
		smContextCreateData.AnType = b
		smContextCreateData.SNssai = &c
		smContextCreateData.UeLocation.NrLocation = d
		smContextCreateData.PduSessionId = &f
		smContextCreateData.Supi = &h

		createRequest = createRequest.JsonData(smContextCreateData)
		createRequest = createRequest.BinaryDataN1SmMessage(g)
		klog.Infoln("N1 SM Container :", g)
		/*
			e = os.WriteFile(hd+"/wipro5gc/testdata/PDUsessionEstNasPdu", g, os.ModeAppend)
			if e != nil {
				klog.Info("cannot write n1SmContainer to file")
			}
		*/
	}
	klog.Infof("json data as string: %+v", smContextCreateData)
	// file, e := os.Open(hd + "/wipro5gc/testdata/n1SmContainer")
	/*if e == nil {
		createRequest = createRequest.BinaryDataN1SmMessage(g)
	} else {
		klog.Info("Error in reading N1 message file:", e)
	}*/
	createdData, response, _ := createRequest.Execute()
	// smContextRef := strconv.Itoa(int(*createdData.PduSessionId)) + smContextCreateData.Guami.AmfId
	// a.sessionDb.Set(ctx, smContextRef, createdData, 0)

	metrics.AmfcreateProcess.Inc()
	metrics.AmfcreateSessionSuccess.Inc()

	klog.Info(createdData)
	klog.Info(response)
}

func (a *ApiClientInfo) ReleaseSMContext() {
	ctx := context.Background()
	smContextRef := "70218A9E"
	if a.sessionDb.Get(ctx, smContextRef).Err() == redis.Nil {
		klog.Info("SM Context Reference does not exist.")
		return
	}
	releaseSmContextRequest := a.openApiClient.IndividualSMContextAPI.ReleaseSmContext(ctx, smContextRef)
	//var smContextReleaseData openapi.SmContextReleaseData
	smContextReleaseJsonFile, err := os.Open("/home/ubuntu/wipro5gc/testdata/smContextRelease.json")
	if err != nil {
		klog.Info("Error in opening sm context release file...", err)
	} else {
		klog.Info("Successfully opened sm context json file smcontextrelease.json...")
	}
	defer smContextReleaseJsonFile.Close()
	// read opened JSON file as a byte array.
	byteValue, err := ioutil.ReadAll(smContextReleaseJsonFile)
	klog.Info("json data as string:", string(byteValue))
	if err != nil {
		klog.Info("Error in reading sm context data...", err)
	} else {
		klog.Info("Successfully read sm context data...")
	}
	var smContextReleaseData openapiclient.SmContextReleaseData
	json.Unmarshal(byteValue, &smContextReleaseData)
	klog.Info("Sm context release data:", smContextReleaseData)
	releaseSmContextRequest = releaseSmContextRequest.SmContextReleaseData(smContextReleaseData)
	binaryDataN2SmInformationFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n2infotest")
	if e == nil {
		releaseSmContextRequest = releaseSmContextRequest.BinaryDataN2SmInformation(binaryDataN2SmInformationFile)
	} else {
		klog.Info("Error in reading N2 information file:", e)
	}
	defer binaryDataN2SmInformationFile.Close()
	response, err := releaseSmContextRequest.Execute()
	klog.Info("Release SM Context response : ", response)
	klog.Info("Release Sm Context err : ", err)
	if err == nil {
		a.sessionDb.Del(ctx, smContextRef)
	}
	metrics.AmfreleaseProcess.Inc()
}

func (a *ApiClientInfo) RetrieveSMContext() {
	ctx := context.Background()
	smContextRef := "66218A9E"
	if a.sessionDb.Get(ctx, smContextRef).Err() == redis.Nil {
		klog.Info("SM Context Reference does not exist.")
		return
	}
	retrieveSmContextRequest := a.openApiClient.IndividualSMContextAPI.RetrieveSmContext(ctx, smContextRef)
	//var smContextRetrieveData openapi.SmContextRetrieveData
	// nonIpSupported := false
	// var smContextRetrieveData = openapiclient.SmContextRetrieveData{
	// 	TargetMmeCap: &openapiclient.MmeCapabilities{
	// 		NonIpSupported: &nonIpSupported,
	// 	},
	// }
	smContextRetrieveJsonFile, err := os.Open("/home/ubuntu/wipro5gc/testdata/smContextRetrieve.json")
	if err != nil {
		klog.Info("Error in opening sm context retrieve file...", err)
	} else {
		klog.Info("Successfully opened sm context json file smcontextretrieve.json...")
	}
	defer smContextRetrieveJsonFile.Close()
	// read opened JSON file as a byte array.
	byteValue, err := ioutil.ReadAll(smContextRetrieveJsonFile)
	klog.Info("json data as string:", string(byteValue))
	if err != nil {
		klog.Info("Error in reading sm context data...", err)
	} else {
		klog.Info("Successfully read sm context data...")
	}
	var smContextRetrieveData openapiclient.SmContextRetrieveData
	json.Unmarshal(byteValue, &smContextRetrieveData)
	klog.Info("Sm context retrieve data:", smContextRetrieveData)
	retrieveSmContextRequest = retrieveSmContextRequest.SmContextRetrieveData(smContextRetrieveData)
	smContextRetrievedData, response, err := retrieveSmContextRequest.Execute()
	klog.Info("Retrieved SM Context data : ", smContextRetrievedData)
	klog.Info("Retrieve SM Context response : ", response)
	klog.Info("Retrieve Sm Context err : ", err)

	metrics.AmfretrieveProcess.Inc()
}

func (a *ApiClientInfo) UpdateSMContext() {
	ctx := context.Background()
	smContextRef := "49218A9E"
	sessionInfo, err := a.sessionDb.Get(ctx, smContextRef).Result()
	if err == redis.Nil {
		klog.Info("SM Context Reference does not exist.")
		return
	}
	if err != nil {
		klog.Info("Internal Server Error : Redis")
		return
	}
	updateSmContextRequest := a.openApiClient.IndividualSMContextAPI.UpdateSmContext(ctx, smContextRef)

	smContextUpdateJsonFile, err := os.Open("/home/ubuntu/wipro5gc/testdata/smContextUpdate.json")
	if err != nil {
		klog.Info("Error in opening sm context update file...", err)
	} else {
		klog.Info("Successfully opened sm context json file smcontextupdate.json...")
	}
	defer smContextUpdateJsonFile.Close()
	// read opened JSON file as a byte array.
	byteValue, err := ioutil.ReadAll(smContextUpdateJsonFile)
	klog.Info("json data as string:", string(byteValue))
	if err != nil {
		klog.Info("Error in reading sm context data...", err)
	} else {
		klog.Info("Successfully read sm context data...")
	}
	var smContextUpdateData openapiclient.SmContextUpdateData
	json.Unmarshal(byteValue, &smContextUpdateData)
	klog.Info("Sm context update data:", smContextUpdateData)
	updateSmContextRequest = updateSmContextRequest.SmContextUpdateData(smContextUpdateData)

	binaryDataN1SmMessageFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n1msgtest")
	if e == nil {
		updateSmContextRequest = updateSmContextRequest.BinaryDataN1SmMessage(binaryDataN1SmMessageFile)
	} else {
		klog.Info("Error in reading N1 message file:", e)
	}
	defer binaryDataN1SmMessageFile.Close()

	binaryDataN2SmInformationFile, e := os.Open("/home/ubuntu/wipro5gc/testdata/n2infotest")
	if e == nil {
		updateSmContextRequest = updateSmContextRequest.BinaryDataN2SmInformation(binaryDataN2SmInformationFile)
	} else {
		klog.Info("Error in reading N2 information file:", e)
	}
	defer binaryDataN2SmInformationFile.Close()

	binaryDataN2SmInformationExt1File, e := os.Open("/home/ubuntu/wipro5gc/testdata/n2infoext1test")
	if e == nil {
		updateSmContextRequest = updateSmContextRequest.BinaryDataN2SmInformationExt1(binaryDataN2SmInformationExt1File)
	} else {
		klog.Info("Error in reading N2 information ext1 file:", e)
	}
	defer binaryDataN2SmInformationExt1File.Close()
	smContextUpdatedData, response, err := updateSmContextRequest.Execute()
	klog.Info("Updated SM Context data : ", smContextUpdatedData)
	klog.Info("Update SM Context response : ", response)
	klog.Info("Update Sm Context err : ", err.Error())
	if err == nil {
		var session openapiclient.SmContextCreatedData
		e := json.Unmarshal([]byte(sessionInfo), &session)
		if e != nil {
			klog.Info("Error in unmarshalling session data")
		}
		session.UpCnxState = smContextUpdatedData.UpCnxState
		session.AllocatedEbiList = smContextUpdatedData.AllocatedEbiList
		a.sessionDb.Set(ctx, smContextRef, session, 0)
	}
	metrics.AmfupdateProcess.Inc()
}

func (a *ApiClientInfo) Start(grpcMsg *grpcserver.GrpcMessage) {
	klog.Infof("Started AMF csp API client")
	switch grpcMsg.MsgType {
	case "create":
		a.CreateSMContext(grpcMsg)
	case "update":
		a.UpdateSMContext()
	case "release":
		a.ReleaseSMContext()
	case "retrieve":
		a.RetrieveSMContext()
	}
}

func (a *ApiClientInfo) WatchApiChannel() chan *SessionMessage {
	return a.apiChannel
}
