package nas

type ULNasModel struct {
	ExtendedProtocolDiscriminator string
	SecurityHeaderType            string
	MessageType                   string
	PayLoadContainerType          string
	PayLoadContainer              interface{}
	PduSessionIdIEI               int
	PduSessionId                  int
}
