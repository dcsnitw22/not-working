package nas

import (
	"errors"
)

//Logic to encode SessionAMBR

func EncodeUnitSessionAMBR(ambr string) (byte, error) {
	switch ambr {
	case "VALUE_NOT_USED":
		return ValNotUsed, nil
	case "MULT_1Kbps":
		return Mult_1kbps, nil
	case "MULT_4Kbps":
		return Mult_4kbps, nil
	case "MULT_16Kbps":
		return Mult_16kbps, nil
	case "MULT_64Kbps":
		return Mult_64kbps, nil
	case "MULT_256kbps":
		return Mult_256kbps, nil
	case "MULT_1Mbps":
		return Mult_1mbps, nil
	case "MULT_4Mbps":
		return Mult_4mbps, nil
	case "MULT_16Mbps":
		return Mult_16mbps, nil
	case "MULT_64Mbps":
		return Mult_64mbps, nil
	case "MULT_256Mbps":
		return Mult_256mbps, nil
	case "MULT_1Gbps":
		return Mult_1gbps, nil
	case "MULT_4Gbps":
		return Mult_4gbps, nil
	case "MULT_16Gbps":
		return Mult_16gbps, nil
	case "MULT_64Gbps":
		return Mult_64gbps, nil
	case "MULT_256Gbps":
		return Mult_256gbps, nil
	case "MULT_1Tbps":
		return Mult_1tbps, nil
	case "MULT_4Tbps":
		return Mult_4tbps, nil
	case "MULT_16Tbps":
		return Mult_16tbps, nil
	case "MULT_64Tbps":
		return Mult_64tbps, nil
	case "MULT_256Tbps":
		return Mult_256tbps, nil
	case "MULT_1Pbps":
		return Mult_1pbps, nil
	case "MULT_4Pbps":
		return Mult_4pbps, nil
	case "MULT_16Pbps":
		return Mult_16pbps, nil
	case "MULT_64Pbps":
		return Mult_64pbps, nil
	case "MULT_256Pbps":
		return Mult_256pbps, nil
	default:
		return 0, errors.New("invalid Session AMBR")
	}
}
