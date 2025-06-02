package pdusmsp

import (
	//"encoding/binary"
	//"net"
	//"reflect"
	//"sync"
	"time"

	"k8s.io/klog"
	//"github.com/gin-gonic/gin"
	"github.com/benbjohnson/clock"

	"w5gc.io/wipro5gcore/stubs/upfgw/api"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc"
	"w5gc.io/wipro5gcore/stubs/upfgw/grpc/grpcserver"

	//"w5gc.io/wipro5gcore/stubs/upfgw/sm/nodes"
	"w5gc.io/wipro5gcore/stubs/upfgw/config"
	"w5gc.io/wipro5gcore/stubs/upfgw/sm/sessions"
	"w5gc.io/wipro5gcore/utils/cache"

	//"w5gc.io/wipro5gcore/pkg/util/db"
	"w5gc.io/wipro5gcore/stubs/upfgw/sm"
	//"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/grpc"
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
	config         *config.PdusmspConfig
	sessionManager sm.SessionManager
	sessionCache   cache.WorkCache
	//sessionDB                     db.SessionDB
	grpc            grpc.Grpc
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
	pdusmsp.apiServer = api.NewApiServer(cfg.NodeInfo)

	// Initialize session manager
	//pdusmsp.sessionManager = sm.NewSessionManager(pdusmsp.config.NodeInfo, pdusmsp.config.n11Nodes, time, pdusmsp.backoffInterval, pdusmsp.timerT1, pdusmsp.retriesN1)

	// Intialize session cache
	pdusmsp.sessionCache = cache.NewCache(pdusmsp.clock)

	// Intialize session db
	// pdusmsp.sessionDb = db.NewDB()

	// Intialize session workers
	pdusmsp.sessionWorkers = NewSessionWorkers(pdusmsp.handleSession, pdusmsp.sessionCache, pdusmsp.backoffInterval)

	// Intialize grpc
	pdusmsp.grpc = grpc.NewGrpc()

	// resyncInterval, backOffPeriod TODO GURU

	return pdusmsp, true
}

// Run starts the Pdusmsp
func (p *Pdusmsp) Run(configChannel <-chan config.PdusmspConfig) {
	// start the session manager
	// p.sessionManager.Start()

	// Start the api handler
	go p.grpc.Start()
	p.apiServer.Start()

	// Start the grpc

	// Start the db Hanlder

	// Start the pdusmsp event handler
	p.pdusmspEvents(configChannel, p)

	return
}

func (p *Pdusmsp) pdusmspEvents(configChannel <-chan config.PdusmspConfig, handler PdusmspHandler) {
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
		case pdusmsMsg := <-sessionChannel:
			switch pdusmsMsg.MsgType {
			case sm.NSMF_CREATE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Create SM Context Request
				klog.Infof("handlePdusmspEvents (CREATE SM CONTEXT REQUEST)")
				sessionId := sessions.SessionId(1)
				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

			case sm.NSMF_UPDATE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Update SM Context Request
				klog.Infof("handlePdusmspEvents (UPDATE SM CONTEXT  REQUEST)")
				sessionId := sessions.SessionId(1)
				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

			case sm.NSMF_RELEASE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Release SM Context Request
				klog.Infof("handlePdusmspEvents (RELEASE SM CONTEXT REQUEST)")
				sessionId := sessions.SessionId(1)
				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

			case sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST:
				// PDU Session management service - Retrieve SM Context Request
				klog.Infof("handlePdusmspEvents (RETRIEVE SM CONTEXT REQUEST)")
				sessionId := sessions.SessionId(1)
				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

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

			// Handle metrics TODO GURU

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
func (p *Pdusmsp) dispatchWork(sessionId sessions.SessionId, sessionMsg sm.SMContextMessage, startTime time.Time) {
	// Run the handle session in an async worker.
	p.sessionWorkers.HandleSessionMessages(&SessionMessageInfo{
		SessionId: sessionId,
		StartTime: startTime,
		PdusmsMsg: sessionMsg,
		OnCompleteFunc: func(err error) {
			// Handle on completion of session update
			if err == nil {
				klog.Infof("Successfully handled pdusms  for session %d", sessionId)
				//metrics.SessionWorkerDuration.WithLabelValues(syncType.String()).Observe(metrics.SinceInSeconds(start))
			} else {
				// Log error and update cause with request rejected
				klog.Errorf("Unable to handle pdusms for session %d", sessionId)
			}
		},
	})
	// Monitor the session and number of qos flows for the session. TODO GURU
	//metrics.QoSFlowsPerSessionCount.Observe(float64(len(session.flows)))
}

// handleSession handles the pdu session message in a session worker
func (p *Pdusmsp) handleSession(msgInfo SessionMessageInfo) error {
	// Get the required message info
	//pdusmsMsg := msgInfo.PdusmsMsg
	msgType := msgInfo.MsgType

	// Process remote node requests/responses
	switch msgType {
	case sm.NSMF_CREATE_SM_CONTEXT_REQUEST:
		// Process pdusms creaye request
		//sm.ProcessNsmfCreateSmContextRequest(pdusmsMsg)
	case sm.NSMF_UPDATE_SM_CONTEXT_REQUEST:
		// Process pdusms update request
		//sm.ProcessNsmfUpdateSmContextRequest(pdusmsMsg)
	case sm.NSMF_RELEASE_SM_CONTEXT_REQUEST:
		// Process pdusms release request
		//sm.ProcessNsmfReleaseSmContextRequest(pdusmsMsg)
	case sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST:
		// Process pdusms retrieve request
		//sm.ProcessNsmfRetrieveSmContextRequest(pdusmsMsg)
	}
	return nil
}
