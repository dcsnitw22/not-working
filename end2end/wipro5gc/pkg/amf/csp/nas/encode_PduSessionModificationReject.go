package nas

func EncodePduSessionModificationReject(pdu PduSessionModificationReject) ([]byte, error) {

	// // Open a file for writing, create it if it doesn't exist, truncate it if it does
	// file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {

	// 	return file, err
	// }
	// defer file.Close()

	// Create an empty byte array
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
	smCause, err := EncodeSMcause(pdu.SMcause)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, smCause)

	// // Write the byte array to the file
	// _, err = file.Write(byteArray)
	// if err != nil {

	// 	return file, err
	// }

	return byteArray, nil

}
