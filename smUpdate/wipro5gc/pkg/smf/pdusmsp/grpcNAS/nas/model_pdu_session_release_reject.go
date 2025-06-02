package nas

// Struct for PDU session Release Reject

type PduSessionReleaseReject struct {
	ExtendedProtocolDiscriminator string
	PDUsessionId                  int // Possible values: 0 to 16
	PTI                           int // Possible values: 0 to 254
	MessageType                   string
	SMcause                       string
}
