package nas

type QoSRuleOperationCode uint8

type PacketFilterList []PacketFilter

type PacketFilterDirection uint8

//type PacketFilterComponentList []PacketFilterComponent

type PacketFilter struct {
	Identifier uint8
	Direction  PacketFilterDirection
	// Components PacketFilterComponentList
}

type QoSRule struct {
	QoSIdentifier    string
	Operation        QoSRuleOperationCode
	DQR              bool
	PacketFilterList PacketFilterList
	Precedence       uint8
	Segregation      bool
	QFI              uint8
}
