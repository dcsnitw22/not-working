package nas

func EncodeDLNAS(pdu DLNasModel) ([]byte, error) {
	// // Open a file for writing, create it if it doesn't exist, truncate it if it does
	// file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {

	// 	return file, err
	// }
	// defer file.Close()

	// Create an empty byte array
	byteArray := make([]byte, 0)

	// First octet is Extended Protocol Discriminator
	epd, err := EncodeEpd(pdu.ExtendedProtocolDiscriminator)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, epd)

	//Second half octet is Security Header Type
	sht, err := EncodeSecurityHeader(pdu.SecurityHeaderType)
	if err != nil {
		return nil, err
	}
	spareOctet := byte(0b0000)
	secondoctet := (sht << 4) | spareOctet
	byteArray = append(byteArray, secondoctet)

	//Third Octet is message type
	messageType, err := EncodeMessageType(pdu.MessageType)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, messageType)

	//Fourth Octet is Payload Container Type
	plcType, err := EncodePayloadContainerType(pdu.PayLoadContainerType)
	if err != nil {
		return nil, err
	}
	// spareOctet = byte(0b0000)
	// secondoctet = (plcType << 4) | spareOctet
	// byteArray = append(byteArray, secondoctet)
	byteArray = append(byteArray, plcType)

	//Rest Of the octets are PayloadContainer
	// payload, err := EncodeDLPayload(pdu.PayLoadContainerType, pdu.PayLoadContainer)
	// if err != nil {
	// 	return nil, err
	// }
	payload := pdu.PayLoadContainer
	byteArray = append(byteArray, byte(0))
	byteArray = append(byteArray, byte(len(payload)))
	byteArray = append(byteArray, payload...)

	// // Write the byte array to the file
	// _, err = file.Write(byteArray)
	// if err != nil {
	// 	return file, err
	// }

	if pdu.PayLoadContainerType == "N1 SM information" || pdu.PayLoadContainerType == "CIoT user data container" {
		byteArray = append(byteArray, byte(pdu.PduSessionIdIEI))
		byteArray = append(byteArray, byte(pdu.PduSessionId))

	}

	return byteArray, nil

}
