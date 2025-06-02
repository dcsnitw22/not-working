package nas

import "errors"

func EncodePacketFilterDirection(pfd string) (byte, error) {
	switch pfd {
	case "DOWNLINK":
		return Downlink, nil
	case "UPLINK":
		return Uplink, nil
	case "BIDIRECTIONAL":
		return Bidirectional, nil
	}
	return 0, errors.New("invalid Input")
}
