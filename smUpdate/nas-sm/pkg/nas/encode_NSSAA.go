package nas

import "errors"

func EncodeNSSAA(nssaa string) (byte, error) {
	switch nssaa {
	case "Network slice-specific authentication and authorization is not to be performed":
		return NotPerformed, nil
	case "Network slice-specific authentication and authorization is to be performed":
		return ToBePerformed, nil
	default:
		return 0, errors.New("invalid NSSAA")
	}
}
