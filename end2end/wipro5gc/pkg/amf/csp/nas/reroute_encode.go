package nas

func ReRouteEncode(epd string, messageType string, msg interface{}) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)
	if epd == "MOBILITY_MANAGEMENT_MESSAGES" {
		switch messageType {
		case "DL_NAS_TRANSPORT":
			res, err := EncodeDLNAS(msg.(DLNasModel))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)
		case "REGISTRATION_ACCEPT":
			res, err := EncodeRegistrationAccept(msg.(RegistrationAcceptModel))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)

		}
	}

	if epd == "SESSION_MANAGEMENT_MESSAGES" {
		switch messageType {
		case "PDU_SESSION_ESTABLISHMENT_ACCEPT":
			res, err := EncodePduSessionEstablishmentAccept(msg.(PduSessionEstablishmentAccept))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)
		case "PDU_SESSION_ESTABLISHMENT_REJECT":
			res, err := EncodePduSessionEstablishmentReject(msg.(PduSessionEstablishmentReject))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)

		case "PDU_SESSION_MODIFICATION_REJECT":
			res, err := EncodePduSessionModificationReject(msg.(PduSessionModificationReject))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)
		case "PDU_SESSION_RELEASE_REJECT":
			res, err := EncodePduSessionReleaseReject(msg.(PduSessionReleaseReject))
			if err != nil {
				return nil, err
			}
			byteArray = append(byteArray, res...)

		}
	}

	return byteArray, nil

}
