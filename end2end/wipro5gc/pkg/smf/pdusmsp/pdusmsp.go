package pdusmsp

import (
	//"encoding/binary"
	//"net"

	//"sync"

	"fmt"
	"strconv"
	"time"

	//	"fmt"
	"k8s.io/klog"
	//"github.com/gin-gonic/gin"
	"github.com/benbjohnson/clock"

	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/api"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/apiclient"
	db "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/database"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc/protos"

	//"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/sm/nodes"

	openapiserver "w5gc.io/wipro5gcore/openapi/openapiserver"

	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/sm/sessions"
	"w5gc.io/wipro5gcore/utils/cache"

	// "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc"
	"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/sm"
)

const (
	retransmitInterval = time.Second * 3
	retransmitRetries  = 3
	backoffInterval    = time.Second * 3
	// Period for performing global cleanup tasks.
	housekeepingPeriod = time.Second * 2
)

// PdusmspHandler is an interface implemented for testability
type PdusmspHandler interface {
	//InitiateSessioReportRequest(sessionMsg *pfcpUdp.Message)
	HandleSessionCleanups() error
}

// Bootstrap is a bootstrapping interface for PDU SMS
type PdusmspBootstrap interface {
	//GetConfiguration()
	//GetContext()
	Run(configChannel <-chan config.PdusmspConfig)
}

type Pdusmsp struct {
	config          *config.PdusmspConfig
	sessionManager  sm.SessionManager
	sessionCache    cache.WorkCache
	dbManager       db.DBManager
	grpc            grpc.Grpc
	apiClient       apiclient.ApiClient
	apiServer       api.ApiServer
	sessionWorkers  SessionWorkers
	clock           clock.Clock
	backoffInterval time.Duration
	timerT1         time.Duration
	retriesN1       uint8
	context         *PdusmspContext
}

type PdusmspContext struct {
	startTime   time.Time
	lastRestart time.Time
	restarts    int // If number of restarts in last 10 sec > 3 reset TODO GURU
	//lock                  sync.Mutex
	NodeId string
}

// Initialize Pdusmsp
func NewPdusmsp(cfg *config.PdusmspConfig, time time.Time) (PdusmspBootstrap, bool) {

	pdusmsp := &Pdusmsp{
		config:          cfg,
		clock:           clock.New(),
		backoffInterval: backoffInterval,
		timerT1:         retransmitInterval,
		retriesN1:       retransmitRetries,
		context: &PdusmspContext{
			startTime: time,
		},
	}

	// Intialize the api handler
	//should both apiClient and apiServer be created in newcsc?
	pdusmsp.apiClient = apiclient.NewApiClient(cfg)
	pdusmsp.apiServer = api.NewApiServer(cfg.NodeInfo)

	// Intialize session db
	pdusmsp.dbManager = db.NewDBManager()
	//temperory to be romoved
	dbInfo := pdusmsp.dbManager.(*db.DBInfo)

	// Initialize session manager
	//pdusmsp.sessionManager = sm.NewSessionManager(pdusmsp.config.NodeInfo, pdusmsp.config.n11Nodes, time, pdusmsp.backoffInterval, pdusmsp.timerT1, pdusmsp.retriesN1)

	// Intialize grpc
	pdusmsp.grpc = grpc.NewGrpc(cfg.GrpcServerInfo, cfg.GrpcClientInfo)
	// s:=(pdusmsp.apiClient).(apiclient.ApiClientInfo)
	// var s sm.SessionManager
	pdusmsp.sessionManager = sm.NewSessionManager(dbInfo, &pdusmsp.grpc, pdusmsp.apiClient)
	//	s:= pdusmsp.sessionManager.NewSMContextAPIService()
	//	fmt.Println(s)
	// Intialize session cache
	pdusmsp.sessionCache = cache.NewCache(pdusmsp.clock)

	// Intialize session workers
	pdusmsp.sessionWorkers = NewSessionWorkers(pdusmsp.handleSession, pdusmsp.sessionCache, pdusmsp.backoffInterval)

	// resyncInterval, backOffPeriod TODO GURU

	return pdusmsp, true
}

// Run starts the Pdusmsp
func (p *Pdusmsp) Run(configChannel <-chan config.PdusmspConfig) {
	// start the session manager
	go p.sessionManager.Start()

	// Start the api handler
	p.apiClient.Start()
	go p.apiServer.Start()

	// Start the grpc
	go p.grpc.Start()

	// Start the db Hanlder
	p.dbManager.Start()

	// Start the pdusmsp event handler
	p.pdusmspEvents(configChannel, p)

	// return
}

func (p *Pdusmsp) pdusmspEvents(configChannel <-chan config.PdusmspConfig, handler PdusmspHandler) {
	klog.Infof("Entered into pdusmspEvents")
	syncTicker := time.NewTicker(time.Second)
	defer syncTicker.Stop()
	housekeepingTicker := time.NewTicker(housekeepingPeriod)
	defer housekeepingTicker.Stop()
	sessionChannel := p.apiServer.WatchApiChannel()
	grpcChannel := p.grpc.WatchGrpcChannel()
	p.handlePdusmspEvents(configChannel, sessionChannel, grpcChannel, syncTicker.C, housekeepingTicker.C, handler)
}

// handlePdusmspEvents is the main loop for processing events in pdusmsp
func (p *Pdusmsp) handlePdusmspEvents(configChannel <-chan config.PdusmspConfig, sessionChannel <-chan *api.SessionMessage,
	grpcChannel <-chan *grpcserver.GrpcMessage, syncCh <-chan time.Time, housekeepingCh <-chan time.Time,
	handler PdusmspHandler) bool {
	klog.Info("Entered into handlePdusmspEvents")
	for {
		select {
		// Handle config updates of nodes - TODO GURU
		//case config := <-configChannel:
		//switch config.Entity {
		//case CPNODES:
		// Handle config updates for CP Nodes i.e. SMF nodes
		//switch config.Type
		//case ADDNODE
		//case UPNODES:
		// Handle config updates for UP Nodes i.e. UPFU nodes
		//}
		case grpcMsg := <-grpcChannel:
			grpcMsgType := sm.MessageType(grpcMsg.MsgType)
			switch grpcMsgType {
			case sm.NSMF_N1_N2_TRANSFER:
				klog.Info("handlePdusmspEvents (N1N2Message Transfer)")
				klog.Infof("%+v", *grpcMsg.GrpcMsg)
				// //TODO ask raghu and verify data we are getting
				data := *grpcMsg.GrpcMsg
				// //TODO raghu's datatype to be used here
				trData := data.(*protos.N1N2MessageTransferDataRequest)

				// refId := strconv.Itoa(int(trData.PduSessionId)) + trData.OldGuami.AmfId
				refId := fmt.Sprintf("%v", trData.SmContextID)
				sessionId := sessions.SessionId(refId)
				p.dispatchWork(sessions.SessionId(sessionId),
					nil,
					grpcMsg,
					grpcMsgType,
					time.Now(),
				)
			}
		case pdusmsMsg := <-sessionChannel:
			switch pdusmsMsg.MsgType {
			case sm.NSMF_CREATE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Create SM Context Request
				klog.Infof("handlePdusmspEvents (CREATE SM CONTEXT REQUEST)")
				//put proper session
				jsonData := pdusmsMsg.SessionMsg
				smData := jsonData.(openapiserver.SmContextCreateData)
				refContext := strconv.Itoa(int(smData.PduSessionId)) + smData.Guami.AmfId
				sessionId := sessions.SessionId(refContext)
				// p.dispatchWork(sessionId, pdusmsMsg, grpcMsg, time.Now())
				p.dispatchWork(
					sessionId,
					pdusmsMsg,
					nil,
					pdusmsMsg.MsgType,
					time.Now())

			case sm.NSMF_UPDATE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Update SM Context Request
				klog.Infof("handlePdusmspEvents (UPDATE SM CONTEXT  REQUEST)")
				// fmt.Println(reflect.TypeOf(pdusmsMsg.SessionMsg))
				//have to change in future
				// id, err := strconv.Atoi(pdusmsMsg.SmContextRefID)
				// if err != nil {
				// 	klog.Error(err.Error())
				// }
				// klog.Info(id)
				sessionId := sessions.SessionId(pdusmsMsg.SmContextRefID)
				// p.dispatchWork(sessionId, pdusmsMsg, grpcMsg, time.Now())
				p.dispatchWork(
					sessionId,
					pdusmsMsg,
					nil,
					pdusmsMsg.MsgType,
					time.Now())

			case sm.NSMF_RELEASE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Release SM Context Request
				klog.Infof("handlePdusmspEvents (RELEASE SM CONTEXT REQUEST)")
				// id, err := strconv.Atoi(pdusmsMsg.SmContextRefID)
				// if err != nil {
				// 	klog.Error(err.Error())
				// }
				sessionId := sessions.SessionId(pdusmsMsg.SmContextRefID)
				// p.dispatchWork(sessionId, pdusmsMsg, grpcMsg, time.Now())
				p.dispatchWork(
					sessionId,
					pdusmsMsg,
					nil,
					pdusmsMsg.MsgType,
					time.Now())

			case sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Retrieve SM Context Request
				klog.Infof("handlePdusmspEvents (RETRIEVE SM CONTEXT REQUEST)")
				// id, err := strconv.Atoi(pdusmsMsg.SmContextRefID)
				// if err != nil {
				// 	klog.Error(err.Error())
				// }
				sessionId := sessions.SessionId(pdusmsMsg.SmContextRefID)
				// p.dispatchWork(sessionId, pdusmsMsg, grpcMsg, time.Now())
				p.dispatchWork(
					sessionId,
					pdusmsMsg,
					nil,
					pdusmsMsg.MsgType,
					time.Now())

			}

		// Handle pdusmsp node events, sesssion events/ notificatiosn TODO GURU
		//case event := <-pdusmspEventChannel:
		// Event for a session.
		/*if session, ok := p.sessionManager.GetSession(event.SessionID); ok {
		                  klog.V(2).Infof("handlePdusmspEvents (EVENT): %q, event: %#v", format.Sessions(session), event)
		                  //handler.HandleSessionEvents()
		          } else {
		                  // If the session no longer exists, ignore the event.
		                  klog.V(4).Infof("handleUpfcEvents (EVENT): ignore irrelevant event: %#v", e)
		          }
		  }

		  if event.Type == sessionEvent.SessionDisconnected {

		  }*/
		//case config := <-ConfigChan:
		// Check node and send to corresponding channel

		// Handle session sync - TODO GURU
		// To handle asyncs during reboot , configuration error etc
		case <-syncCh:
		// Sync sessions waiting for sync

		/*sessionsToSync := p.getSessionsToSync()
		  if len(sessionsToSync) == 0 {
		          break
		  }
		  klog.V(4).Infof("SyncLoop (SYNC): %d sessions", len(sessionsToSync))
		  //handler.HandleSessionSyncs(sessionsToSync)*/

		// Handle house keeping of sessions TODO GURU
		case <-housekeepingCh:
			klog.V(4).Infof("SyncLoop (housekeeping)")
			if err := handler.HandleSessionCleanups(); err != nil {
				klog.Errorf("Failed cleaning session: %v", err)
			}

		}
	}
}

/*func (p *Pdusmsp) getSessionsToSync() {
        allSessions := p.n4Manager.pfcp.sessions
        sessionIds := p.workCache.GetItem()

        var sessionsToSync []*sm.Session
        for _, session := range allSessions {
                if session.sessionId in sessionIds {
                        sessionsToSync = append(sessionsToSync, session)
                        continue
                }
        }
        return sessionsToSync
}*/

func (p *Pdusmsp) HandleSessionCleanups() error {
	deletedSessions := make(map[sessions.SessionId]struct{})
	err := p.sessionWorkers.RemoveSessionWorkers(deletedSessions)
	return err
}

// dispatchWork handles the session in a session worker
func (p *Pdusmsp) dispatchWork(sessionId sessions.SessionId,
	sessionMsg *api.SessionMessage,
	grpcMsg *grpcserver.GrpcMessage,
	msgType sm.MessageType,
	startTime time.Time) {

	if msgType == sm.NSMF_N1_N2_TRANSFER {
		p.sessionWorkers.HandleSessionMessages(&SessionMessageInfo{
			SessionId: sessionId,
			StartTime: startTime,
			PdusmsMsg: api.SessionMessage{},
			GrpcMsg:   *grpcMsg,
			MsgType:   msgType,
			OnCompleteFunc: func(err error) {
				// Handle on completion of session update
				if err == nil {
					klog.Infof("Successfully handled pdusms  for session %s", sessionId)
					//metrics.SessionWorkerDuration.WithLabelValues(syncType.String()).Observe(metrics.SinceInSeconds(start))
				} else {
					// Log error and update cause with request rejected
					klog.Info(err)
					klog.Errorf("Unable to handle pdusms for session %s", sessionId)
				}
			},
		})
		return
	}
	// Run the handle session in an async worker.
	p.sessionWorkers.HandleSessionMessages(&SessionMessageInfo{
		SessionId: sessionId,
		StartTime: startTime,
		PdusmsMsg: *sessionMsg,
		GrpcMsg:   grpcserver.GrpcMessage{},
		MsgType:   msgType,
		OnCompleteFunc: func(err error) {
			// Handle on completion of session update
			if err == nil {
				klog.Infof("Successfully handled pdusms  for session %s", sessionId)
				//metrics.SessionWorkerDuration.WithLabelValues(syncType.String()).Observe(metrics.SinceInSeconds(start))
			} else {
				// Log error and update cause with request rejected
				klog.Info(err)
				klog.Errorf("Unable to handle pdusms for session %s", sessionId)
			}
		},
	})
	// Monitor the session and number of qos flows for the session. TODO GURU
	//metrics.QoSFlowsPerSessionCount.Observe(float64(len(session.flows)))
}

// handleSession handles the pdu session message in a session worker
func (p *Pdusmsp) handleSession(msgInfo SessionMessageInfo) error {
	// Get the required message info
	klog.Info("inside")
	msgType := msgInfo.MsgType
	pdusmsMsg := msgInfo.PdusmsMsg
	// grpc->myfunction(process/dbupdate)->callruchi'sfunction
	// Process remote node requests/responses
	switch msgType {
	case sm.NSMF_CREATE_SM_CONTEXT_REQUEST:
		// Process pdusms create request
		//sm.ProcessNsmfCreateSmContextRequest(pdusmsMsg)
		n1SmMessage := pdusmsMsg.BinaryDataN1SmMessage
		jsonData := pdusmsMsg.SessionMsg
		smData := jsonData.(openapiserver.SmContextCreateData)
		resp, err := p.sessionManager.ProcessNsmfCreateSmContextRequest(smData, n1SmMessage)
		chanRec := p.apiServer.WatchRecChannel()
		chanRec <- &api.Receiver{RecievedResponse: resp, RecievedErr: err}
		return err
	case sm.NSMF_UPDATE_SM_CONTEXT_REQUEST:
		// Process pdusms update request
		n1SmMessage := pdusmsMsg.BinaryDataN1SmMessage
		n2SmMessage := pdusmsMsg.BinaryDataN2SmInformation
		n2SmExt1Message := pdusmsMsg.BinaryDataN2SmInformationExt1
		smcontextref := pdusmsMsg.SmContextRefID
		jsonData := pdusmsMsg.SessionMsg
		smData := jsonData.(openapiserver.SmContextUpdateData)
		resp, err := p.sessionManager.ProcessNsmfUpdateSmContextRequest(smcontextref, smData, n1SmMessage, n2SmMessage, n2SmExt1Message)
		//sm.ProcessNsmfUpdateSmContextRequest(pdusmsMsg)
		chanRec := p.apiServer.WatchRecChannel()
		chanRec <- &api.Receiver{RecievedResponse: resp, RecievedErr: err}
		return err
	case sm.NSMF_RELEASE_SM_CONTEXT_REQUEST:
		// Process pdusms release request
		n2SmMessage := pdusmsMsg.BinaryDataN2SmInformation
		smcontextref := pdusmsMsg.SmContextRefID
		jsonData := pdusmsMsg.SessionMsg
		smData := jsonData.(openapiserver.SmContextReleaseData)
		resp, err := p.sessionManager.ProcessNsmfReleaseSmContextRequest(smcontextref, smData, n2SmMessage)
		//sm.ProcessNsmfReleaseSmContextRequest(pdusmsMsg)
		chanRec := p.apiServer.WatchRecChannel()
		chanRec <- &api.Receiver{RecievedResponse: resp, RecievedErr: err}
		return err
	case sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST:
		// Process pdusms retrieve request
		smcontextref := pdusmsMsg.SmContextRefID
		jsonData := pdusmsMsg.SessionMsg
		smData := jsonData.(openapiserver.SmContextRetrieveData)
		resp, err := p.sessionManager.ProcessNsmfRetrieveSmContextRequest(smcontextref, smData)
		//sm.ProcessNsmfRetrieveSmContextRequest(pdusmsMsg)
		chanRec := p.apiServer.WatchRecChannel()
		chanRec <- &api.Receiver{RecievedResponse: resp, RecievedErr: err}
		return err
	case sm.NSMF_N1_N2_TRANSFER:
		recGrpcMsg := (*msgInfo.GrpcMsg.GrpcMsg).(*protos.N1N2MessageTransferDataRequest)
		//hardcoded for now
		ip := "csp-service.csp.svc.cluster.local" //change to amf IP based on certain conditions // to ask Guru
		p.sessionManager.ProcessN1N2Message(recGrpcMsg, ip)
	}
	return nil
}
