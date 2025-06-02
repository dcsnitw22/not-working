package nas

func DecodePduSessionEstablishmentRequest(byteArray []byte) (PduSessionEstablishmentRequest, error) {
	var pdu PduSessionEstablishmentRequest

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
		//Second Octet is PDU Session ID
		case 1:
			pduSessionID, err := DecodePduSessionID(byteArray[1])
			if err != nil {
				return pdu, err
			}
			pdu.PDUsessionId = pduSessionID
		//Third octet is PTI
		case 2:
			pdu.PTI = int(byteArray[2])
		//Fouth Octet is message type
		case 3:
			messageType, err := DecodeMessageType(byteArray[3])
			if err != nil {
				return pdu, err
			}
			pdu.MessageType = messageType
		//Fifth Octet is Max Integrity Protected DataRate for UL
		case 4:
			maxIntegrityProtectedDataRateUl, err := DecodeMaximumDataRatePerUeForUserPlaneIntegrityProtectionForUplink(byteArray[4])
			if err != nil {
				return pdu, err
			}
			pdu.MaxIntegrityProtectedDataRateUl = maxIntegrityProtectedDataRateUl
		//Sixth Octet is Max Integrity Protected DataRate for DL
		case 5:
			maxIntegrityProtectedDataRateDl, err := DecodeMaximumDataRatePerUeForUserPlaneIntegrityProtectionForDownlink(byteArray[5])
			if err != nil {
				return pdu, err
			}
			pdu.MaxIntegrityProtectedDataRateDl = maxIntegrityProtectedDataRateDl

		}
		index++
	}

	return pdu, nil
}
