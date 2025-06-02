package nas

// Struct for PDU session Establishment Accept

type PduSessionEstablishmentAccept struct {
	ExtendedProtocolDiscriminator string
	PDUsessionId                  int // Possible values: 0 to 16
	PTI                           int // Possible values: 0 to 254
	MessageType                   string
	PduSessionType                string
	SSCmode                       string
	QosRuleIEI                    int
	AuthorizedQoSRules            []QoSRule
	SessionAmbr                   SessionAMBR
	// Nssai                         Nssai
}
