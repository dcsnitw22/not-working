package nas

func EncodePacketFilterList(qos QoSRule) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)

	if qos.Operation == "Modify existing QoS rule and delete packet filters" {
		bytes, err := EncodeModDelPF(qos.PacketFilterList)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, bytes...)
	} else {
		bytes, err := EncodePF(qos.PacketFilterList)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, bytes...)
	}

	return byteArray, nil

}

func EncodeModDelPF(pfl []PacketFilter) ([]byte, error) {
	byteArray := make([]byte, 0)
	for index := 0; index < len(pfl); index++ {
		// list, ok := pfl[index].(PacketFilterModDel)
		// if ok {
		//Only octet per item - Identifier
		byteArray = append(byteArray, byte(pfl[index].Identifier))

		// } else {
		// 	return nil, errors.New("error converting the packet filter list")
		// }

	}

	return byteArray, nil

}

func EncodePF(pfl []PacketFilter) ([]byte, error) {
	byteArray := make([]byte, 0)
	for index := 0; index < len(pfl); index++ {
		// list, ok := pfl[index].(PacketFilter)
		list := pfl[index]
		// if ok {
		//First Octet - spare, direction, identifier
		first2bits := (byte(0b00000000) & 0b00000011) << 6

		dir, err := EncodePacketFilterDirection(list.Direction)
		if err != nil {
			return nil, err
		}
		next2bits := (dir & 0b00000011) << 4

		id := byte(list.Identifier)

		last4bits := id & 0b00001111

		resultOctet := first2bits | next2bits | last4bits

		byteArray = append(byteArray, resultOctet)

		//Second Octet - Length
		byteArray = append(byteArray, byte(len(list.Components)))

		//Rest of the octets - Packet Filter Component
		idlist, err := EncodeIDList(list.Components)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, idlist...)

		// } else {
		// 	return nil, errors.New("error converting the packet filter list")
		// }

	}

	return byteArray, nil

}

// Packet Filter Component List
func EncodeIDList(ids []string) ([]byte, error) {
	byteArray := make([]byte, 0)

	for _, str := range ids {
		comp, err := EncodePacketFilterComponentTypeIdentifier(str)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, comp)

	}
	return byteArray, nil
}
