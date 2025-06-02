package nas

func DecodeRegistrationRequest(byteArray []byte) (RegistrationRequestModel, error) {
	var pdu RegistrationRequestModel

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
			regBits := byteArray[3] >> 4
			regValBits := regBits >> 1
			forBits := regBits & 0x01

			nasKsiBits := byteArray[3] & 0x0F
			tscBits := nasKsiBits >> 3
			ksiBits := nasKsiBits & 0x07

			regV, regF, err := Decode5GSRegistrationType(regValBits, forBits)
			if err != nil {
				return pdu, err
			}
			pdu.FORValue = regF
			pdu.RegistrationType = regV

			tsc, ksi, err := DecodeNASKSI(tscBits, ksiBits)
			if err != nil {
				return pdu, err
			}
			pdu.NASTSC = tsc
			pdu.NASKSI = ksi

		}

		index++

	}

	return pdu, nil
}
