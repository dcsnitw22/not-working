package nas

func EncodePduSessionAddress(si6lla string, pdusessionType string, ipAddr string) ([]byte, error) {
	byteArray := make([]byte, 0)

	//First Octet is IEI - 29
	byteArray = append(byteArray, byte(41))

	//Second Octet is length
	if pdusessionType == "IPV4" {
		byteArray = append(byteArray, byte(5))
	}

	//Third Octet is Pdu Session Type & SI6LLA
	if pdusessionType == "IPV4" {
		pdusType, err := EncodePduSessionType(pdusessionType)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, pdusType)
	}

	//Next Bytes are Pdu Session Address
	pduAddrInfo, err := EncodePduAddress(ipAddr)
	if err != nil {
		return nil, err
	}

	byteArray = append(byteArray, pduAddrInfo...)

	return byteArray, nil
}
