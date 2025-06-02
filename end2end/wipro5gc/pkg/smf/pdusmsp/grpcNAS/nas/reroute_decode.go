package nas

type NasDecodedMsg interface{}

func ReRouteDecode(epd string, messageType string, byteArray []byte) (NasDecodedMsg, error) {
	var nasDecodedMsg NasDecodedMsg

	if epd == "MOBILITY_MANAGEMENT_MESSAGES" {
		switch messageType {
		case "UL_NAS_TRANSPORT":
			nasMsg, err := DecodeULNas(byteArray)
			if err != nil {
				return nil, err
			}
			nasDecodedMsg = nasMsg
		case "REGISTRATION_REQUEST":
			nasMsg, err := DecodeRegistrationRequest(byteArray)
			if err != nil {
				return nil, err
			}
			nasDecodedMsg = nasMsg

		}
	}
	if epd == "SESSION_MANAGEMENT_MESSAGES" {
		switch messageType {
		case "PDU_SESSION_ESTABLISHMENT_REQUEST":
			nasMsg, err := DecodePduSessionEstablishmentRequest(byteArray)
			if err != nil {
				return nil, err
			}
			nasDecodedMsg = nasMsg
		case "PDU_SESSION_MODIFICATION_REQUEST":
			nasMsg, err := DecodePduSessionModificationRequest(byteArray)
			if err != nil {
				return nil, err
			}
			nasDecodedMsg = nasMsg
		case "PDU_SESSION_RELEASE_REQUEST":
			nasMsg, err := DecodePduSessionReleaseRequest(byteArray)
			if err != nil {
				return nil, err
			}
			nasDecodedMsg = nasMsg
		}

	}

	return nasDecodedMsg, nil

}
