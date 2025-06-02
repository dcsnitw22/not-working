package nas

import (
	"errors"
)

//Logic for decoding Extended Protocol Discriminator

func DecodeEpd(epd byte) (string, error) {
	switch epd {
	case MobilityManagementEPD:
		return "MOBILITY_MANAGEMENT_MESSAGES", nil
	case SessionManagementEPD:
		return "SESSION_MANAGEMENT_MESSAGES", nil
	default:
		return "", errors.New("unknown EPD")
	}
}
