package nas

import (
	"errors"
)

//Logic to decode MaximumDataRatePerUeForUserPlaneIntegrityProtectionForDownlink

func DecodeMaximumDataRatePerUeForUserPlaneIntegrityProtectionForDownlink(dl byte) (string, error) {
	switch dl {
	case Kbps:
		return "SIXTY_FOUR_KBPS", nil
	case FullDataRate:
		return "FULL_DATA_RATE", nil
	default:
		return "", errors.New("invalid MaximumDataRatePerUeForUserPlaneIntegrityProtectionForDownlink")
	}
}
