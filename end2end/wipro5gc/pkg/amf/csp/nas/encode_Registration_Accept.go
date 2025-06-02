package nas

func EncodeRegistrationAccept(regAccept RegistrationAcceptModel) ([]byte, error) {

	// // Open a file for writing, create it if it doesn't exist, truncate it if it does
	// file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {

	// 	return file, err
	// }
	// defer file.Close()

	// Create an empty byte array
	byteArray := make([]byte, 0)

	//First octet is Extended Protocol Discriminator
	epd, err := EncodeEpd(regAccept.ExtendedProtocolDiscriminator)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, epd)

	//Second half octet is Security Header Type

	sht, err := EncodeSecurityHeader(regAccept.SecurityHeaderType)
	if err != nil {
		return nil, err
	}
	spareOctet := byte(0b0000)
	secondoctet := (sht << 4) | spareOctet
	byteArray = append(byteArray, secondoctet)

	//Third Octet is Message Type
	messageType, err := EncodeMessageType(regAccept.MessageType)
	if err != nil {
		return nil, err

	}
	byteArray = append(byteArray, messageType)

	//Fourth Octet is registration result

	byteArray = append(byteArray, OneLength)

	resSpare := byte(0b0)

	regRes, err := EncodeRegistrationResult(regAccept.RegResult)
	if err != nil {
		return nil, err

	}

	sms, err := EncodeSMS(regAccept.Sms)
	if err != nil {
		return nil, err

	}

	nssaa, err := EncodeNSSAA(regAccept.NssaPerformed)
	if err != nil {
		return nil, err

	}

	emeReg, err := EncodeEmergencyReg(regAccept.EmergencyReg)
	if err != nil {
		return nil, err

	}

	roaming, err := EncodeRoamingReg(regAccept.RoamingReg)
	if err != nil {
		return nil, err

	}

	fourthOctet := resSpare | roaming | emeReg | nssaa | sms | regRes

	byteArray = append(byteArray, fourthOctet)

	// // Write the byte array to the file
	// _, err = file.Write(byteArray)
	// if err != nil {

	// 	return file, err
	// }

	return byteArray, nil

}
