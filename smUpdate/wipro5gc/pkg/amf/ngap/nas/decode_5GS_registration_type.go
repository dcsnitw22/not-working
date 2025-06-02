package nas

import "errors"

func Decode5GSRegistrationValue(regValue byte) (string, error) {
	switch regValue {
	case InitialReg:
		return "Initial registration", nil
	case MobilityReg:
		return "Mobility registration updating", nil
	case PeriodicReg:
		return "Periodic registration updating", nil
	case EmergencyReg:
		return "Emergency registration", nil
	case SNPNReg:
		return "SNPN onboarding registration", nil
	case RoamingMobilityReg:
		return "Disaster roaming mobility registration updating", nil
	case RoamingInitialReg:
		return "Disaster roaming initial registration", nil
	default:
		return "Initial Registration", nil

	}
}

func Decode5GSRegistrationFOR(fOR byte) (string, error) {
	switch fOR {
	case NoFOR:
		return "No follow-on request pending", nil
	case FOR:
		return "Follow-on request pending", nil
	default:
		return "", errors.New("invalid FOR")
	}
}

func Decode5GSRegistrationType(regsValue byte, fORe byte) (string, string, error) {
	decodedRegValue, _ := Decode5GSRegistrationValue(regsValue)
	decodedFOR, err := Decode5GSRegistrationFOR(fORe)
	return decodedRegValue, decodedFOR, err
}
