package nas

type ReleaseCommandModel struct {
	ExtendedProtocolDiscriminator string
	PDUsessionId                  int // Possible values: 0 to 16
	PTI                           int // Possible values: 0 to 254
	MessageType                   string
	SMCause                       string
}
