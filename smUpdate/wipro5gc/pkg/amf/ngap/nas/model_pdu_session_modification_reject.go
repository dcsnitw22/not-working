package nas

// Struct for PDU session Modification Reject

type PduSessionModificationReject struct {
	ExtendedProtocolDiscriminator string
	PDUsessionId                  string // Possible values: 0 to 16
	PTI                           int    // Possible values: 0 to 254
	MessageType                   string
	SMcause                       string
}
