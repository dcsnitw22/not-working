package nas

import "errors"

func EncodeRoamingReg(roReg string) (byte, error) {
	switch roReg {
	case "No additional information":
		return NotRegistered, nil
	case "Request for registration for disaster roaming services accepted as registration not for disaster roaming services":
		return Registered, nil
	default:
		return 0, errors.New("invalid Roaming request type")
	}
}
