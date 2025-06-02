package nas

import (
	"errors"
	"strconv"
)

func DecodeKSITSC(tsc byte) (string, error) {
	switch tsc {
	case Native:
		return "Native security context", nil
	case Mapped:
		return "Mapped security context", nil
	default:
		return "", errors.New("invalid TSC")
	}
}

func DecodeKSI(ksi byte) (string, error) {
	switch ksi {
	case NoKey:
		return "No Key is available", nil
	default:
		return strconv.Itoa(int(ksi)), nil
	}
}

func DecodeNASKSI(tscVal byte, ksiVal byte) (string, string, error) {
	decodedTSC, err := DecodeKSITSC(tscVal)
	decodedKSI, _ := DecodeKSI(ksiVal)
	return decodedTSC, decodedKSI, err
}
