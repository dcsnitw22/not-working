// TODO separate session code and sm code
// TODO decode NAS
package sm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	openapiserver "w5gc.io/wipro5gcore/openapi/openapiserver"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/apiclient"
	db "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/database"
	redisClient "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/database/redis"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"
	grpcnas "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/grpcSmfNas/pb"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpcNAS/nas"
)

var (
	createProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_create_operations_total",
		Help: "The total number of SM create operations",
	})

	updateProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_update_operations_total",
		Help: "The total number of SM update operations",
	})

	retrieveProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_retrieve_operations_total",
		Help: "The total number of SM retrieve operations",
	})

	releaseProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_release_operations_total",
		Help: "The total number of SM release operations",
	})
	UsersPerSession = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "users_per_session_id",
			Help: "Tracks the number of users per session ID",
		},
		[]string{"session_id", "user_id", "pei"},
	)
	createSessionSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_create_session_success_total",
		Help: "The total number of successful session creations",
	})

	createSessionAttempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sm_create_session_attempts_total",
		Help: "The total number of session creation attempts",
	})
)

type SessionManager interface {
	Start()
	ProcessNsmfReleaseSmContextRequest(smContextRef string, smContextReleaseData openapiserver.SmContextReleaseData, binaryN2SMInformation []byte) (openapiserver.ImplResponse, error)
	ProcessNsmfRetrieveSmContextRequest(smContextRef string, smContextRetrieveData openapiserver.SmContextRetrieveData) (openapiserver.ImplResponse, error)
	ProcessNsmfUpdateSmContextRequest(smContextRef string, smContextUpdateData openapiserver.SmContextUpdateData, binaryDataN1SmMessage []byte, binaryN2SMInformation []byte, binaryDataN2SmInformationExt1 []byte) (openapiserver.ImplResponse, error)
	ProcessNsmfCreateSmContextRequest(jsonData openapiserver.SmContextCreateData, binaryDataN1SmMessage []byte) (openapiserver.ImplResponse, error)
	ProcessN1N2Message(grpcMsg *protos.N1N2MessageTransferDataRequest, ip string) error
}

type SMContextMessage interface{}

type MessageType uint8

const (
	NSMF_CREATE_SM_CONTEXT_REQUEST    MessageType = 1
	NSMF_CREATE_SM_CONTEXT_RESPONSE   MessageType = 2
	NSMF_UPDATE_SM_CONTEXT_REQUEST    MessageType = 3
	NSMF_UPDATE_SM_CONTEXT_RESPONSE   MessageType = 4
	NSMF_RELEASE_SM_CONTEXT_REQUEST   MessageType = 5
	NSMF_RELEASE_SM_CONTEXT_RESPONSE  MessageType = 6
	NSMF_RETRIEVE_SM_CONTEXT_REQUEST  MessageType = 7
	NSMF_RETRIEVE_SM_CONTEXT_RESPONSE MessageType = 8
	NSMF_N1_N2_TRANSFER               MessageType = 11
)

// try to have a single client
// make context as a part of sminfo
// temp struct
type SmInfo struct {
	// ctx context.Context
	// sessionDb *redis.Client // Redis client for database 0
	// userDb    *redis.Client // Redis client for database 1
	// dbClient *redis.Client
	DbClient  *db.DBInfo
	grpc      *grpc.Grpc
	apiClient apiclient.ApiClient
}

func NewSessionManager(info *db.DBInfo, grpc *grpc.Grpc, apiclient apiclient.ApiClient) *SmInfo {
	// ctx := context.Background()
	// sessionClient, err := redisClient.NewRedisClient(ctx, 0)
	// if err != nil {
	// 	// Handle error
	// 	klog.Error("unable to connect to session Database")
	// 	klog.Info(err)
	// }

	// userClient, err := redisClient.NewRedisClient(ctx, 1)
	// if err != nil {
	// 	// Handle error
	// 	klog.Error("unable to connect to user Database")
	// 	klog.Info(err)
	// }

	// return &SmInfo{
	// 	// smcontextTable:   make(map[string]SessionContext),
	// 	// userContextTable: make(map[string]UserContext),
	// 	ctx:       ctx,
	// 	sessionDb: sessionClient,
	// 	userDb:    userClient,
	// }
	// client, err := redisClient.NewRedisClient(ctx, redisClient.SessionDb)
	// if err != nil {
	// 	klog.Errorf("unable to connect to session Database,%s", err.Error())
	// 	return &SmInfo{
	// 		ctx:      ctx,
	// 		dbClient: nil,
	// 	}
	// }
	return &SmInfo{
		DbClient:  info,
		grpc:      grpc,
		apiClient: apiclient,
	}
}

func (s *SmInfo) Start() {
	klog.Info("starting session manager")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

//	func (s *SessionManagerInfo) NewSMContextAPIService() *SMContextAPIService {
//		return s.smContextAPI
//
// Helper function to convert SessionContext to JSON
// }
func SessionContextToJSON(sessionContext SessionContext) string {
	data, err := json.Marshal(sessionContext)
	if err != nil {
		klog.Errorf("Error converting SessionContext to JSON: %v", err)
		return ""
	}
	return string(data)
}

// Helper function to convert UserContext to JSON
func UserContextToJSON(userContext UserContext) string {
	data, err := json.Marshal(userContext)
	if err != nil {
		klog.Errorf("Error converting UserContext to JSON: %v", err)
		return ""
	}
	return string(data)
}

// requestType - emergeny | not taken TODO
func validateData(jsonData openapiserver.SmContextCreateData) error {

	//TODO check condition - if amf is not the sender - for later
	supiValid, _ := regexp.MatchString("^(imsi-[0-9]{5,15}|nai-.+|.+)$", jsonData.Supi)
	if !supiValid {
		return errors.New("invalid supi")
	}
	// gpsiValid, _ := regexp.MatchString("^(msisdn-[0-9]{5,15}|extid-.+@.+|.+)$", jsonData.Gpsi)
	// if !gpsiValid {
	// 	return errors.New("invalid gpsi")
	// }
	peiValid, _ := regexp.MatchString("^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$", jsonData.Pei)
	if !peiValid {
		return errors.New("invalid pei")
	}
	amfidValid, _ := regexp.MatchString("^[A-Fa-f0-9]{6}$", jsonData.Guami.AmfId)
	gMccValid, _ := regexp.MatchString("^[0-9]{3}$", jsonData.Guami.PlmnId.Mcc)
	gMncValid, _ := regexp.MatchString("^[0-9]{2,3}$", jsonData.Guami.PlmnId.Mnc)
	guamiValid := amfidValid && gMccValid && gMncValid
	if !guamiValid {
		return errors.New("invalid guami")
	}
	snMccValid, _ := regexp.MatchString("^[0-9]{3}$", jsonData.ServingNetwork.Mcc)
	snMncValid, _ := regexp.MatchString("^[0-9]{2,3}$", jsonData.ServingNetwork.Mnc)
	servingNetworkValid := snMccValid && snMncValid
	if !servingNetworkValid {
		return errors.New("invalid servingNetwork")
	}
	return nil
}

// func (s *SmInfo) changeDatabase(redisDb redisClient.Database) {
// 	s.dbClient = redisClient.ChangeRedisDatabase(s.dbClient, redisDb)
// }

// func upfSessionCreate(data openapiserver.RefToBinaryData) (openapiserver.ImplResponse, error) {
// 	return openapiserver.Response(200, nil), nil
// }

// ProcessNsmfCreateSmContextRequest - Create SM Context
func (s *SmInfo) ProcessNsmfCreateSmContextRequest(
	// ctx context.Context,
	jsonData openapiserver.SmContextCreateData, binaryDataN1SmMessage []byte) (openapiserver.ImplResponse, error) {
	// TODO - update ProcessNsmfCreateSmContextRequest with the required logic for this service method.
	// Add api_sm_contexts_collection_service.go to the .openapiserver-generator-ignore to avoid overwriting this service implementation when updating open api gen>
	klog.Info("create function initiated")
	createSessionAttempts.Inc()

	if jsonData.PduSessionId == 0 {
		return openapiserver.Response(403, openapiserver.SmContextCreateError{Error: openapiserver.ProblemDetails{Title: "invalid request data",
			Type: "validityErr", Status: 403, Detail: "invalid PduSessionId "}, N1SmMsg: jsonData.N1SmMsg}), errors.New("invalid PduSessionId")
	}

	validityErr := validateData(jsonData)
	klog.Info(validityErr)

	if validityErr != nil {
		return openapiserver.Response(403, openapiserver.SmContextCreateError{Error: openapiserver.ProblemDetails{Title: "invalid request data",
			Type: "validityErr", Status: 403, Detail: validityErr.Error()}, N1SmMsg: jsonData.N1SmMsg}), validityErr
	}

	// get connection attributes related to the session
	/*	user, found := s.userContextTable[jsonData.Supi]
		if found {
			//assumingmax allowed devices is 15
			if user.NoSession == 15 {
				// recheck in specs
				return openapiserver.Response(401, openapiserver.SmContextCreateError{Error: openapiserver.ProblemDetails{Title: "Connection limit exceeded",
					Type: "ConnectionLimitErr", Status: 401, Detail: "device not allowed to have more than 15 connections"},
					N1SmMsg: jsonData.N1SmMsg}), errors.New("device not allowed to have more than 15 connections")
			}
		}
	*/
	//contextRefId is PduSessionId concatenated with SUPI
	klog.Info("Creating ContextRefId")
	// contextRefId := strconv.Itoa(int(jsonData.PduSessionId)) + jsonData.Guami.AmfId
	contextRefId := strconv.Itoa(int(jsonData.PduSessionId)) + jsonData.Supi
	var sessionData SessionContext

	// sData, exists := s.sessionDb.Get(s.ctx, contextRefId).Result()
	// s.changeDatabase(redisClient.SessionDb)
	// sData, exists := s.dbClient.Get(s.ctx, contextRefId).Result()
	// sData, exists := redisClient.Read(contextRefId, redisClient.SessionDb, s.dbClient)
	sData, exists := s.DbClient.Redis.Read(contextRefId, redisClient.SessionDb)
	klog.Info("sessionContext: ", sData)
	json.Unmarshal([]byte(sData), &sessionData)

	//TODO Should request be buffered to handle later
	if exists == nil {
		if sessionData.State == openapiserver.ACTIVATING {
			// klog.Info("inside")
			return openapiserver.Response(403, openapiserver.SmContextCreateError{
				Error: openapiserver.ProblemDetails{
					Title:  "Already in Progress",
					Type:   "ALreadyInProgressErr",
					Status: 403,
					Detail: "request already in progress",
				},
				N1SmMsg: jsonData.N1SmMsg,
			}), errors.New("request already in progress")
		}
	}

	var user UserContext

	//TODO only use one call to both function (put it in redisClient)
	// dbData, err := s.userDb.Get(s.ctx, jsonData.Supi).Result()
	// s.changeDatabase(redisClient.UserDb)
	// userData, err := s.dbClient.Get(s.ctx, jsonData.Supi).Result()
	userData, err := s.DbClient.Redis.Read(jsonData.Supi, redisClient.UserDb)
	// klog.Infof("UserContext,err: %+v,%+v", userData, )

	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "failed to get data in Redis: redis: nil" {
			user = UserContext{NoSession: 0}
		} else {
			klog.Info("iwashere")
			return openapiserver.Response(500, openapiserver.SmContextCreateError{
				Error: openapiserver.ProblemDetails{
					Title:  "Unable to access redis data",
					Type:   "RedisDatabaseError",
					Status: 500,
					Detail: err.Error(),
				},
				N1SmMsg: jsonData.N1SmMsg,
			}), err
		}
	} else {
		json.Unmarshal([]byte(userData), &user)
	}
	//assumingmax allowed devices is 15
	if user.NoSession == 15 {
		return openapiserver.Response(401, openapiserver.SmContextCreateError{
			Error: openapiserver.ProblemDetails{
				Title:  "Connection limit exceeded",
				Type:   "ConnectionLimitErr",
				Status: 401,
				Detail: "device not allowed to have more than 15 connections",
			},
			N1SmMsg: jsonData.N1SmMsg,
		}), errors.New("device not allowed to have more than 15 connections")
	}

	ueContextId := jsonData.Supi
	//grpc code
	createData := protos.SmContextCreateDataRequest{
		SmContextID: contextRefId,
		N4ReqType:   "N4 Establishment Request",
		// Supi:         jsonData.Supi,
		// PduSessionId: jsonData.PduSessionId,
		// Guami: &protos.Guami{
		// 	PlmnId: &protos.PlmnId{
		// 		Mcc: jsonData.Guami.PlmnId.Mcc,
		// 		Mnc: jsonData.Guami.PlmnId.Mnc,
		// 	},
		// 	AmfId: jsonData.Guami.AmfId,
		// },
		// ServingNfId: jsonData.ServingNfId,
		// UnauthenticatedSupi: jsonData.UnauthenticatedSupi,
		// Pei:                 jsonData.Pei,
		// Gpsi:                jsonData.Gpsi,
		// Dnn:                 jsonData.Dnn,
		// ServingNetwork: &protos.PlmnId{
		// 	Mcc: jsonData.ServingNetwork.Mcc,
		// 	Mnc: jsonData.ServingNetwork.Mnc,
		// },
		// RequestType: string(jsonData.RequestType),
		// N1SmMessage: &protos.N1SmMessage{
		// 	PduSessionEstablishmentRequest: &protos.PduSessionEstablishmentRequest{
		// 		PduSessionId: jsonData.PduSessionId,
		// 	},
		// },
		// AnType:             string(jsonData.AnType),
		// RatType:            string(jsonData.RatType),
		// SmContextStatusUri: jsonData.SmContextStatusUri,
	}
	(*s.grpc).SendSmContextCreateData(&createData)

	//Decode Nas Message
	// var nasMsg interface{}

	anyMessage, err := anypb.New(&pb.ByteDataWrapper{ByteArray: binaryDataN1SmMessage})
	if err != nil {
		klog.Fatalf("Failed to create Any message: %v", err)
	}

	reqType := "" // Empty string for the request type

	client, ctx := grpcnas.CreateGRPCNasClient()

	// Send the request
	resp, err := client.SendSMData(ctx, &pb.SMDataRequest{
		NasMessage: anyMessage,
		TypeReq:    reqType,
	})
	if err != nil {
		klog.Fatalf("Error sending SMData: %v", err)
	}

	klog.Infof("Decoded binaryDataN1SmMessage: %v", resp.NasResponse)

	var nasResp pb.PDUSEstReqModel

	err = json.Unmarshal(resp.NasResponse.Value, &nasResp)
	if err != nil {
		klog.Fatalf("Error unmarshalling NAS decoded Message: %v", err)
	}

	sessionData = SessionContext{
		Supi:                               jsonData.Supi,
		Pei:                                jsonData.Pei,
		ServingNfId:                        jsonData.ServingNfId,
		State:                              openapiserver.ACTIVATING,
		SubState:                           ACTIVATING_CREATE_IN_PROGRESS,
		UnauthenticatedSupi:                jsonData.UnauthenticatedSupi,
		Gpsi:                               jsonData.Gpsi,
		FDnn:                               jsonData.Dnn,
		Guami:                              jsonData.Guami,
		ServiceName:                        jsonData.ServiceName,
		ServingNetwork:                     jsonData.ServingNetwork,
		AnType:                             jsonData.AnType,
		RatType:                            jsonData.RatType,
		SmContextStatusUri:                 jsonData.SmContextStatusUri,
		PduSessionId:                       jsonData.PduSessionId,
		ContextRefId:                       contextRefId,
		UeContextId:                        ueContextId,
		NASepd:                             nasResp.Epd,
		NASpduSessionId:                    int(*nasResp.PdusessionID),
		NASpti:                             int(*nasResp.Pti),
		NASmsgType:                         nasResp.MsgType,
		NASmaxIntegrityProtectedDataRateUL: nasResp.DatarateUL,
		NASmaxIntegrityProtectedDataRateDL: nasResp.DatarateDL,
	}

	// Store data in Redis database 0
	//        err := s.sessionDb.Set(ctx,, "Keshav", 0).Err()
	// contextRefId := strconv.Itoa(int(jsonData.PduSessionId)) + jsonData.Guami.AmfId
	// s.changeDatabase(redisClient.SessionDb)
	// err = s.dbClient.Set(s.ctx, contextRefId, SessionContextToJSON(sessionData), 0).Err()

	// err = s.sessionDb.Set(s.ctx, contextRefId, SessionContextToJSON(sessionData), 0).Err()

	_, err = s.DbClient.Redis.Create(
		contextRefId,
		SessionContextToJSON(sessionData),
		redisClient.SessionDb,
	)

	if err != nil {
		klog.Errorf("Error storing data in Redis database 0: %v", err)
		return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	//send to channel to upf (ask raghu)

	user = UserContext{user.NoSession + 1}
	// s.changeDatabase(redisClient.UserDb)
	// err = s.dbClient.Set(s.ctx, jsonData.Supi, UserContextToJSON(user), 0).Err()
	// err = s.userDb.Set(s.ctx, jsonData.Supi, UserContextToJSON(user), 0).Err()

	_, err = s.DbClient.Redis.Create(
		jsonData.Supi,
		UserContextToJSON(user),
		redisClient.UserDb,
	)

	if err != nil {
		return openapiserver.Response(500, openapiserver.SmContextCreateError{
			Error: openapiserver.ProblemDetails{
				Title:  "Unable to Store Session data",
				Type:   "RedisDatabaseError",
				Status: 500,
				Detail: err.Error(),
			},
			N1SmMsg: jsonData.N1SmMsg,
		}), err
	}

	createProcess.Inc()

	pduSessionIdStr := fmt.Sprintf("%d", sessionData.PduSessionId)
	UsersPerSession.WithLabelValues(pduSessionIdStr, sessionData.Supi, sessionData.Pei).Inc()

	createSessionSuccess.Inc()

	// Increment the UsersPerSession metric
	// UsersPerSession.WithLabelValues(pduSessionIdStr, sessionData.Supi+"_"+sessionData.Pei).Inc()
	//On success, "201 Created" shall be returned, the payload body of the POST response shall contain the representation describing the status of the request and the "Location" header shall be present and shall contain the URI of the created resource. The authority and/or deployment-specific string of the apiRoot of the created resource URI may differ from the authority and/or deployment-specific string of the apiRoot of the request URI received in the POST request.
	return openapiserver.Response(201, openapiserver.SmContextCreatedData{
		UpCnxState:   sessionData.State,
		PduSessionId: sessionData.PduSessionId,
	}), nil

}

// ReleaseSmContext - Release SM Context
func (s *SmInfo) ProcessNsmfReleaseSmContextRequest(smContextRef string, smContextReleaseData openapiserver.SmContextReleaseData, binaryDataN2SmInformation []byte) (openapiserver.ImplResponse, error) {
	klog.Info("Release initiated")
	// Returns a `400 Bad Request` response if the smContextRef is empty.
	if smContextRef == "" {
		return openapiserver.Response(http.StatusBadRequest, nil), errors.New("smContextRef must not be empty")
	}
	//TODO activated_release_in_progress case handling

	// Check if the SM context exists
	// Returns a `404 Not Found` response if the SM context does not exist.
	var smcontext SessionContext

	// dbData, err := s.sessionDb.Get(s.ctx, smContextRef).Result()
	// s.changeDatabase(redisClient.SessionDb)
	// dbData, err := s.dbClient.Get(s.ctx, smContextRef).Result()
	sData, err := s.DbClient.Redis.Read(smContextRef, redisClient.SessionDb)

	klog.Info(sData, err)

	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "failed to get data in Redis: redis: nil" {
			klog.Errorf("Incoorect smcontextRef: %v", smContextRef)
			return openapiserver.Response(
				http.StatusNotFound,
				openapiserver.ProblemDetails{},
			), errors.New("smContextRef not a valid input")
		} else {
			klog.Error(err.Error())
			return openapiserver.Response(
				http.StatusInternalServerError,
				openapiserver.ProblemDetails{},
			), errors.New("Internal server error" + err.Error())
		}

	} else {
		json.Unmarshal([]byte(sData), &smcontext)
	}
	if smcontext.SubState == ACTIVATED_RELEASE_IN_PROGRESS {
		return openapiserver.Response(
			http.StatusNotFound,
			nil,
		), errors.New("release already in progress")
	}

	//TODO Release the IP addresses/prefixes and User Plane resources
	//TODO smContextReleaseData.Cause=="PDU_SESSION_STATUS_MISMATCH"
	//get some info from ngap message
	releaseData := protos.SmContextReleaseDataRequest{
		ServingNfId:  smcontext.ServingNfId,
		Pei:          smcontext.Pei,
		PduSessionId: smcontext.PduSessionId,
		Guami: &protos.Guami{
			PlmnId: &protos.PlmnId{
				Mcc: smcontext.Guami.PlmnId.Mcc,
				Mnc: smcontext.Guami.PlmnId.Mnc,
			},
			AmfId: smcontext.Guami.AmfId,
		},
		SmContextStatusUri: smcontext.SmContextStatusUri,
		ServingNetwork: &protos.PlmnId{
			Mcc: smcontext.ServingNetwork.Mcc,
			Mnc: smcontext.ServingNetwork.Mnc,
		},
	}
	(*s.grpc).SendSmContextReleaseData(&releaseData)

	// err = s.ReleaseUPFResources(smcontext.N4SessionID)
	// if err != nil {
	// 	// Handle error releasing UPF resources
	// 	return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("error releasing UPF resources")
	// }

	var user UserContext

	// val, _ := s.userDb.Get(s.ctx, smcontext.Supi).Result()
	// s.changeDatabase(redisClient.UserDb)
	// val, _ := s.dbClient.Get(s.ctx, smcontext.Supi).Result()

	uData, err := s.DbClient.Redis.Read(smcontext.Supi, redisClient.UserDb)
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "failed to get data in Redis: redis: nil" {
			klog.Info("New User")
			user = UserContext{NoSession: 0}
		} else {
			klog.Error(err.Error())
			return openapiserver.Response(
				http.StatusNotFound,
				openapiserver.ProblemDetails{},
			), err
		}
	} else {
		json.Unmarshal([]byte(uData), &user)
	}
	user = UserContext{user.NoSession - 1}

	_, err = s.DbClient.Redis.Create(smcontext.Supi, UserContextToJSON(user), redisClient.UserDb)

	//err:= s.dbClient.Set(s.ctx, smContextRef, UserContextToJSON(user), 0).Err()
	// err := s.userDb.Set(s.ctx, smcontext.Supi, UserContextToJSON(user), 0).Err()
	//	err = s.userDb.Set(s.ctx, smContextRef, UserContextToJSON(UserContext{user.NoSession - 1}), 0).Err()

	if err != nil {
		klog.Error(err.Error())
		return openapiserver.Response(500, nil), err
	}

	// Release the SM context

	// Remove data from Redis database 0

	// err = s.sessionDb.Del(s.ctx, smContextRef).Err()
	// s.changeDatabase(redisClient.SessionDb)
	// err= s.dbClient.Del(s.ctx, smContextRef).Err()
	_, err = s.DbClient.Redis.Delete(smContextRef, redisClient.SessionDb)

	if err != nil {
		klog.Errorf("Error removing data from Redis database 0 for smContextRef %s: %v", smContextRef, err)
		return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	releaseProcess.Inc()

	pduSessionIdStr := fmt.Sprintf("%d", smcontext.PduSessionId)

	// Decrement the Prometheus metric with session_id and user_id
	UsersPerSession.WithLabelValues(pduSessionIdStr, smcontext.Supi, smcontext.Pei).Dec()

	//why is this commented ask deep

	// Update data in Redis database 1
	/*		val1, err := UserContextToJSON(UserContext{val.NoSession - 1})
		if err != nil {
			klog.Errorf("Error converting UserContext to JSON for updating in Redis database 1: %v", err)
			return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		err = s.userDb.HSet(s.ctx,"userContextTable", smcontext.Supi, val1).Err()
		if err != nil {
			klog.Errorf("Error updating data in Redis database 1 for Supi %s: %v", smcontext.Supi, err)
			return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	*/
	// Return a success response
	// On success, the SMF shall return a "200 OK" with message body containing the representation of the ReleasedData when information needs to be returned to the NF Service Consumer, or a "204 No Content" response with an empty payload body in the POST response.
	return openapiserver.Response(http.StatusNoContent, nil), nil
}

// ReleaseUPFResources is a dummy function to release UPF resources
func (s *SmInfo) ReleaseUPFResources(n4SessionID string) error {
	// Simulate releasing UPF resources
	fmt.Println("Releasing UPF resources for N4 Session ID:", n4SessionID)
	time.Sleep(2 * time.Second)
	return nil
}

// RetrieveSmContext - Retrieve SM Context
func (s *SmInfo) ProcessNsmfRetrieveSmContextRequest(smContextRef string, smContextRetrieveData openapiserver.SmContextRetrieveData) (openapiserver.ImplResponse, error) {
	klog.Info("Retrieve is working")
	if smContextRef == "" {
		return openapiserver.Response(http.StatusBadRequest, nil), errors.New("smContextRef must not be empty")
	}
	var data SessionContext

	// dbData, err := s.sessionDb.Get(s.ctx, smContextRef).Result()
	// s.changeDatabase(redisClient.SessionDb)
	// dbData, err := s.dbClient.Get(s.ctx, smContextRef).Result()
	sData, err := s.DbClient.Redis.Read(smContextRef, redisClient.SessionDb)

	json.Unmarshal([]byte(sData), &data)
	klog.Info(data)
	if err != nil {
		klog.Info(smContextRef)
		klog.Info(err)
		return openapiserver.Response(http.StatusNotFound, nil), errors.New("smContextRef not a valid input")
	}
	retrieveProcess.Inc()
	//send as retrieved data TODO
	return openapiserver.Response(http.StatusOK, openapiserver.SmContextRetrievedData{
		UeEpsPdnConnection: "",
	}), nil // return Response(http.StatusNotImplemented, nil), errors.New("RetrieveSmContext method not implemented")
}

func validateDataUpdate(jsonData openapiserver.SmContextUpdateData) error {

	peiValid, _ := regexp.MatchString("^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$", jsonData.Pei)
	if !peiValid {
		return errors.New("invalid pei")
	}
	amfidValid, _ := regexp.MatchString("^[A-Fa-f0-9]{6}$", jsonData.Guami.AmfId)
	gMccValid, _ := regexp.MatchString("^[0-9]{3}$", jsonData.Guami.PlmnId.Mcc)
	gMncValid, _ := regexp.MatchString("^[0-9]{2,3}$", jsonData.Guami.PlmnId.Mnc)
	guamiValid := amfidValid && gMccValid && gMncValid
	if !guamiValid {
		return errors.New("invalid guami")
	}
	snMccValid, _ := regexp.MatchString("^[0-9]{3}$", jsonData.ServingNetwork.Mcc)
	snMncValid, _ := regexp.MatchString("^[0-9]{2,3}$", jsonData.ServingNetwork.Mnc)
	servingNetworkValid := snMccValid && snMncValid
	if !servingNetworkValid {
		return errors.New("invalid servingNetwork")
	}
	return nil

}

// UpdateSmContext - Update SM Context
func (s *SmInfo) ProcessNsmfUpdateSmContextRequest(smContextRef string, smContextUpdateData openapiserver.SmContextUpdateData, binaryDataN1SmMessage []byte, binaryDataN2SmInformation []byte, binaryDataN2SmInformationExt1 []byte) (openapiserver.ImplResponse, error) {

	klog.Info("Update is working")
	// Validate input parameters
	// Returns a `400 Bad Request` response if the smContextRef is empty.
	if smContextRef == "" {
		return openapiserver.Response(http.StatusBadRequest, openapiserver.SmContextUpdateError{
			Error: openapiserver.ProblemDetails{
				Title:  "Context Reference is empty",
				Type:   "ValidityErr",
				Detail: "smContextRef must not be empty",
				Status: 400,
				Cause:  "MANDATORY_QUERY_PARMS_MISSING",
			},
			N1SmMsg:      smContextUpdateData.N1SmMsg,
			N2SmInfo:     smContextUpdateData.N2SmInfo,
			N2SmInfoType: smContextUpdateData.N2SmInfoType,
		}), errors.New("smContextRef must not be empty")
	}

	// Validate the smContextUpdateData
	err := validateDataUpdate(smContextUpdateData)
	if err != nil {
		// Handle error from data validation
		return openapiserver.Response(http.StatusBadRequest, openapiserver.SmContextUpdateError{
			Error: openapiserver.ProblemDetails{
				Title:  "Invalid data sent",
				Type:   "ValidityErr",
				Detail: err.Error(),
				Status: 400,
				Cause:  "",
			},
			N1SmMsg:      smContextUpdateData.N1SmMsg,
			N2SmInfo:     smContextUpdateData.N2SmInfo,
			N2SmInfoType: smContextUpdateData.N2SmInfoType,
		}), err
	}

	// Check if the SM context exists
	var sessionData SessionContext

	// dbData, err := s.sessionDb.Get(s.ctx, smContextRef).Result()
	// s.changeDatabase(redisClient.SessionDb)
	// dbData, err := s.dbClient.Get(s.ctx, smContextRef).Result()

	sData, err := s.DbClient.Redis.Read(smContextRef, redisClient.SessionDb)

	// klog.Info(sData, err)
	json.Unmarshal([]byte(sData), &sessionData)
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "failed to get data in Redis: redis: nil" {
			return openapiserver.Response(http.StatusNotFound, openapiserver.SmContextUpdateError{
				Error: openapiserver.ProblemDetails{
					Title:  "Session Context Not found",
					Type:   "NotFoundErr",
					Detail: "Session context corresponding to smContextRef not found",
					Status: 404,
					Cause:  "CONTEXT_NOT_FOUND",
				},
				N1SmMsg:      smContextUpdateData.N1SmMsg,
				N2SmInfo:     smContextUpdateData.N2SmInfo,
				N2SmInfoType: smContextUpdateData.N2SmInfoType,
			}), errors.New("session context corresponding to smContextRef not found")
		} else {
			return openapiserver.Response(http.StatusNotFound, openapiserver.SmContextUpdateError{
				Error: openapiserver.ProblemDetails{
					Title:  "Redis Database Error",
					Type:   "InternalSernerError",
					Detail: err.Error(),
					Status: 404,
					Cause:  "Database_Error",
				},
				N1SmMsg:      smContextUpdateData.N1SmMsg,
				N2SmInfo:     smContextUpdateData.N2SmInfo,
				N2SmInfoType: smContextUpdateData.N2SmInfoType,
			}), err
		}
	}

	/*
		if sessionData.State == openapiserver.ACTIVATING {
			return openapiserver.Response(403, openapiserver.SmContextUpdateError{
				Error: openapiserver.ProblemDetails{
					Title:  "Already in Progress",
					Type:   "AlreadyInProgress",
					Detail: "Session already in Progress",
					Status: 403,
				},
				N1SmMsg:      smContextUpdateData.N1SmMsg,
				N2SmInfo:     smContextUpdateData.N2SmInfo,
				N2SmInfoType: smContextUpdateData.N2SmInfoType,
				UpCnxState:   sessionData.State,
			}), errors.New("request already in progress")
		}
	*/ // Update the SM context
	// Updates the SM context in the
	//TODO check if field exists then only change

	//fill pduSessionResourceSetup... after decode
	updateData := protos.SmContextUpdateDataRequest{
		ServingNfId:  sessionData.ServingNfId,
		Pei:          sessionData.Pei,
		PduSessionId: sessionData.PduSessionId,
		Guami: &protos.Guami{
			PlmnId: &protos.PlmnId{
				Mcc: sessionData.Guami.PlmnId.Mcc,
				Mnc: sessionData.Guami.PlmnId.Mnc,
			},
			AmfId: sessionData.Guami.AmfId,
		},
		SmContextStatusUri: sessionData.SmContextStatusUri,
		ServingNetwork: &protos.PlmnId{
			Mcc: sessionData.ServingNetwork.Mcc,
			Mnc: sessionData.ServingNetwork.Mnc,
		},
	}
	(*s.grpc).SendSmContextUpdateData(&updateData)

	anyN1Message, err := anypb.New(&pb.ByteDataWrapper{ByteArray: binaryDataN1SmMessage})
	if err != nil {
		klog.Fatalf("Failed to create Any message: %v", err)
	}

	reqType := "" // Empty string for the request type

	clientN1, ctxN1 := grpcnas.CreateGRPCNasClient()

	// Send the request
	respN1, err := clientN1.SendSMData(ctxN1, &pb.SMDataRequest{
		NasMessage: anyN1Message,
		TypeReq:    reqType,
	})
	if err != nil {
		klog.Fatalf("Error sending SMData: %v", err)
	}

	klog.Infof("Decoded binaryDataN1SmMessage: %v", respN1.NasResponse.Value)

	var nasResp pb.PDUSModReqModel

	err = json.Unmarshal(respN1.NasResponse.Value, &nasResp)
	if err != nil {
		klog.Fatalf("Error unmarshalling NAS decoded Message: %v", err)
	}

	// anyN2Message, err := anypb.New(&pb.ByteDataWrapper{ByteArray: binaryDataN2SmInformation})
	// if err != nil {
	// 	klog.Fatalf("Failed to create Any message: %v", err)
	// }
	// clientN2, ctxN2 := grpcnas.CreateGRPCNasClient()

	// // Send the request
	// respN2, err := clientN2.SendSMData(ctxN2, &pb.SMDataRequest{
	// 	NasMessage: anyN2Message,
	// 	TypeReq:    reqType,
	// })
	// if err != nil {
	// 	klog.Fatalf("Error sending SMData: %v", err)
	// }
	// klog.Infof("Decoded binaryDataN2SmInformation: %v", respN2.NasResponse.Value)

	// anyN2Ext1Message, err := anypb.New(&pb.ByteDataWrapper{ByteArray: binaryDataN2SmInformationExt1})
	// if err != nil {
	// 	klog.Fatalf("Failed to create Any message: %v", err)
	// }
	// clientN2Ext1, ctxN2Ext1 := grpcnas.CreateGRPCNasClient()

	// // Send the request
	// respN2Ext1, err := clientN2Ext1.SendSMData(ctxN2Ext1, &pb.SMDataRequest{
	// 	NasMessage: anyN2Ext1Message,
	// 	TypeReq:    reqType,
	// })
	// if err != nil {
	// 	klog.Fatalf("Error sending SMData: %v", err)
	// }
	// klog.Infof("Decoded binaryDataN2SmInformation: %v", respN2Ext1.NasResponse.Value)

	session := SessionContext{
		Pei:                 smContextUpdateData.Pei,
		ServingNfId:         smContextUpdateData.ServingNfId,
		ServingNetwork:      smContextUpdateData.ServingNetwork,
		Supi:                sessionData.Supi,
		Gpsi:                sessionData.Gpsi,
		ServiceName:         sessionData.ServiceName,
		RatType:             smContextUpdateData.RatType,
		AnType:              smContextUpdateData.AnType,
		UnauthenticatedSupi: sessionData.UnauthenticatedSupi,
		// PduSessionId: sessionData.PduSessionId,
		FDnn:               sessionData.FDnn,
		Guami:              smContextUpdateData.Guami,
		SmContextStatusUri: smContextUpdateData.SmContextStatusUri,
		NASepd:             nasResp.Epd,
		NASpduSessionId:    int(*nasResp.PdusessionID),
		NASpti:             int(*nasResp.Pti),
		NASmsgType:         nasResp.MsgType,
	}

	// Update data in Redis database 0

	// err = s.sessionDb.Set(s.ctx, smContextRef, SessionContextToJSON(session), 0).Err()
	// s.changeDatabase(redisClient.SessionDb)
	// err = s.dbClient.Set(s.ctx, smContextRef, SessionContextToJSON(session), 0).Err()

	_, err = s.DbClient.Redis.Create(
		smContextRef,
		SessionContextToJSON(session),
		redisClient.SessionDb,
	)

	if err != nil {
		klog.Errorf("Error updating data in Redis database 0 for smContextRef %s: %v", smContextRef, err)
		//TODO proper response as per openapi
		return openapiserver.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	updateProcess.Inc()
	// Return a success response
	// Returns a `200 OK` response.
	// sessionContext := s.smcontextTable[smContextRef]
	return openapiserver.Response(http.StatusOK, openapiserver.SmContextUpdatedData{
		UpCnxState:   session.State,
		N1SmMsg:      smContextUpdateData.N1SmMsg,
		N2SmInfo:     smContextUpdateData.N2SmInfo,
		N2SmInfoType: smContextUpdateData.N2SmInfoType,
	}), nil

}

func (s *SmInfo) ProcessN1N2Message(grpcMsg *protos.N1N2MessageTransferDataRequest, ip string) error {

	reqType := "PDU_SESSION_ESTABLISHMENT_ACCEPT"

	n1msgByteArray, n1msgClass, pduSessionId, err := buildNASmsg(s, reqType, grpcMsg.SmContextID)
	if err != nil {
		return err
	}

	/*hd := ""
	if hd, err = os.UserHomeDir(); err != nil {
		return err
	}

	err = os.WriteFile(hd+"/wipro5gc/n1n2data/n1MsgContent", n1msgByteArray, os.ModeAppend)
	if err != nil {
		return err
	}*/

	n2InfoClass := "SM"
	n2Info := &protos.N2Information{
		GtpTunnelEndpointIpAddress: grpcMsg.TunnelEndpointAdr,
		GtpTeid:                    grpcMsg.Gtpteid,
		//QoS info hardcoded for now
		Qfi: 1,
		QosFlowLevelQosParameters: &protos.QosFlowLevelQosParameters{
			Fqi: 1,
			Arp: &protos.Arp{
				PriorityLevel:           "10",
				PreemptionCapability:    "NOT_PREEMPT",
				PreemptionVulnerability: "NOT_PREEMPTABLE",
			},
			GbrQosFlowInformation: &protos.GbrQoSFlowInformation{
				MaximumFlowBitrateDownlink:    100,
				MaximumFlowBitrateUplink:      100,
				GuaranteedFlowBitrateDownlink: 100,
				GuaranteedFlowBitrateUplink:   100,
			},
		},
	}

	// get encoded n2 info from ngap
	encodedn2Info := (*s.grpc).SendN2Info(n2Info)
	if encodedn2Info.Error != "nil" {
		return errors.New(encodedn2Info.Error)
	}
	n2InfoByteArray := encodedn2Info.EncodedData

	// construct n1n2 message struct
	n1n2Message := openapi_commn_client.N1N2MessageTransferReqData{
		N1MessageContainer: &openapi_commn_client.N1MessageContainer{
			N1MessageClass: openapi_commn_client.N1MessageClass{
				String: &n1msgClass,
			},
			N1MessageContent: openapi_commn_client.RefToBinaryData{
				ContentId: "binaryDataN1MessageContent",
			},
			//nfid and serviceinstanceid not given
		},
		N2InfoContainer: &openapi_commn_client.N2InfoContainer{
			N2InformationClass: openapi_commn_client.N2InformationClass{
				String: &n2InfoClass,
			},
			SmInfo: &openapi_commn_client.N2SmInformation{
				//harcoded pdusessionid
				PduSessionId: pduSessionId,
				N2InfoContent: &openapi_commn_client.N2InfoContent{
					NgapIeType: &openapi_commn_client.NgapIeType{
						String: &grpcMsg.NgapleType,
					},
					NgapData: openapi_commn_client.RefToBinaryData{
						ContentId: "binaryDataN2InfoContent",
					},
				},
			},
		},
		PduSessionId: &pduSessionId,
		OldGuami: &openapi_commn_client.Guami{
			PlmnId: openapi_commn_client.PlmnId{
				Mcc: "286",
				Mnc: "01",
			},
			AmfId: "218A9E",
		},
	}

	//hardcoded for now
	supi := "imsi-128600000001"
	s.apiClient.N1N2MessageTransfer(ip, supi, n1n2Message, n1msgByteArray, n2InfoByteArray)

	return nil
}

func buildNASmsg(s *SmInfo, reqType string, smContextID string) ([]byte, string, int32, error) {
	var sessionData SessionContext
	//Read the data from redis
	sData, err := s.DbClient.Redis.Read(smContextID, redisClient.SessionDb)
	if err != nil {
		return nil, "", -1, err
	}
	json.Unmarshal([]byte(sData), &sessionData)

	//Hard Coded the following values for time being
	comp := []string{"Match-all type"}
	pf := nas.PacketFilter{Identifier: 15, Direction: "BIDIRECTIONAL", Components: comp}
	pflist := []nas.PacketFilter{pf}
	qos := nas.QoSRule{QoSIdentifier: "QRI 1", Operation: "Create new QoS rule", DQR: "DEFAULT_QoS_RULE", Precedence: 255, Segregation: "Segregation not requested", QFI: "QFI 1", PacketFilterList: pflist}
	ses := nas.SessionAMBR{IEI: 1, UnitUL: "MULT_1Kbps", RateUL: 30, UnitDL: "MULT_1Kbps", RateDL: 40}
	qosList := []nas.QoSRule{qos, qos}
	//Build the message
	establishmentAccept := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: sessionData.NASepd, PDUsessionId: sessionData.NASpduSessionId, PTI: sessionData.NASpti, MessageType: reqType, PduSessionType: "IPV4", SSCmode: "SSC_MODE_1", AuthorizedQoSRules: qosList, SessionAmbr: ses}

	//Proto Equivalents for above data
	pbPf := pb.PacketFilters{Identifier: uint32(pf.Identifier), Direction: pf.Direction, Components: pf.Components}

	pbQos := pb.QosRules{Qosidentifier: qos.QoSIdentifier, Operation: qos.Operation, Dqr: qos.DQR, Precidence: uint32(qos.Precedence), Seg: qos.Segregation, Qfi: qos.QFI, Pf: []*pb.PacketFilters{&pbPf}}
	pbQosList := []*pb.QosRules{&pbQos}

	pbSes := pb.Sessionambr{Iei: int32(ses.IEI), UnitUL: ses.UnitUL, RateUL: int32(ses.RateUL), UnitDL: ses.UnitDL, RateDL: int32(ses.RateDL)}

	esPb := pb.PDUSEstAccModel{Epd: establishmentAccept.ExtendedProtocolDiscriminator, PdusessionID: int32(establishmentAccept.PDUsessionId), Pti: int32(establishmentAccept.PTI), MsgType: establishmentAccept.MessageType, PdusType: establishmentAccept.PduSessionType, SscMode: establishmentAccept.SSCmode, QosIEI: int32(establishmentAccept.QosRuleIEI), Qosrules: pbQosList, Sessionambr: &pbSes}
	anyEstAccMsg, err := anypb.New(&esPb)
	if err != nil {
		klog.Fatalf("Failed to create Any message: %v", err)
		return nil, "", int32(sessionData.NASpduSessionId), err
	}

	reqTypeEstAcc := reqType

	// Send the request
	client, ctx := grpcnas.CreateGRPCNasClient()
	respEnc, err := client.SendSMData(ctx, &pb.SMDataRequest{
		NasMessage: anyEstAccMsg,
		TypeReq:    reqTypeEstAcc,
	})
	if err != nil {
		klog.Fatalf("Error sending SMData: %v", err)
		return nil, "", int32(sessionData.NASpduSessionId), err
	}

	n1msgresp := respEnc.NasResponse.Value[2:]

	return n1msgresp, "SM", int32(sessionData.NASpduSessionId), nil

}
