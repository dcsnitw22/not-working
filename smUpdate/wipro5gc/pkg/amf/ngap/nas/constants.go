package nas

// Constants for Extended Protocol Discriminator values
const (
	MobilityManagementEPD byte = 0b01111110
	SessionManagementEPD  byte = 0b00101110
)

// Constants for MaximumDataRatePerUeForUserPlaneIntegrityProtectionForUplink and Downlink
const (
	Kbps         byte = 0b00000000
	FullDataRate byte = 0b11111111
)

// Constants for message type
const (
	RegistrationRequest              byte = 0b01000001
	RegistrationAccept               byte = 0b01000010
	RegistrationComplete             byte = 0b01000011
	RegistrationReject               byte = 0b01000100
	DeregistrationRequestUeOrigin    byte = 0b01000101
	DeregistrationAcceptUeOrigin     byte = 0b01000110
	DeregistrationRequestUeTerminate byte = 0b01000111
	DeregistrationAcceptUeTerminate  byte = 0b01001000
	ServiceRequest                   byte = 0b01001100
	ServiceReject                    byte = 0b01001101
	ServiceAccept                    byte = 0b01001110
	ConfigurationUpdateCommand       byte = 0b01010100
	ConfigurationUpdateComplete      byte = 0b01010101
	AuthenticationRequest            byte = 0b01010110
	AuthenticationResponse           byte = 0b01010111
	AuthenticationReject             byte = 0b01011000
	AuthenticationFailure            byte = 0b01011001
	AuthenticationResult             byte = 0b01011010
	IdentityRequest                  byte = 0b01011011
	IdentityResponse                 byte = 0b01011100
	SecurtiyModeCommand              byte = 0b01011101
	SecurtiyModeComplete             byte = 0b01011110
	SecurtiyModeReject               byte = 0b01011111
	FiveGMMStatus                    byte = 0b01100100
	Notification                     byte = 0b01100101
	NotificationResponse             byte = 0b01100110
	UlNasTransport                   byte = 0b01100111
	DlNasTransport                   byte = 0b01101000
	PduSEstablishmentRequest         byte = 0b11000001
	PduSEstablishmentAccept          byte = 0b11000010
	PduSEstablishmentReject          byte = 0b11000011
	PduSAuthenticationCommand        byte = 0b11000101
	PduSAutheticationComplete        byte = 0b11000110
	PduSAuthenticationResult         byte = 0b11000111
	PduSModificationRequest          byte = 0b11001001
	PduSModificationReject           byte = 0b11001010
	PduSModificationCommand          byte = 0b11001011
	PduSModificationComplete         byte = 0b11001100
	PduSModificationCommandReject    byte = 0b11001101
	PduSReleaseRequest               byte = 0b11010001
	PduSReleaseReject                byte = 0b11010010
	PduSReleaseCommand               byte = 0b11010011
	PduSReleaseComplete              byte = 0b11010100
	FiveGSMStatus                    byte = 0b11010110
)

// Constants for PDU Session ID
const (
	NoVal byte = 0b00000000
	Val1  byte = 0b00000001
	Val2  byte = 0b00000010
	Val3  byte = 0b00000011
	Val4  byte = 0b00000100
	Val5  byte = 0b00000101
	Val6  byte = 0b00000110
	Val7  byte = 0b00000111
	Val8  byte = 0b00001000
	Val9  byte = 0b00001001
	Val10 byte = 0b00001010
	Val11 byte = 0b00001011
	Val12 byte = 0b00001100
	Val13 byte = 0b00001101
	Val14 byte = 0b00001110
	Val15 byte = 0b00001111
)

// Constants for DQR
const (
	NotDefaultQoS byte = 0b0
	DefaultQoS    byte = 0b1
)

// Constants for Packet Filter Component Type Identifier
const (
	MatchAllType               byte = 0b00000001
	IPv4RemoteAddressType      byte = 0b00010000
	IPv4LocalAddressType       byte = 0b00010001
	IPv6RemoteAddress          byte = 0b00100001
	IPv6LocalAddress           byte = 0b00100011
	ProtocolIdentifier         byte = 0b00110000
	SingleLocalPortType        byte = 0b01000000
	LocalPortRangeType         byte = 0b01000001
	SingleRemotePortType       byte = 0b01010000
	RemotePortRangeType        byte = 0b01010001
	SecurityParameterIndexType byte = 0b01100000
	TypeOfService              byte = 0b01110000
	FlowLabelType              byte = 0b10000000
	DestinationMACaddressType  byte = 0b10000001
	SourceMACAddressType       byte = 0b10000010
	CTAGVIDtype                byte = 0b10000011
	STAGVIDtype                byte = 0b10000100
	CTAGPCPtype                byte = 0b10000101
	STAGPCPtype                byte = 0b10000110
	Ethertype                  byte = 0b10000111
	DestinationMAC             byte = 0b10001000
	SourceMAC                  byte = 0b10001001
)

// Constants for Packet Filter Direction
const (
	Uplink        byte = 0b10
	Downlink      byte = 0b01
	Bidirectional byte = 0b11
)

// Constants for PDU Session Type
const (
	IpV4         byte = 0b001
	IpV6         byte = 0b010
	IPV4V6       byte = 0b011
	Unstructured byte = 0b100
	Ethernet     byte = 0b101
)

// Constants for QoS Rule Operation Code
const (
	CreateNewQoSRule                                   byte = 0b001
	DeleteExistingQoSRule                              byte = 0b010
	ModifyExistingQoSRuleAndAddPacketFilters           byte = 0b011
	ModifyExistingQoSRuleAndReplaceAllPacketFilters    byte = 0b100
	ModifyExistingQoSRuleAndDeletePacketFilters        byte = 0b101
	ModifyExistingQoSRuleWithoutModifyingPacketFilters byte = 0b110
)

// Constants for Seggregation
const (
	Requested    byte = 0b1
	NotRequested byte = 0b0
)

// Constants for SSC Mode
const (
	SSCMode1 byte = 0b001
	SSCMode2 byte = 0b010
	SSCMode3 byte = 0b011
	Unused1  byte = 0b100
	Unused2  byte = 0b101
	Unused3  byte = 0b110
)

// Constants for SessionAMBR
const (
	ValNotUsed   byte = 0b00000000
	Mult_1kbps   byte = 0b00000001
	Mult_4kbps   byte = 0b00000010
	Mult_16kbps  byte = 0b00000011
	Mult_64kbps  byte = 0b00000100
	Mult_256kbps byte = 0b00000101
	Mult_1mbps   byte = 0b00000110
	Mult_4mbps   byte = 0b00000111
	Mult_16mbps  byte = 0b00001000
	Mult_64mbps  byte = 0b00001001
	Mult_256mbps byte = 0b00001010
	Mult_1gbps   byte = 0b00001011
	Mult_4gbps   byte = 0b00001100
	Mult_16gbps  byte = 0b00001101
	Mult_64gbps  byte = 0b00001110
	Mult_256gbps byte = 0b00001111
	Mult_1tbps   byte = 0b00010000
	Mult_4tbps   byte = 0b00010001
	Mult_16tbps  byte = 0b00010010
	Mult_64tbps  byte = 0b00010011
	Mult_256tbps byte = 0b00010100
	Mult_1pbps   byte = 0b00010101
	Mult_4pbps   byte = 0b00010110
	Mult_16pbps  byte = 0b00010111
	Mult_64pbps  byte = 0b00011000
	Mult_256pbps byte = 0b00011001
)

// Constants for SM Cause
const (
	InsufficientResources            byte = 0b00011010
	MissingOrUnknownDNN              byte = 0b00011011
	UnknownPduSType                  byte = 0b00011100
	UserAutheticationFail            byte = 0b00011101
	RequestRejectedUnspecified       byte = 0b00011111
	ServiceOptionOutOfOrder          byte = 0b00100010
	PTIAlreadyInUse                  byte = 0b00100011
	RegularDeactivation              byte = 0b00100100
	NetworkFailure                   byte = 0b00100110
	ReactivationRequested            byte = 0b00100111
	SemanticErrorTFT                 byte = 0b00101001
	SyntacticalErrorTFT              byte = 0b00101010
	InvalidPDUSIdentity              byte = 0b00101011
	SemanticErrorPacketFilters       byte = 0b00101100
	SyntacticalErrorPacketFilters    byte = 0b00101101
	OutOfLADNArea                    byte = 0b00101110
	PTIMismatch                      byte = 0b00101111
	PDUSTypeIPV4                     byte = 0b00110010
	PDUSTypeIPV6                     byte = 0b00110011
	PDUSDoesnotExist                 byte = 0b00110110
	PDUSTypeIPV4V6                   byte = 0b00111001
	PDUSTypeUnstructured             byte = 0b00111010
	Unsupported5QI                   byte = 0b00111011
	PDUSTypeEthernet                 byte = 0b00111101
	InsufficientResourcesSliceAndDNN byte = 0b01000011
	UnsupportedSSCMode               byte = 0b01000100
	InsufficientResourcesSlice       byte = 0b01000101
	MissingOrUnknownDNNSLice         byte = 0b01000110
	InvalidPTI                       byte = 0b01010001
	MaxDataRateLow                   byte = 0b01010010
	SemanticErrorQoS                 byte = 0b01010011
	SyntacticalErrorQoS              byte = 0b01010100
	SemanticallyIncorrectMsg         byte = 0b01011111
	InvalidMandatoryInfo             byte = 0b01100000
	MessageTypeNonExistent           byte = 0b01100001
	MesageTypeNotCompatible          byte = 0b01100010
	IENonExistent                    byte = 0b01100011
	ConditionalIEError               byte = 0b01100100
	MsgNotCompatibleProtocolState    byte = 0b01100101
	ProtocolErrorUnspecified         byte = 0b01101111
)

// Constants for NSSAI Length
const (
	Sst           byte = 0b00000001
	SstHPLMNSst   byte = 0b00000010
	SstSD         byte = 0b00000100
	SstSDHPLMNSst byte = 0b00000101
	All           byte = 0b00001000
)

// Constants for Security Header Type
const (
	PlainNAS           byte = 0000
	IntegrityProtected byte = 0001
	IntegrityCipher    byte = 0010
	IntegrityNew       byte = 0011
	IntegrityCipherNew byte = 0100
)

// Constants for 5GS Registration Type
const (
	InitialReg         byte = 001
	MobilityReg        byte = 010
	PeriodicReg        byte = 011
	EmergencyReg       byte = 100
	SNPNReg            byte = 101
	RoamingMobilityReg byte = 110
	RoamingInitialReg  byte = 111
)

// Constants for Follow On request
const (
	NoFOR byte = 0
	FOR   byte = 1
)

// Constants for NAS KSI
const (
	Native byte = 0
	Mapped byte = 1
	NoKey  byte = 111
)

// Constants for Registration Result
const (
	ThreeGPPAccess byte = 001
	Non3GPPAccess  byte = 010
	Both           byte = 011
	Reserved       byte = 111
)

// Constants for SMS
const (
	Allowed    byte = 1
	NotAllowed byte = 0
)

// Constants for NSSAA
const (
	NotPerformed  byte = 0
	ToBePerformed byte = 1
)

// Constants for Emergency Services
const (
	Registered    byte = 1
	NotRegistered byte = 0
)

const (
	SpareOctet byte = 0000
	Spare      byte = 0
	OneLength  byte = 00000001
)

// Constants for Payload Container Type
const (
	N1SmInfo      byte = 0b0001
	Sms           byte = 0b0010
	LppMsg        byte = 0b0011
	SorTran       byte = 0b0100
	UEPolCont     byte = 0b0101
	UEParaUpdTran byte = 0b0110
	LocaServMsg   byte = 0b0111
	CIoTUser      byte = 0b1000
	SLAACont      byte = 0b1001
	EventNotif    byte = 0b1010
	MultiPay      byte = 0b1111
)
