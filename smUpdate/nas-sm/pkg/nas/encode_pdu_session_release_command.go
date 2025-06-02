package nas

func EncodeReleaseCommand(pdu ReleaseCommandModel) ([]byte, error) {
	//Create an empty byte array
	byteArray := make([]byte, 0)

	//First octet is Extended Protocol Discriminator
	epd, err := EncodeEpd(pdu.ExtendedProtocolDiscriminator)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, epd)

	//Second Octet is PDU Session ID
	pduSessionID, err := EncodePduSessionID(pdu.PDUsessionId)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, pduSessionID)

	//Third octet is PTI
	byteArray = append(byteArray, byte(pdu.PTI))

	//Fouth Octet is message type
	messageType, err := EncodeMessageType(pdu.MessageType)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, messageType)

	//Fifth Octet is SM Cause
	smCause, err := EncodeSMcause(pdu.SMCause)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, smCause)

	return byteArray, nil

}
