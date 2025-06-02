package sm

type SessionManager interface {
	Start()
}

type SMContextMessage *interface{}

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
)
