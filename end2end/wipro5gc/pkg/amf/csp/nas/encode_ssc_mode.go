package nas

import (
	"errors"
)

//Logic to encode SSC Mode

func EncodeSSCmode(sscMode string) (byte, error) {
	switch sscMode {
	case "SSC_MODE_1":
		return SSCMode1, nil
	case "SSC_MODE_2":
		return SSCMode2, nil
	case "SSC_MODE_3":
		return SSCMode3, nil
	case "UNUSED_1":
		return Unused1, nil
	case "UNUSED_2":
		return Unused2, nil
	case "UNUSED_3":
		return Unused3, nil
	default:
		return 0, errors.New("invalid SSC Mode")
	}
}
