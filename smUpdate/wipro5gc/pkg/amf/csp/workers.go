package csp

import (
	"sync"
	"time"

	"w5gc.io/wipro5gcore/pkg/amf/csp/sm"
	"w5gc.io/wipro5gcore/pkg/amf/csp/sm/nodes"
	"w5gc.io/wipro5gcore/pkg/amf/csp/sm/sessions"
	"w5gc.io/wipro5gcore/utils/cache"
)

type SessionId = sessions.SessionId
type NodeType = nodes.NodeType
type OnCompleteFunction func(err error)

// SessionWorkers is an abstract interface for testability.
type SessionWorkers interface {
	HandleSessionMessages(msgInfo *SessionMessageInfo)
	RemoveSessionWorkers(map[SessionId]struct{}) error
}

type SessionFunction func(msgInfo SessionMessageInfo) error

type SessionMessageInfo struct {
	SessionId      SessionId          // Session Id of the message handled
	StartTime      time.Time          // Start of message handling
	OnCompleteFunc OnCompleteFunction // callback function when operation completes
	MsgType        sm.MessageType
	PdusmsMsg      sm.SMContextMessage // Message
}

type CspSessionWorkers struct {
	sessionLock   sync.Mutex                            // Protects all per worker fields.
	sessionWorker map[SessionId]chan SessionMessageInfo // Per-session worker to process messages received through
	// its corresponding channel.
	isMessageHandled       map[SessionId]bool               // Per-session goroutines state
	lastUndeliveredMessage map[SessionId]SessionMessageInfo // Last undelivered message for this session
	sessionFunction        SessionFunction                  // FUnction handling the session messages
	workCache              cache.WorkCache                  // Cache to store the unhandled messsages
	backoffInterval        time.Duration                    // Interval after which to retry the handling the message
}

// Initialize the session wprkers
func NewSessionWorkers(function SessionFunction, cache cache.WorkCache, backoffInterval time.Duration) SessionWorkers {
	return &CspSessionWorkers{
		sessionWorker:          make(map[SessionId]chan SessionMessageInfo),
		isMessageHandled:       make(map[SessionId]bool),
		lastUndeliveredMessage: make(map[SessionId]SessionMessageInfo),
		sessionFunction:        function,
		workCache:              cache,
		backoffInterval:        backoffInterval,
	}
}

// HandleSessionMessages start a session worker for each session to handle its messages
func (p *CspSessionWorkers) HandleSessionMessages(msgInfo *SessionMessageInfo) {

	var sessionWorker chan SessionMessageInfo
	var exists bool

	sessionId := msgInfo.SessionId

	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()

	// Check session worker exists
	if sessionWorker, exists = p.sessionWorker[sessionId]; !exists {
		// Create worker to handle the messages
		sessionWorker = make(chan SessionMessageInfo, 1)
		p.sessionWorker[sessionId] = sessionWorker

		// Create a new worker for new session or when amf-upfgw restarted.
		go func() {
			//defer ??
			p.sessionWorkerHandler(sessionWorker)
		}()
	}

	// Continue handling the session message if no other other messsage is currently handled
	if !p.isMessageHandled[sessionId] {
		// Set message is handled and provide the message info to the worker
		p.isMessageHandled[sessionId] = true
		sessionWorker <- *msgInfo
	} else {
		// Drop the message?
		// Do we need to increase the channel buffer?? TODO GURU
		p.lastUndeliveredMessage[sessionId] = *msgInfo
	}
}

// sessionWorkerHandler starts a new worker for a new session to handle its messages
func (p *CspSessionWorkers) sessionWorkerHandler(sessionMessages <-chan SessionMessageInfo) {

	// Handle the session messages in a continuos loop
	for msgInfo := range sessionMessages {
		sessionId := msgInfo.SessionId
		err := func() error {

			err := p.sessionFunction(SessionMessageInfo{
				SessionId: sessionId,
				StartTime: msgInfo.StartTime,
				PdusmsMsg: msgInfo.PdusmsMsg,
			})
			return err
		}()

		// Notify the call-back function on handling the message
		if msgInfo.OnCompleteFunc != nil {
			msgInfo.OnCompleteFunc(err)
		}

		// On error Add the work item again to the work cache
		if err != nil {
			p.workCache.AddItem(sessionId, p.backoffInterval)
		}
		if msgInfo, exists := p.lastUndeliveredMessage[sessionId]; exists {
			p.sessionWorker[sessionId] <- msgInfo
			delete(p.lastUndeliveredMessage, sessionId)
		} else {
			p.isMessageHandled[sessionId] = false
		}
	}
}

func (p *CspSessionWorkers) RemoveSessionWorkers(deletedSessions map[SessionId]struct{}) error {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	for sessionId, _ := range p.sessionWorker {
		if _, exists := deletedSessions[sessionId]; !exists {
			if ch, ok := p.sessionWorker[sessionId]; ok {
				close(ch)
				delete(p.sessionWorker, sessionId)
				delete(p.lastUndeliveredMessage, sessionId)
			}
		}
	}
	return nil
}
