package nas

func DecodeReleaseComplete(byteArray []byte) (ReleaseCompleteModel, error) {

	var pdu ReleaseCompleteModel

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
		case 4:
			if byteArray[4] == byte(89) {
				smCause, err := DecodeSmCause(byteArray[5])
				if err != nil {
					return pdu, err
				}
				pdu.SMCause = smCause
			}
		}
		index++
	}

	return pdu, nil

}
