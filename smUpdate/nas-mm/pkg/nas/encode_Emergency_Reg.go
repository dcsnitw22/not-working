package nas

import "errors"

func EncodeEmergencyReg(emReg string) (byte, error) {
	switch emReg {
	case "Not registered for emergency services":
		return NotRegistered, nil
	case "Registered for emergency services":
		return Registered, nil
	default:
		return 0, errors.New("invalid Emergency request type")
	}
}
