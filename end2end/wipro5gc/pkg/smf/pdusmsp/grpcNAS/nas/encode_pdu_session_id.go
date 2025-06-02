package nas

//Logic to encode PDUSessionID

func EncodePduSessionID(pdusid int) (byte, error) {
	return byte(pdusid), nil
	// switch pdusid {
	// case "NO_VAL":
	// 	return NoVal, nil
	// case "VAL_1":
	// 	return Val1, nil
	// case "VAL_2":
	// 	return Val2, nil
	// case "VAL_3":
	// 	return Val3, nil
	// case "VAL_4":
	// 	return Val4, nil
	// case "VAL_5":
	// 	return Val5, nil
	// case "VAL_6":
	// 	return Val6, nil
	// case "VAL_7":
	// 	return Val7, nil
	// case "VAL_8":
	// 	return Val8, nil
	// case "VAL_9":
	// 	return Val9, nil
	// case "VAL_10":
	// 	return Val10, nil
	// case "VAL_11":
	// 	return Val11, nil
	// case "VAL_12":
	// 	return Val12, nil
	// case "VAL_13":
	// 	return Val13, nil
	// case "VAL_14":
	// 	return Val14, nil
	// case "VAL_15":
	// 	return Val15, nil
	// default:
	// 	return 0, errors.New("invalid PDU Session ID")
	// }

}
