package nas

import "errors"

func EncodeRegistrationResult(reg string) (byte, error) {
	switch reg {
	case "3GPP access":
		return ThreeGPPAccess, nil
	case "Non-3GPP access":
		return Non3GPPAccess, nil
	case "3GPP access and non-3GPP access":
		return Both, nil
	case "Reserved":
		return Reserved, nil
	default:
		return 0, errors.New("invalid Registration Result")
	}
}
