package sm

import (
	openapiserver "w5gc.io/wipro5gcore/openapi/openapiserver"
)

type upCnxSubState int

const (
	ACTIVATING_CREATE_IN_PROGRESS upCnxSubState = 1
	ACTIVATING_CREATED            upCnxSubState = 2
	ACTIVATING_UPDATE_IN_PROGRESS upCnxSubState = 3
	ACTIVATED_RELEASE_IN_PROGRESS upCnxSubState = 4
)

type SessionContext struct {
	//change to upcnxstatus TODO
	State                              openapiserver.UpCnxState      `json:"state,omitempty"`
	SubState                           upCnxSubState                 `json:"substate,omitempty"`
	Supi                               string                        `json:"supi,omitempty"`
	UnauthenticatedSupi                bool                          `json:"unauthenticatedSupi,omitempty"`
	Pei                                string                        `json:"pei,omitempty"`
	Gpsi                               string                        `json:"gpsi,omitempty"`
	PduSessionId                       int32                         `json:"pduSessionId,omitempty"`
	FDnn                               string                        `json:"dnn,omitempty"`
	ServingNfId                        string                        `json:"servingNfId"`
	Guami                              openapiserver.Guami           `json:"guami,omitempty"`
	ServiceName                        openapiserver.ServiceName     `json:"serviceName,omitempty"`
	ServingNetwork                     openapiserver.PlmnId          `json:"servingNetwork"`
	RequestType                        openapiserver.RequestType     `json:"requestType,omitempty"`
	N1SmMsg                            openapiserver.RefToBinaryData `json:"n1SmMsg,omitempty"`
	AnType                             openapiserver.AccessType      `json:"anType"`
	RatType                            openapiserver.RatType         `json:"ratType,omitempty"`
	SmContextStatusUri                 string                        `json:"smContextStatusUri"`
	N4SessionID                        string
	UeContextId                        string
	ContextRefId                       string
	NASepd                             string
	NASpduSessionId                    int
	NASpti                             int
	NASmsgType                         string
	NASmaxIntegrityProtectedDataRateUL string
	NASmaxIntegrityProtectedDataRateDL string
	NASSmCause                         string
}

// uecontext newname
type UserContext struct {
	NoSession int
}
