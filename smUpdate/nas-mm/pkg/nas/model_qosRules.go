package nas

type QoSRule struct {
	QoSIdentifier    string
	Operation        string
	DQR              string
	PacketFilterList []PacketFilter
	Precedence       uint8
	Segregation      string
	QFI              string
}
