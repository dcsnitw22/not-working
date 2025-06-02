package nas

func DecodeULNas(byteArray []byte) (ULNasModel, error) {
	var pdu ULNasModel

	// // Get the file size
	// fileInfo, err := binaryFile.Stat()
	// if err != nil {

	// 	return pdu, err
	// }
	// fileSize := fileInfo.Size()

	// // Create a byte array with the size of the file
	// byteArray := make([]byte, fileSize)

	// // Read the bytes from the file into the byte array
	// _, err = binaryFile.Read(byteArray)
	// if err != nil {

	// 	return pdu, err
	// }

	index := 0
	for index < int(len(byteArray)) {
		switch index {
		//First octet is Extended Protocol Discriminator
		case 0:
			epd, err := DecodeEpd(byteArray[0])
			if err != nil {
				return pdu, err
			}
			pdu.ExtendedProtocolDiscriminator = epd
		case 1:
			sht, err := DecodeSecurityHeader(byteArray[1] & 0x0F)
			if err != nil {
				return pdu, err
			}
			pdu.SecurityHeaderType = sht
		case 2:
			messageType, err := DecodeMessageType(byteArray[2])
			if err != nil {
				return pdu, err
			}
			pdu.MessageType = messageType
		case 3:
			payloadType, err := DecodePayloadContainerType(byteArray[3] & 0x0F)
			if err != nil {
				return pdu, err
			}
			pdu.PayLoadContainerType = payloadType
		case 4:
			payload, err := DecodeULPayload(pdu.PayLoadContainerType, byteArray[6:])
			if err != nil {
				return pdu, err
			}
			pdu.PayLoadContainer = payload
		case 5:
			if pdu.PayLoadContainerType == "N1 SM information" || pdu.PayLoadContainerType == "CIoT user data container" {
				len := int(byteArray[5])
				pdu.PduSessionIdIEI = int(byteArray[len+6])
				pdu.PduSessionId = int(byteArray[len+7])
			}
		}
		index++
	}

	return pdu, nil

}
