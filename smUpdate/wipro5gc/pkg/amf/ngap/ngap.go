package ngap

import (
	//"encoding/binary"
	//"net"
	//"reflect"
	//"sync"
	"net"
	"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	"github.com/benbjohnson/clock"

	//"w5gc.io/wipro5gcore/pkg/amf/csp/apiclient"
	//"w5gc.io/wipro5gcore/pkg/amf/csp/apiserver"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/message"

	//"w5gc.io/wipro5gcore/pkg/amf/csp/sm/nodes"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
	//"w5gc.io/wipro5gcore/pkg/amf/csp/sm/sessions"
	//"w5gc.io/wipro5gcore/utils/cache"
	//"w5gc.io/wipro5gcore/pkg/util/db"
	//"w5gc.io/wipro5gcore/pkg/amf/csp/sm"
	//"w5gc.io/wipro5gcore/pkg/amf/csp/grpc"

	"github.com/ishidawataru/sctp"
	logger "w5gc.io/wipro5gcore/pkg/amf/ngap/log"
)

const (
	retransmitInterval = time.Second * 3
	retransmitRetries  = 3
	backoffInterval    = time.Second * 3
	// Period for performing global cleanup tasks.
	// housekeepingPeriod = time.Second * 2
)

// CspHandler is an interface implemented for testability
/*type CspHandler interface {
	//InitiateSessioReportRequest(sessionMsg *pfcpUdp.Message)
	HandleSessionCleanups() error
}*/

// Bootstrap is a bootstrapping interface for NGAP
type NgapBootstrap interface {
	//GetConfiguration()
	//GetContext()
	Run()
}

// type NgapServer interface{
// 	Start()
// }

type Ngap struct {
	config *config.NgapConfig
	//sessionManager sm.SessionManager
	// ngapServer NgapServer
	//sessionCache   cache.WorkCache
	//sessionDB                     db.SessionDB
	grpc grpc.Grpc
	//apiClient       apiclient.ApiClient
	//apiServer       apiserver.ApiServer
	//sessionWorkers  SessionWorkers
	clock           clock.Clock
	backoffInterval time.Duration
	timerT1         time.Duration
	retriesN1       uint8
	context         *NgapContext
	// ueDb		db.UeDb
}

type NgapContext struct {
	startTime   time.Time
	lastRestart time.Time
	restarts    int // If number of restarts in last 10 sec > 3 reset TODO GURU
	//lock                  sync.Mutex
	NodeId string
}

// Initialize Ngap
func NewNgap(cfg *config.NgapConfig, time time.Time) (NgapBootstrap, bool) {

	ngap := &Ngap{
		config:          cfg,
		clock:           clock.New(),
		backoffInterval: backoffInterval,
		timerT1:         retransmitInterval,
		retriesN1:       retransmitRetries,
		context: &NgapContext{
			startTime: time,
		},
	}

	// Intialize the api handler
	// csp.apiClient = apiclient.NewApiClient(cfg)
	// csp.apiServer = apiserver.NewApiServer(cfg.NodeInfo)

	// Initialize session manager
	//csp.sessionManager = sm.NewSessionManager(csp.config.NodeInfo, csp.config.n11Nodes, time, csp.backoffInterval, csp.timerT1, csp.retriesN1)

	// Intialize session cache
	// csp.sessionCache = cache.NewCache(csp.clock)

	// Intialize session db
	// csp.sessionDb = db.NewDB()

	// Intialize session workers
	// csp.sessionWorkers = NewSessionWorkers(csp.handleSession, csp.sessionCache, csp.backoffInterval)

	// Initialize db
	// ngap.ueDb = db.NewDBManager()
	// Intialize grpc
	//commented for now
	ngap.grpc = grpc.NewGrpc(ngap.config)

	// resyncInterval, backOffPeriod TODO GURU

	return ngap, true
}

// Run starts the Ngap
// func (p *Ngap) Run(configChannel <-chan config.NgapConfig) {
func (p *Ngap) Run() {
	// fmt.Println("inside run function")
	cfg := p.config.NodeInfo
	// fmt.Println(cfg.ApiPort)
	// fmt.Println(cfg.NodeId)
	// start the session manager
	// p.sessionManager.Start()
	// p.ngapServer.Start()

	// Start the api handler
	//TODO accept input variable for type of request here
	// switch requestType {
	// case "create":
	// 	p.apiClient.Start(1)
	// case "update":
	// 	p.apiClient.Start(3)
	// case "release":
	// 	p.apiClient.Start(5)
	// case "retrieve":
	// 	p.apiClient.Start(7)
	// }
	// p.apiClient.Start()

	// p.apiServer.Start()

	// Start the grpc
	p.grpc.Start()

	// Start the db Hanlder

	// Start the csp event handler
	// p.cspEvents(configChannel, p)
	//commented go routine for now
	//go func(cfg config.AmfNodeInfo) {
	// fmt.Println("go func started")
	nodeIpStr := cfg.NodeId
	ngapPortStr := cfg.ApiPort
	ngapPort, _ := strconv.Atoi(ngapPortStr)
	ips := []net.IPAddr{}
	netAddr, err := net.ResolveIPAddr("ip", nodeIpStr)
	if err != nil {
		logger.AppLog.Errorf("Error in resolving address %s : %v", nodeIpStr, err)
		// e := string("Error in resolving address %s : %v", addr, err)
		//return err
	} else {
		logger.AppLog.Infof("Resolved address %s to %s", nodeIpStr, netAddr)
		ips = append(ips, *netAddr)
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ngapPort,
	}
	err = listenAndServe(addr, &(p.grpc))
	if err != nil {
		logger.AppLog.Errorf("Error in listening on given address %s : %v", nodeIpStr, err)
	}

}

/*func Run(addresses []string, port int) error {
	ips := []net.IPAddr{}

	for _, addr := range addresses {
		netAddr, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			logger.AppLog.Errorf("Error in resolving address %s : %v", addr, err)
			// e := string("Error in resolving address %s : %v", addr, err)
			return err
		} else {
			logger.AppLog.Infof("Resolved address %s to %s", addr, netAddr)
			ips = append(ips, *netAddr)
		}
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    port,
	}

	err := listenAndServe(addr)
	if err != nil {
		return err
	}
	return nil
}*/

const rcvbuf int = 26144

var Connection net.Conn

func listenAndServe(addr *sctp.SCTPAddr, grpc *grpc.Grpc) error {
	listener, err := sctp.ListenSCTP("sctp", addr)
	if err != nil {
		logger.AppLog.Fatalf("failed to listen: %v", err)
		return err
	}
	logger.AppLog.Printf("App listening on %s", listener.Addr())
	for {
		conn, err := listener.Accept()
		Connection = conn
		if err != nil {
			logger.AppLog.Fatalf("failed to accept: %v", err)
			return err
		}
		logger.AppLog.Printf("Accepted Connection from RemoteAddr: %s", conn.RemoteAddr())
		wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn.(*sctp.SCTPConn))
		if rcvbuf != 0 {
			err = wconn.SetReadBuffer(rcvbuf)
			if err != nil {
				logger.AppLog.Fatalf("failed to set read buf: %v", err)
				return err
			}
		}
		go handleConnection(conn, rcvbuf, grpc)
	}
}

func handleConnection(conn net.Conn, bufsize int, grpc *grpc.Grpc) {
	for {
		buf := make([]byte, bufsize+128)
		_, err := conn.Read(buf)
		if err != nil {
			logger.AppLog.Errorf("Read failed : %v", err)
			return
		}
		//handle message now
		message.HandleMessage(conn, buf, grpc)
	}
}

// func (p *Ngap) ngapEvents(configChannel <-chan config.NgapConfig, handler csp.cspHandler) {
// 	syncTicker := time.NewTicker(time.Second)
// 	defer syncTicker.Stop()
// 	housekeepingTicker := time.NewTicker(housekeepingPeriod)
// 	defer housekeepingTicker.Stop()
// 	sessionChannel := p.apiClient.WatchApiChannel()
// 	grpcChannel := p.grpc.WatchGrpcChannel()
// 	p.handleCspEvents(configChannel, sessionChannel, grpcChannel, syncTicker.C, housekeepingTicker.C, handler)
// }

// handleCspEvents is the main loop for processing events in csp
// func (p *Csp) handleCspEvents(configChannel <-chan config.CspConfig, sessionChannel <-chan *apiclient.SessionMessage,
// 	grpcChannel <-chan *grpc.GrpcMessage, syncCh <-chan time.Time, housekeepingCh <-chan time.Time,
// 	handler CspHandler) bool {
// 	for {
// 		select {
// 		// Handle config updates of nodes - TODO GURU
// 		//case config := <-configChannel:
// 		//switch config.Entity {
// 		//case CPNODES:
// 		// Handle config updates for CP Nodes i.e. SMF nodes
// 		//switch config.Type
// 		//case ADDNODE
// 		//case UPNODES:
// 		// Handle config updates for UP Nodes i.e. UPFU nodes
// 		//}
// 		case pdusmsMsg := <-sessionChannel:
// 			switch pdusmsMsg.MsgType {
// 			case sm.NSMF_CREATE_SM_CONTEXT_REQUEST:
// 				// PDU Session management service - Create SM Context Request
// 				klog.Infof("handleCspEvents (CREATE SM CONTEXT REQUEST)")
// 				sessionId := sessions.SessionId(1)
// 				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

// 			case sm.NSMF_UPDATE_SM_CONTEXT_REQUEST:
// 				// PDU Session management service - Update SM Context Request
// 				klog.Infof("handleCspEvents (UPDATE SM CONTEXT  REQUEST)")
// 				sessionId := sessions.SessionId(1)
// 				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

// 			case sm.NSMF_RELEASE_SM_CONTEXT_REQUEST:
// 				// PDU Session management service - Release SM Context Request
// 				klog.Infof("handleCspEvents (RELEASE SM CONTEXT REQUEST)")
// 				sessionId := sessions.SessionId(1)
// 				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

// 			case sm.NSMF_RETRIEVE_SM_CONTEXT_REQUEST:
// 				// PDU Session management service - Retrieve SM Context Request
// 				klog.Infof("handleCspEvents (RETRIEVE SM CONTEXT REQUEST)")
// 				sessionId := sessions.SessionId(1)
// 				p.dispatchWork(sessionId, pdusmsMsg.SessionMsg, time.Now())

// 			}

// 		// Handle csp node events, sesssion events/ notificatiosn TODO GURU
// 		//case event := <-cspEventChannel:
// 		// Event for a session.
// 		/*if session, ok := p.sessionManager.GetSession(event.SessionID); ok {
// 		                  klog.V(2).Infof("handleCspEvents (EVENT): %q, event: %#v", format.Sessions(session), event)
// 		                  //handler.HandleSessionEvents()
// 		          } else {
// 		                  // If the session no longer exists, ignore the event.
// 		                  klog.V(4).Infof("handleUpfcEvents (EVENT): ignore irrelevant event: %#v", e)
// 		          }
// 		  }

// 		  if event.Type == sessionEvent.SessionDisconnected {

// 		  }*/
// 		//case config := <-ConfigChan:
// 		// Check node and send to corresponding channel

// 		// Handle session sync - TODO GURU
// 		// To handle asyncs during reboot , configuration error etc
// 		case <-syncCh:
// 		// Sync sessions waiting for sync

// 		/*sessionsToSync := p.getSessionsToSync()
// 		  if len(sessionsToSync) == 0 {
// 		          break
// 		  }
// 		  klog.V(4).Infof("SyncLoop (SYNC): %d sessions", len(sessionsToSync))
// 		  //handler.HandleSessionSyncs(sessionsToSync)*/

// 		// Handle house keeping of sessions TODO GURU
// 		case <-housekeepingCh:
// 			klog.V(4).Infof("SyncLoop (housekeeping)")
// 			if err := handler.HandleSessionCleanups(); err != nil {
// 				klog.Errorf("Failed cleaning session: %v", err)
// 			}

// 			// Handle metrics TODO GURU

// 		}
// 	}
// }

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

/*func (p *Csp) HandleSessionCleanups() error {
	deletedSessions := make(map[sessions.SessionId]struct{})
	err := p.sessionWorkers.RemoveSessionWorkers(deletedSessions)
	return err
}*/

// dispatchWork handles the session in a session worker
/*func (p *Csp) dispatchWork(sessionId sessions.SessionId, sessionMsg sm.SMContextMessage, startTime time.Time) {
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
}*/

// handleSession handles the pdu session message in a session worker
/*func (p *Csp) handleSession(msgInfo SessionMessageInfo) error {
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
}*/
