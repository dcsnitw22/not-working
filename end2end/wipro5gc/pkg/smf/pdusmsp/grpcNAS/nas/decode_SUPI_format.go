package nas

import "errors"

func DecodeSUPIFormat(format byte) (string, error) {
	switch format {
	case IMSI:
		return "IMSI", nil
	case NetworkSpecId:
		return "Network specific identifier", nil
	case GCI:
		return "GCI", nil
	case GLI:
		return "GLI", nil
	default:
		return "", errors.New("unknown SUPI Format")
	}

}
