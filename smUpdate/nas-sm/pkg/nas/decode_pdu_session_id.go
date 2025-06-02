package nas

//Logic to Decode PDUSessionID

func DecodePduSessionID(pdusid byte) (int, error) {
	return int(pdusid), nil

	// switch pdusid {
	// case NoVal:
	// 	return "NO_VAL", nil
	// case Val1:
	// 	return "VAL_1", nil
	// case Val2:
	// 	return "VAL_2", nil
	// case Val3:
	// 	return "VAL_3", nil
	// case Val4:
	// 	return "VAL_4", nil
	// case Val5:
	// 	return "VAL_5", nil
	// case Val6:
	// 	return "VAL_6", nil
	// case Val7:
	// 	return "VAL_7", nil
	// case Val8:
	// 	return "VAL_8", nil
	// case Val9:
	// 	return "VAL_9", nil
	// case Val10:
	// 	return "VAL_10", nil
	// case Val11:
	// 	return "VAL_11", nil
	// case Val12:
	// 	return "VAL_12", nil
	// case Val13:
	// 	return "VAL_13", nil
	// case Val14:
	// 	return "VAL_14", nil
	// case Val15:
	// 	return "VAL_15", nil
	// default:
	// 	return "", errors.New("invalid PDU Session ID")
	// }

}
