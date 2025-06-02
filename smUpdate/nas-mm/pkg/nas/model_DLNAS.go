package nas

type DLNasModel struct {
	ExtendedProtocolDiscriminator string
	SecurityHeaderType            string
	MessageType                   string
	PayLoadContainerType          string
	// PayLoadContainer              interface{}
	PayLoadContainer []byte
	PduSessionIdIEI  int
	PduSessionId     int
}
