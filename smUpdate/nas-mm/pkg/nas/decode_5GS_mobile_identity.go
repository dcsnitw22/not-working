package nas

// Only for SUPI - ISMI - No scheme
func Decode5GSMobilityIdentity(ba []byte) (SupiIdentity, error) {
	var supi SupiIdentity
	var supiRes SupiIdentity
	var err error
	var length int
	index := 0
	for index < int(len(ba)) {
		switch index {
		case 0:
			iei := DecodeIEI(ba[0])
			supi.IEI = iei
		case 1:
			length = int(ba[1])

		case 2:
			supiRes, err = GetDetails(ba[2 : 2+length])
			if err != nil {
				return supi, err
			}
			supi.IdentityType = supiRes.IdentityType
			supi.MCC = supiRes.MCC
			supi.MNC = supiRes.MNC
			supi.ProtectionSchemeID = supiRes.ProtectionSchemeID
			supi.RoutingIndicator = supiRes.RoutingIndicator
			supi.SUPIFormat = supiRes.SUPIFormat
			supi.SchemeOutput = supiRes.SchemeOutput
			supi.HNPKIdentifier = supiRes.HNPKIdentifier
		}

		index++
	}

	return supi, nil
}

func GetDetails(byteArray []byte) (SupiIdentity, error) {
	var supi SupiIdentity
	index := 0
	for index < int(len(byteArray)) {
		switch index {
		case 0:
			last3bits := byteArray[0] & 0b00000111
			id, err := Decode5GSIdentityType(last3bits)
			if err != nil {
				return supi, err
			}
			supi.IdentityType = id

			// firstbits := (byteArray[0] & 0b11100000) >> 5
			firstbits := (byteArray[0] >> 5) & 0b00000111

			supiFormat, err := DecodeSUPIFormat(firstbits)
			if err != nil {
				return supi, err
			}
			supi.SUPIFormat = supiFormat
		case 1:
			mcc1bits := byteArray[1] & 0b00001111
			mcc2bits := (byteArray[1] & 0b11110000) >> 4
			mcc3bits := byteArray[2] & 0b00001111
			mnc1bits := byteArray[3] & 0b00001111
			mnc2bits := (byteArray[3] & 0b11110000) >> 4
			mnc3bits := (byteArray[2] & 0b11110000) >> 4
			var mccarray = []int{int(mcc1bits), int(mcc2bits), int(mcc3bits)}
			var mncarray []int
			if int(mnc3bits) == 15 {
				mncarray = []int{int(mnc1bits), int(mnc2bits)}
			} else {
				mncarray = []int{int(mnc1bits), int(mnc2bits), int(mnc3bits)}

			}
			supi.MCC = MakeString(mccarray)
			supi.MNC = MakeString(mncarray)
		case 4:
			routing1bits := byteArray[4] & 0b00001111
			routing2bits := (byteArray[4] & 0b11110000) >> 4
			routing3bits := byteArray[5] & 0b00001111
			routing4bits := (byteArray[5] & 0b11110000) >> 4
			var routingarray = []int{int(routing1bits), int(routing2bits), int(routing3bits), int(routing4bits)}
			supi.RoutingIndicator = MakeString(routingarray)
		case 6:
			supi.ProtectionSchemeID = DecodeProtectionSchemeID(byteArray[6])
		case 7:
			supi.HNPKIdentifier = DecodeHNPKI(byteArray[7])
		case 8:
			supi.SchemeOutput = DecodeSchemaOutput(byteArray[8:])
		}

		index++
	}

	// fmt.Println(supi)
	return supi, nil

}
