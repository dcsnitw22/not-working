package nas

func EncodeQoSRules(qos QoSRule) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)

	//Rule Identifier is first octet
	ruleIdentifier, err := EncodeQoSIdentifier(qos.QoSIdentifier)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, ruleIdentifier)

	//Length is the second octet
	if qos.Operation == "Delete Existing QoS rule" {
		length := 2 + len(qos.PacketFilterList)
		byteArray = append(byteArray, byte(length))

	} else {
		length := 3 + len(qos.PacketFilterList)
		byteArray = append(byteArray, byte(length))
	}

	//Rule operation code, DQR Bit, Packer filter length is octet 3
	ruleOpCode, err := EncodeQoSRuleOperationCode(qos.Operation)
	if err != nil {
		return nil, err
	}

	dqrBit, err := EncodeDQR(qos.DQR)
	if err != nil {
		return nil, err
	}

	var packetFilterLen byte
	if qos.Operation == "Delete existing QoS rule" || qos.Operation == "Modify existing QoS rule without modifying packet filters" {
		packetFilterLen = byte(0)
	} else {
		packetFilterLen = byte(len(qos.PacketFilterList))
	}

	first3bits := (ruleOpCode & 0b00000111) << 5
	fourthBit := (dqrBit & 0b00000001) << 4
	last4Bits := (packetFilterLen & 0b00001111)

	octet3 := first3bits | fourthBit | last4Bits

	byteArray = append(byteArray, octet3)

	//Next Octets are packet filter
	if qos.Operation != "Delete existing QoS rule" || qos.Operation != "Modify existing QoS rule without modifying packet filters" {

		pfOctet, err := EncodePacketFilterList(qos)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, pfOctet...)
	}

	//QosRule Precedence Octet
	if qos.Operation != "Delete Existing QoS rule" {
		byteArray = append(byteArray, byte(qos.Precedence))
	}

	//Last Octet
	seg, err := EncodeSegregation(qos.Segregation)
	if err != nil {
		return nil, err
	}

	qfi, err := EncodeQFI(qos.QFI)
	if err != nil {
		return nil, err
	}

	firstBit := (SpareOctet & 0b00000001) << 7
	secondBit := (seg & 0b00000001) << 6
	last6Bits := (qfi & 0b00111111)

	lastOctet := firstBit | secondBit | last6Bits

	byteArray = append(byteArray, lastOctet)

	return byteArray, nil
}
