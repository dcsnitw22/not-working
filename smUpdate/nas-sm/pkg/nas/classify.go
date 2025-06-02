package nas

func Classify(byteArray []byte) (epd string, messageType string, e error) {
	epd = ""
	messageType = ""
	// // Get the file size
	// fileInfo, err := binaryFile.Stat()
	// if err != nil {

	// 	return epd, messageType, err
	// }
	// fileSize := fileInfo.Size()

	// // Create a byte array with the size of the file
	// byteArray := make([]byte, fileSize)

	// // Read the bytes from the file into the byte array
	// _, err = binaryFile.Read(byteArray)
	// if err != nil {

	// 	return epd, messageType, err
	// }
	epd, err := DecodeEpd(byteArray[0])
	if err != nil {
		return epd, messageType, err
	}

	if epd == "MOBILITY_MANAGEMENT_MESSAGES" {
		messageType, err = DecodeMessageType(byteArray[2])
		if err != nil {
			return epd, messageType, err
		}
	} else if epd == "SESSION_MANAGEMENT_MESSAGES" {
		messageType, err = DecodeMessageType(byteArray[3])
		if err != nil {
			return epd, messageType, err
		}
	}

	return epd, messageType, nil

}
