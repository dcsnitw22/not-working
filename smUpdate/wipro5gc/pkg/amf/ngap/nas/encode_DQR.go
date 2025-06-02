package nas

import "errors"

func EncodeDQR(dqr string) (byte, error) {
	switch dqr {
	case "NOT_DEFAULT_QoS_RULE":
		return NotDefaultQoS, nil
	case "DEFAULT_QoS_RULE":
		return DefaultQoS, nil
	}

	return 0, errors.New("invalid Input")

}
