package nas

import (
	"errors"
)

//Logic for encoding Extended Protocol Discriminator

func EncodeEpd(epd string) (byte, error) {
	switch epd {
	case "MOBILITY_MANAGEMENT_MESSAGES":
		return MobilityManagementEPD, nil
	case "SESSION_MANAGEMENT_MESSAGES":
		return SessionManagementEPD, nil
	default:
		return 0, errors.New("unknown EPD")
	}
}
