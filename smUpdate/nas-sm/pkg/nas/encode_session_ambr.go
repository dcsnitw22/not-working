package nas

func EncodeSessionAmbr(sessionAMBR SessionAMBR) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)
	//First Octet - Session AMBR IEI
	// byteArray = append(byteArray, byte(sessionAMBR.IEI))
	//Second Octet - Length
	byteArray = append(byteArray, byte(6))
	//Third Octet - Unit for DL
	dl, err := EncodeUnitSessionAMBR(sessionAMBR.UnitDL)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, dl)
	//Next 2 Octets - Max Rate for DL
	// byteArray = append(byteArray, byte(sessionAMBR.RateDL))
	dlArr := decimalTo16BitBytes(sessionAMBR.RateDL)
	byteArray = append(byteArray, dlArr...)

	//Fifth Octet - Unit for UL
	ul, err := EncodeUnitSessionAMBR(sessionAMBR.UnitUL)
	if err != nil {
		return nil, err
	}
	byteArray = append(byteArray, ul)
	//Next 2 Octets - Max Rate for UL
	// byteArray = append(byteArray, byte(sessionAMBR.RateUL))
	ulArr := decimalTo16BitBytes(sessionAMBR.RateUL)
	byteArray = append(byteArray, ulArr...)

	return byteArray, err

}
