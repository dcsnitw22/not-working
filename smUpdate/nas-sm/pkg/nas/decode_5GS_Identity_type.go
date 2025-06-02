package nas

import (
	"errors"
)

func Decode5GSIdentityType(idtype byte) (string, error) {
	switch idtype {
	case NoIdentity:
		return "No Identity", nil
	case SUCI:
		return "SUCI", nil
	case FiveGGUTI:
		return "5G-GUTI", nil
	case IEMI:
		return "IEMI", nil
	case TMSI:
		return "5G-S-TMSI", nil
	case IEMISV:
		return "IMEISV", nil
	case MAC:
		return "MAC address", nil
	case EUI:
		return "EUI-64", nil
	default:
		return "", errors.New("unknown 5GS Identity Type")
	}

}
