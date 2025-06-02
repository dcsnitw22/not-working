package nas

func EncodeAuthorizedQoSRules(pdu PduSessionEstablishmentAccept) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)
	//First Octet is QoS IEI
	byteArray = append(byteArray, byte(pdu.QosRuleIEI))
	//Second  is length
	byteArray = append(byteArray, byte(len(pdu.AuthorizedQoSRules)))
	// byteArray = append(byteArray, byte(len(pdu.AuthorizedQoSRules)))
	//Rest of the Octets are QoSRules
	for index := 0; index < len(pdu.AuthorizedQoSRules); index++ {
		qos, err := EncodeQoSRules(pdu.AuthorizedQoSRules[index])
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, qos...)
	}

	byteArray[1] = byte(len(byteArray) - 2)

	// lenBytes := decimalTo16BitBytes(len(byteArray) - 3)
	// byteArray[1] = lenBytes[0]
	// byteArray[2] = lenBytes[1]

	return byteArray, nil
}
