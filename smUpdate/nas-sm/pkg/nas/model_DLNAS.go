package nas

type DLNasModel struct {
	ExtendedProtocolDiscriminator string
	SecurityHeaderType            string
	MessageType                   string
	PayLoadContainerType          string
	PayLoadContainer              interface{}
	PduSessionIdIEI               int
	PduSessionId                  int
}
