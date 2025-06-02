package nas

import (
	"errors"
)

//Logic to decode MaximumDataRatePerUeForUserPlaneIntegrityProtectionForUplink

func DecodeMaximumDataRatePerUeForUserPlaneIntegrityProtectionForUplink(dl byte) (string, error) {
	switch dl {
	case Kbps:
		return "SIXTY_FOUR_KBPS", nil
	case FullDataRate:
		return "FULL_DATA_RATE", nil
	default:
		return "", errors.New("invalid MaximumDataRatePerUeForUserPlaneIntegrityProtectionForUplink")
	}
}
