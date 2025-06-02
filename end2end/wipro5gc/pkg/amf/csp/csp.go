package csp

import (
	//"encoding/binary"
	//"net"
	//"reflect"
	//"sync"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"k8s.io/klog"

	//"github.com/gin-gonic/gin"
	"github.com/benbjohnson/clock"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"w5gc.io/wipro5gcore/pkg/amf/csp/apiclient"
	"w5gc.io/wipro5gcore/pkg/amf/csp/apiserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/grpcserver"
	"w5gc.io/wipro5gcore/pkg/amf/csp/grpc/protos/ngapNas"

	"w5gc.io/wipro5gcore/pkg/amf/csp/config"
	"w5gc.io/wipro5gcore/pkg/amf/csp/sm/sessions"
	"w5gc.io/wipro5gcore/utils/cache"

	"w5gc.io/wipro5gcore/pkg/amf/csp/sm"
)

const (
	retransmitInterval = time.Second * 3
	retransmitRetries  = 3
	backoffInterval    = time.Second * 3
	// Period for performing global cleanup tasks.
	housekeepingPeriod = time.Second * 2
	metricsPort        = ":2113"
)

// CspHandler is an interface implemented for testability
type CspHandler interface {
	//InitiateSessioReportRequest(sessionMsg *pfcpUdp.Message)
	HandleSessionCleanups() error
}

// Bootstrap is a bootstrapping interface for PDU SMS
type CspBootstrap interface {
	//GetConfiguration()
	//GetContext()
	Run(configChannel <-chan config.CspConfig)
}

type Csp struct {
	config         *config.CspConfig
	sessionManager sm.SessionManager
	sessionCache   cache.WorkCache
	//sessionDB                     db.SessionDB
	grpc            grpc.Grpc
	apiClient       apiclient.ApiClient
	apiServer       apiserver.ApiServer
	sessionWorkers  SessionWorkers
	clock           clock.Clock
	backoffInterval time.Duration
	timerT1         time.Duration
	retriesN1       uint8
	context         *CspContext
}

type CspContext struct {
	startTime   time.Time
	lastRestart time.Time
	restarts    int // If number of restarts in last 10 sec > 3 reset TODO GURU
	//lock                  sync.Mutex
	NodeId string
}

func (p *Csp) runClient(grpcMsg *grpcserver.GrpcMessage) {
	if grpcMsg != (&grpcserver.GrpcMessage{}) {
		switch grpcMsg.MsgType {
		case "create":
			p.apiClient.Start(grpcMsg)
		case "update":
			p.apiClient.Start(grpcMsg)
		case "release":
			p.apiClient.Start(grpcMsg)
		case "retrieve":
			p.apiClient.Start(grpcMsg)
		}

	}
}

// Initialize Csp
func NewCsp(cfg *config.CspConfig, time time.Time) (CspBootstrap, bool) {

	csp := &Csp{
		config:          cfg,
		clock:           clock.New(),
		backoffInterval: backoffInterval,
		timerT1:         retransmitInterval,
		retriesN1:       retransmitRetries,
		context: &CspContext{
			startTime: time,
		},
	}

	// Intialize the api handler
	csp.apiClient = apiclient.NewApiClient(cfg)
	csp.apiServer = apiserver.NewApiServer(cfg.NodeInfo)

	// Initialize session manager
	//csp.sessionManager = sm.NewSessionManager(csp.config.NodeInfo, csp.config.n11Nodes, time, csp.backoffInterval, csp.timerT1, csp.retriesN1)

	// Intialize session cache
	csp.sessionCache = cache.NewCache(csp.clock)

	// Intialize session db
	// csp.sessionDb = db.NewDB()

	// Intialize session workers
	csp.sessionWorkers = NewSessionWorkers(csp.handleSession, csp.sessionCache, csp.backoffInterval)

	// Intialize grpc
	csp.grpc = grpc.NewGrpc(*cfg)

	// resyncInterval, backOffPeriod TODO GURU

	return csp, true
}

// Run starts the Csp
func (p *Csp) Run(configChannel <-chan config.CspConfig) {
	g := []byte{46, 1, 1, 193, 255, 255, 145, 161, 40, 1, 0, 123, 0, 7, 128, 0, 10, 0, 0, 13, 0, 18, 1, 129, 34, 4, 1, 0, 0, 1, 37, 9, 8, 105, 110, 116, 101, 114, 110, 101, 116}
	hd, _ := os.UserHomeDir()
	e := os.WriteFile(hd+"/wipro5gc/testdata/PDUsessionEstNasPdu", g, os.ModeAppend)
	if e != nil {
		klog.Info("cannot write n1SmContainer to file")
	}
	go p.apiServer.Start(&(p.grpc))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("Metrics server panicked: %v", r)
			}
		}()
		p.startMetricsServer()
	}()
	// start the session manager
	// p.sessionManager.Start()

	// Start the api handler
	//TODO accept input variable for type of request here
	// RequestType = make(chan string)
	fmt.Println("BEFORE RUN CLIENT FUNCTION IN CSP")
	x := &grpcserver.GrpcMessage{}
	p.runClient(x)
	fmt.Println("AFTER RUN CLIENT FUNCTION IN CSP")

	// Start the grpc
	klog.Info("before csp grpc function")
	p.grpc.Start()

	// Start the db Hanlder

	// Start the csp event handler
	//commented temporarily
	klog.Info("before csp events function")
	p.cspEvents(configChannel, p)
	klog.Info("after csp events function")
	//testing locally
	/*cspNssai := create_sm_context_grpc.Snssai{Sst: 0, Sd: "0"}
	cspPlmnID := create_sm_context_grpc.PlmnId{Mcc: "1", Mnc: "2"}
	cspTai := create_sm_context_grpc.Tai{PlmnId: &cspPlmnID, Tac: "3"}
	cspNcgi := create_sm_context_grpc.Ncgi{PlmnId: cspTai.PlmnId, NrCellId: "4"}

	cspNrLoc := create_sm_context_grpc.NrLocation{Tai: &cspTai, Ncgi: &cspNcgi, AgeOfLocationInformation: 1, UeLocationTimestamp: &timestamppb.Timestamp{Seconds: 4}, GeographicalInformation: "india", GeodeticInformation: "43E10N", GlobalGnbId: 2}
	cspMsg := create_sm_context_grpc.CreateSmContextDataFromNasMod{AnType: "7", Snssai: &cspNssai, NrLocation: &cspNrLoc, PduSessionId: 0, N1SmContainer: []byte{10, 2, 1}}
	client := grpcclient.NewGrpcClient()
	client.Start()
	client.SendDataForCreateSmContext(&cspMsg)*/
	// return
}

// Run starts the Csp
/*func (p *Csp) Run(configChannel <-chan config.CspConfig, requestType string) {
	// start the session manager
	// p.sessionManager.Start()

	// Start the api handler
	//TODO accept input variable for type of request here
	switch requestType {
	case "create":
		p.apiClient.Start(1)
	case "update":
		p.apiClient.Start(3)
	case "release":
		p.apiClient.Start(5)
	case "retrieve":
		p.apiClient.Start(7)
	}

	p.apiServer.Start()

	// Start the grpc
	p.grpc.Start()

	// Start the db Hanlder

	// Start the csp event handler
	p.cspEvents(configChannel, p)

	return
}*/

func (p *Csp) cspEvents(configChannel <-chan config.CspConfig, handler CspHandler) {
	syncTicker := time.NewTicker(time.Second)
	defer syncTicker.Stop()
	housekeepingTicker := time.NewTicker(housekeepingPeriod)
	defer housekeepingTicker.Stop()
	sessionChannel := p.apiClient.WatchApiChannel()
	grpcChannel := p.grpc.WatchGrpcChannel()
	n1n2Channel := p.grpc.WatchN1N2DataGrpcChannel()
	p.handleCspEvents(configChannel, sessionChannel, grpcChannel, n1n2Channel, syncTicker.C, housekeepingTicker.C, handler)
}

// handleCspEvents is the main loop for processing events in csp
func (p *Csp) handleCspEvents(configChannel <-chan config.CspConfig, sessionChannel <-chan *apiclient.SessionMessage,
	grpcChannel <-chan *grpcserver.GrpcMessage, n1n2Channel <-chan *ngapNas.DLRequest, syncCh <-chan time.Time, housekeepingCh <-chan time.Time,
	handler CspHandler) bool {
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
			klog.Info("received data from nasmm in grpc channel")
			klog.Infof("%+v", *grpcMsg.GrpcMsg)
			//hardcoded message type for now
			//read a parameter from grpcMsg when other message types are implemented
			grpcMsg.MsgType = "create"
			p.runClient(grpcMsg)

		case n1n2Data := <-n1n2Channel:
			if n1n2Data != (&ngapNas.DLRequest{}) {
				// n1sm := n1n2Data.GetN1DataBytes()
				// n2 := n1n2Data.GetN2DataBytes()
				// send n1n2data to nasmm, nasmm will construct nasmm layer, and forward new n1n2data to ngap
				err := p.grpc.SendN1N2Data(n1n2Data)
				klog.Info("Error in sending n1n2 data bytes : ", err.Error)
			}

		// Handle csp node events, sesssion events/ notificatiosn TODO GURU
		//case event := <-cspEventChannel:
		// Event for a session.
		/*if session, ok := p.sessionManager.GetSession(event.SessionID); ok {
		                  klog.V(2).Infof("handleCspEvents (EVENT): %q, event: %#v", format.Sessions(session), event)
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

/*func (p *Csp) getSessionsToSync() {
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

func (p *Csp) HandleSessionCleanups() error {
	deletedSessions := make(map[sessions.SessionId]struct{})
	err := p.sessionWorkers.RemoveSessionWorkers(deletedSessions)
	return err
}

// dispatchWork handles the session in a session worker
func (p *Csp) dispatchWork(sessionId sessions.SessionId, sessionMsg sm.SMContextMessage, startTime time.Time) {
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
func (p *Csp) handleSession(msgInfo SessionMessageInfo) error {
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

func (p *Csp) startMetricsServer() {
	log.Println("Starting metrics CSP server at", metricsPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Metrics handler registered.")
	if err := http.ListenAndServe(metricsPort, nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
