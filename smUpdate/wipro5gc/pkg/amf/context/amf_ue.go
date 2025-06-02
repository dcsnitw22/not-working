package context

import (
	"sync"

	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/nas"
)

type OnGoingProcedure string

const (
	OnGoingProcedureNothing      OnGoingProcedure = "Nothing"
	OnGoingProcedurePaging       OnGoingProcedure = "Paging"
	OnGoingProcedureN2Handover   OnGoingProcedure = "N2Handover"
	OnGoingProcedureRegistration OnGoingProcedure = "Registration"
)

type StateType string

// State is a thread-safe structure that represents the state of a object
// and it can be used for FSM
type State struct {
	// current state of the State object
	current StateType
	// stateMutex ensures that all operations to current is thread-safe
	stateMutex sync.RWMutex
}

// GMM state for UE
const (
	Deregistered            StateType = "Deregistered"
	DeregistrationInitiated StateType = "DeregistrationInitiated"
	Authentication          StateType = "Authentication"
	SecurityMode            StateType = "SecurityMode"
	ContextSetup            StateType = "ContextSetup"
	Registered              StateType = "Registered"
)

type AmfUe struct {
	/* the AMF which serving this AmfUe now */
	servingAMF *AMFContext // never nil

	/* Gmm State */
	State map[openapi_commn_client.AccessType]*State
	/* Registration procedure related context */
	RegistrationType5GS             uint8
	IdentityTypeUsedForRegistration uint8
	RegistrationRequest             *nas.RegistrationRequestModel
	// ServingAmfChanged                  bool
	// DeregistrationTargetAccessType     uint8 // only used when deregistration procedure is initialized by the network
	RegistrationAcceptForNon3GPPAccess []byte
	NasPduValue                        []byte
	RetransmissionOfInitialNASMsg      bool
	RequestIdentityType                uint8
	/* Used for AMF relocation */
	// TargetAmfProfile *models.NfProfile
	// TargetAmfUri     string
	/* Ue Identity */
	PlmnId openapi_commn_client.PlmnId
	// Suci                   string
	Supi                string
	UnauthenticatedSupi bool
	// Gpsi                   string
	// Pei                    string
	// Tmsi                   int32 // 5G-Tmsi
	// Guti                   string
	// GroupID                string
	// EBI                    int32
	// EventSubscriptionsInfo map[string]*AmfUeEventSubscription
	/* User Location */
	RatType  openapi_commn_client.RatType
	Location openapi_commn_client.UserLocation
	Tai      openapi_commn_client.Tai
	// LocationChanged          bool
	LastVisitedRegisteredTai openapi_commn_client.Tai
	TimeZone                 string // "[+-]HH:MM[+][1-2]", Refer to TS 29.571 - 5.2.2 Simple Data Types
	/* context about udm */
	/*	UdmId                             string
		NudmUECMUri                       string
		NudmSDMUri                        string
		ContextValid                      bool
		Reachability                      models.UeReachability
		SubscribedData                    models.SubscribedData
		SmfSelectionData                  *models.SmfSelectionSubscriptionData
		UeContextInSmfData                *models.UeContextInSmfData
		TraceData                         *models.TraceData
		UdmGroupId                        string
		SubscribedNssai                   []models.SubscribedSnssai
		AccessAndMobilitySubscriptionData *models.AccessAndMobilitySubscriptionData
		BackupAmfInfo                     []models.BackupAmfInfo*/
	/* contex abut ausf */
	/*	AusfGroupId                       string
		AusfId                            string
		AusfUri                           string
		RoutingIndicator                  string
		AuthenticationCtx                 *models.UeAuthenticationCtx
		AuthFailureCauseSynchFailureTimes int
		IdentityRequestSendTimes          int
		ABBA                              []uint8
		Kseaf                             string
		Kamf                              string*/
	/* context about PCF */
	/*	PcfId                        string
		PcfUri                       string
		PolicyAssociationId          string
		AmPolicyUri                  string
		AmPolicyAssociation          *models.PolicyAssociation
		RequestTriggerLocationChange bool // true if AmPolicyAssociation.Trigger contains RequestTrigger_LOC_CH*/
	/* UeContextForHandover */
	// HandoverNotifyUri string
	/* N1N2Message */
	N1N2MessageIDGenerator          *IDGenerator
	N1N2Message                     *N1N2Message
	N1N2MessageSubscribeIDGenerator *IDGenerator
	// map[int64]models.UeN1N2InfoSubscriptionCreateData; use n1n2MessageSubscriptionID as key
	N1N2MessageSubscription sync.Map
	/* Pdu Sesseion context */
	SmContextList sync.Map // map[int32]*SmContext, pdu session id as key
	/* Related Context */
	RanUe map[openapi_commn_client.AccessType]*RanUe
	/* other */
	// onGoing                         map[models.AccessType]*OnGoing
	UeRadioCapability string // OCTET string
	// Capability5GMM    nasType.Capability5GMM
	// ConfigurationUpdateIndication   nasType.ConfigurationUpdateIndication
	// ConfigurationUpdateCommandFlags *ConfigurationUpdateCommandFlags
	/* context related to Paging */
	/*	UeRadioCapabilityForPaging                 *UERadioCapabilityForPaging
		InfoOnRecommendedCellsAndRanNodesForPaging *InfoOnRecommendedCellsAndRanNodesForPaging
		UESpecificDRX                              uint8*/
	/* Security Context */
	/*SecurityContextAvailable bool
	UESecurityCapability     nasType.UESecurityCapability // for security command
	NgKsi                    models.NgKsi
	MacFailed                bool      // set to true if the integrity check of current NAS message is failed
	KnasInt                  [16]uint8 // 16 byte
	KnasEnc                  [16]uint8 // 16 byte
	Kgnb                     []uint8   // 32 byte
	Kn3iwf                   []uint8   // 32 byte
	NH                       []uint8   // 32 byte
	NCC                      uint8     // 0..7
	ULCount                  security.Count
	DLCount                  security.Count
	CipheringAlg             uint8
	IntegrityAlg             uint8*/
	/* Registration Area */
	RegistrationArea map[openapi_commn_client.AccessType][]openapi_commn_client.Tai
	// LadnInfo         []factory.Ladn
	/* Network Slicing related context and Nssf */
	NssfId  string
	NssfUri string
	// NetworkSliceInfo                  *models.AuthorizedNetworkSliceInfo
	AllowedNssai                      map[openapi_commn_client.AccessType][]openapi_commn_client.AllowedSnssai
	ConfiguredNssai                   []openapi_commn_client.ConfiguredSnssai
	NetworkSlicingSubscriptionChanged bool
	SdmSubscriptionId                 string
	UeCmRegistered                    map[openapi_commn_client.AccessType]bool
	/* T3513(Paging) */
	// T3513 *Timer // for paging
	/* T3565(Notification) */
	// T3565 *Timer // for NAS Notification
	/* T3560 (for authentication request/security mode command retransmission) */
	// T3560 *Timer
	/* T3550 (for registration accept retransmission) */
	// T3550 *Timer
	/* T3522 (for deregistration request) */
	// T3522 *Timer
	/* T3570 (for identity request) */
	// T3570 *Timer
	/* T3555 (for configuration update command) */
	// T3555 *Timer
	/* Ue Context Release Cause */
	// ReleaseCause map[models.AccessType]*CauseAll
	/* T3502 (Assigned by AMF, and used by UE to initialize registration procedure) */
	T3502Value             int        // Second
	T3512Value             int        // default 54 min
	Non3gppDeregTimerValue int        // default 54 min
	Lock                   sync.Mutex // Update context to prevent race condition

	// logger
	/*NASLog      *logrus.Entry
	GmmLog      *logrus.Entry
	ProducerLog *logrus.Entry*/
}

type N1N2Message struct {
	Request     openapi_commn_client.N1N2MessageTransferReqData
	Status      openapi_commn_client.N1N2MessageTransferCause
	ResourceUri string
}
