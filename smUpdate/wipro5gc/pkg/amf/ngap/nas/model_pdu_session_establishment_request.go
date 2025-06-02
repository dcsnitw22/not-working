package nas

// Struct for PDU session Establishment Request

type PduSessionEstablishmentRequest struct {
	ExtendedProtocolDiscriminator string
	PDUsessionId                  string // Possible values: 0 to 16
	PTI                           int    // Possible values: 0 to 254
	MessageType                   string
	//Integrity protection maximum data rate - 2 Octets
	MaxIntegrityProtectedDataRateUl string
	MaxIntegrityProtectedDataRateDl string
}
